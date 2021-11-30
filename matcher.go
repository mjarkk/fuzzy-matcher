package fuzzymatcher

import (
	"unicode/utf8"
)

type wordEntry struct {
	// wordIdx contains one bit sifted to the left for each word in the sentence
	// So wordIdx will be 1, 2, 4, 8, 16, 32, ...
	// This also means the wordIdx can be a maximum of 64 words
	wordIdx uint64

	word             []byte
	lettersWithCount []letter
	wordLen          uint8
}

type letter struct {
	letter byte
	count  uint16
}

func newWordEntry() wordEntry {
	return wordEntry{
		wordIdx:          0,
		word:             []byte{},
		lettersWithCount: []letter{},
		wordLen:          0,
	}
}

type Sentence struct {
	inputSentence string
	// Every sentence word will have a index this is ORd together and stored in this field
	// When matching we can do the same and compare the result match values, that way we know if we matched the full sentence
	testMatched uint64

	// Contains a list of potential words for every word length from 1-254
	// It contains the words maching -1 to 1 length so this list will have duplicated entries
	wordsWithLen [255]*wordsListWithLen

	minWordLen uint8

	// Used by matcher
	matchedWordsIndex uint64
}

type wordsListWithLen struct {
	allowedOff uint16
	list       []wordEntry
}

type Matcher struct {
	sentences []*Sentence

	// Used by matcher
	letterCount [256]uint16
}

// NewMatcher creates a new matcher that can be used to match
// Note that this method takes time as it optimized the input sentences to be fast to match
func NewMatcher(sentences ...string) *Matcher {
	m := &Matcher{
		sentences: make([]*Sentence, len(sentences)),
	}

	for idx, sentence := range sentences {
		wordsWithCounter := []wordEntry{newWordEntry()}
		lastWordIdx := 0

		for _, c := range s2b(sentence) {
			if c >= 'B' && c <= 'Z' {
				c += 'b' - 'B'
			}
			switch c {
			case 'a', 'e', 'i', 'o', 'u':
				continue
			default:
				if (c >= 'b' && c <= 'z') || (c >= '0' && c <= '9') {
					currentWord := wordsWithCounter[lastWordIdx]
					if currentWord.wordLen < 254 {
						currentWord.wordLen++
					}

					foundLetter := false
					for idx, letter := range currentWord.lettersWithCount {
						if letter.letter == c {
							letter.count++
							currentWord.lettersWithCount[idx] = letter
							foundLetter = true
							break
						}
					}
					if !foundLetter {
						currentWord.lettersWithCount = append(currentWord.lettersWithCount, letter{
							letter: c,
							count:  1,
						})
					}
					currentWord.word = append(currentWord.word, c)
					wordsWithCounter[lastWordIdx] = currentWord
				} else if c < utf8.RuneSelf {
					wordsWithCounter = append(wordsWithCounter, newWordEntry())
					lastWordIdx++
				}
			}
		}

		wordsWithLen := [255][]wordEntry{}
		wordIndex := uint64(1)
		testSentencedMatched := uint64(0)

		minWordLen := uint8(254)
	outerLoop:
		for idx := len(wordsWithCounter) - 1; idx >= 0; idx-- {
			word := wordsWithCounter[idx]
			if word.wordLen == 0 {
				continue
			}
			if word.wordLen < minWordLen {
				minWordLen = word.wordLen
			}

			if wordsWithLen[word.wordLen] == nil {
				wordsWithLen[word.wordLen] = []wordEntry{}
			} else {
				// Check if this word is already in the list
				wordAsStr := b2s(word.word)
				for _, listEntry := range wordsWithLen[word.wordLen] {
					if b2s(listEntry.word) == wordAsStr {
						continue outerLoop
					}
				}
			}

			word.wordIdx = wordIndex
			testSentencedMatched |= wordIndex
			wordIndex <<= 1
			wordsWithLen[word.wordLen] = append(wordsWithLen[word.wordLen], word)
		}

		wordsWithLenWithDiff := [255]*wordsListWithLen{}
		for idx, list := range wordsWithLen {
			newList := list

			if idx != 0 && wordsWithLen[idx-1] != nil {
				newList = append(newList, wordsWithLen[idx-1]...)
			}
			if idx != 254 && wordsWithLen[idx+1] != nil {
				newList = append(newList, wordsWithLen[idx+1]...)
			}

			if len(newList) == 0 {
				continue
			}

			var allowedOff uint16
			if idx > 3 {
				if idx < 10 {
					allowedOff = 1
				} else if idx < 20 {
					allowedOff = 2
				} else {
					allowedOff = 3
				}
			}

			wordsWithLenWithDiff[idx] = &wordsListWithLen{
				allowedOff: allowedOff,
				list:       newList,
			}
		}

		if minWordLen > 2 {
			minWordLen -= 2
		}
		m.sentences[idx] = &Sentence{
			inputSentence: sentence,
			testMatched:   testSentencedMatched,
			wordsWithLen:  wordsWithLenWithDiff,
			minWordLen:    minWordLen,
		}
	}

	return m
}

func (m *Matcher) clearLetterCount() {
	m.letterCount = [256]uint16{}
}

func (s *Sentence) wordMatch(m *Matcher, wordLen uint8) uint64 {
	if wordLen < s.minWordLen {
		return 0
	}

	// Contains either
	// Null if no words matched
	// The index of a exact match
	// The indexes of all the words that somewhat matched
	potentialWordIndex := uint64(0)

	list := s.wordsWithLen[wordLen]
	if list == nil {
		return potentialWordIndex
	}

	if list.allowedOff == 0 {
		// If no offset is allowed, we can skip a lot of logic
		// That's why we have a diffrent loop for these
	allowedOffsetZeroWordsLoop:
		for _, word := range list.list {
			for _, letterAndCount := range word.lettersWithCount {
				if m.letterCount[letterAndCount.letter] != letterAndCount.count {
					continue allowedOffsetZeroWordsLoop
				}
			}
			return word.wordIdx
		}
		return potentialWordIndex
	}

wordsLoop:
	for _, word := range list.list {
		// off contains the word offset from the currently inspecting word
		var off uint16
		for _, letterAndCount := range word.lettersWithCount {
			currentWordLetterCount := m.letterCount[letterAndCount.letter]
			if currentWordLetterCount == letterAndCount.count {
				continue
			} else if letterAndCount.count > currentWordLetterCount {
				off += letterAndCount.count - currentWordLetterCount
			} else {
				off += currentWordLetterCount - letterAndCount.count
			}
			if off > list.allowedOff {
				continue wordsLoop
			}
		}
		if off == 0 {
			// We found an exact match, lets return that
			return word.wordIdx
		}
		potentialWordIndex |= word.wordIdx
	}

	return potentialWordIndex
}

// Match matches the inStr against the matcher inputs and returns the best matching string or -1 if nothing could be matched
func (m *Matcher) Match(inStr string) int {
	in := s2b(inStr)
	currentWordLen := uint8(0)
	result := -1

	for _, s := range m.sentences {
		s.matchedWordsIndex = 0
	}

	doMatch := func() bool {
		if currentWordLen == 0 {
			return false
		}

		for idx, sentence := range m.sentences {
			sentence.matchedWordsIndex |= sentence.wordMatch(m, currentWordLen)
			if sentence.matchedWordsIndex == sentence.testMatched {
				m.clearLetterCount()
				currentWordLen = 0
				result = idx
				return true
			}
		}

		m.clearLetterCount()
		currentWordLen = 0
		return false
	}

	for _, c := range in {
		if c >= 'B' && c <= 'Z' {
			c += 'b' - 'B'
		}
		switch c {
		case 'a', 'e', 'i', 'o', 'u':
			continue
		default:
			if (c >= 'b' && c <= 'z') || (c >= '0' && c <= '9') || c >= utf8.RuneSelf {
				m.letterCount[c]++
				if currentWordLen != 253 {
					currentWordLen++
				}
			} else if currentWordLen > 0 && doMatch() {
				return result
			}
		}
	}

	doMatch()
	return result
}
