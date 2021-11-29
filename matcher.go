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
	lettersWithCount map[byte]uint16
	wordLen          uint8
}

func newWordEntry() wordEntry {
	return wordEntry{
		wordIdx:          0,
		word:             []byte{},
		lettersWithCount: map[byte]uint16{},
		wordLen:          0,
	}
}

type Sentence struct {
	inputSentence string
	// Every sentence word will have a index this is ORd together and stored in this field
	// When matching we can do the same and compare the result match values, that way we know if we matched the full sentence
	testMatched  uint64
	wordsWithLen [255]*wordsListWithLen
	minWordLen   uint8

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
	letterCount map[byte]uint16
}

// NewMatcher creates a new matcher that can be used to match
// Note that this method takes time as it optimized the input sentences to be fast to match
func NewMatcher(sentences ...string) *Matcher {
	m := &Matcher{
		sentences:   make([]*Sentence, len(sentences)),
		letterCount: map[byte]uint16{},
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
					currentWord.lettersWithCount[c]++
					currentWord.word = append(currentWord.word, c)
					wordsWithCounter[lastWordIdx] = currentWord
				} else if c < utf8.RuneSelf {
					wordsWithCounter = append(wordsWithCounter, newWordEntry())
					lastWordIdx++
				}
			}
		}

		wordsWithLen := [255]*wordsListWithLen{}
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
				// set the allowed offset
				var allowedOff uint16
				if word.wordLen > 3 {
					if word.wordLen < 10 {
						allowedOff = 1
					} else if word.wordLen < 20 {
						allowedOff = 2
					} else if word.wordLen < 30 {
						allowedOff = 3
					}
				}

				wordsWithLen[word.wordLen] = &wordsListWithLen{
					list:       []wordEntry{},
					allowedOff: allowedOff,
				}
			} else {
				// Check if this word is already in the list
				wordAsStr := b2s(word.word)
				for _, listEntry := range wordsWithLen[word.wordLen].list {
					if b2s(listEntry.word) == wordAsStr {
						continue outerLoop
					}
				}
			}

			list := wordsWithLen[word.wordLen]
			word.wordIdx = wordIndex
			testSentencedMatched |= wordIndex
			wordIndex <<= 1
			list.list = append(list.list, word)
		}

		if minWordLen > 2 {
			minWordLen -= 2
		}
		m.sentences[idx] = &Sentence{
			inputSentence: sentence,
			testMatched:   testSentencedMatched,
			wordsWithLen:  wordsWithLen,
			minWordLen:    minWordLen,
		}
	}

	return m
}

func (m *Matcher) clearLetterCount() {
	if len(m.letterCount) != 0 {
		for k := range m.letterCount {
			delete(m.letterCount, k)
		}
	}
}

func (s *Sentence) wordMatch(m *Matcher, wordLen uint8) uint64 {
	if wordLen < s.minWordLen {
		return 0
	}

	wordsListsToMatchWith := []*wordsListWithLen{
		s.wordsWithLen[wordLen],
		s.wordsWithLen[wordLen-1],
		s.wordsWithLen[wordLen+1],
	}

	// Contains either
	// Null if no words matched
	// The index of a exact match
	// The indexes of all the words that somewhat matched
	potentialWordIndex := uint64(0)

listsLoop:
	for _, list := range wordsListsToMatchWith {
		if list == nil {
			continue
		}

		if list.allowedOff == 0 {
		allowedOffsetZeroWordsLoop:
			for _, word := range list.list {
				for letter, letterCount := range word.lettersWithCount {
					if m.letterCount[letter] != letterCount {
						continue allowedOffsetZeroWordsLoop
					}
				}
				potentialWordIndex = word.wordIdx
				break listsLoop
			}
		} else {
		wordsLoop:
			for _, word := range list.list {
				// off contains the word offset from the currently inspecting word
				var off uint16
				for letter, letterCount := range word.lettersWithCount {
					currentWordLetterCount := m.letterCount[letter]
					if currentWordLetterCount == letterCount {
						continue
					} else if letterCount > currentWordLetterCount {
						off += letterCount - currentWordLetterCount
					} else {
						off += currentWordLetterCount - letterCount
					}
					if off > list.allowedOff {
						continue wordsLoop
					}
				}
				if off == 0 {
					// We found an exact match!
					// As potentialWordIndex is ORed to the matchedWordsIdx later on we write the word index in
					// there so later on the index gets ORed to the matchedWordsIdx
					potentialWordIndex = word.wordIdx
					break listsLoop
				}
				potentialWordIndex |= word.wordIdx
			}
		}

		if potentialWordIndex != 0 {
			// Break the loop early if we have some potential matches
			// Don't waist cpu cycles
			break
		}
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
