package fuzzymatcher

import (
	"unicode/utf8"
)

type wordEntry struct {
	// wordIdx contains one bit sifted to the left for each word in the sentence
	// So wordIdx will be 1, 2, 4, 8, 16, 32, ...
	// This also means the wordIdx can be a maximum of 64 words
	WordIdx uint64
	Letters []rune

	len               int
	allowedOffset     int
	FuzzyFirstLetter  [3]rune
	FuzzyLettersOrder [][3]rune
}

func (we *wordEntry) letterAt(idx int) rune {
	if idx < 0 || idx >= len(we.Letters) {
		return 0
	}
	return we.Letters[idx]
}

func (we *wordEntry) calculateFuzzyLetterOrder() {
	we.len = len(we.Letters)
	if we.len <= 4 {
		we.allowedOffset = 1
	} else if we.len <= 7 {
		we.allowedOffset = 2
	} else {
		we.allowedOffset = 3
	}

	we.FuzzyLettersOrder = [][3]rune{}

	lettersLen := len(we.Letters)
	possibleNextLetters := [3]rune{}
	for i := 0; i < lettersLen-1; i++ {
		possibleNextLetters = [3]rune{we.letterAt(i + 1), we.letterAt(i + 2), we.letterAt(i + 3)}
		if we.allowedOffset < 3 {
			possibleNextLetters[2] = 0
			if we.allowedOffset < 2 {
				possibleNextLetters[1] = 0
			}
		}
		we.FuzzyLettersOrder = append(we.FuzzyLettersOrder, possibleNextLetters)
	}

	we.FuzzyFirstLetter = [3]rune{we.letterAt(0)}
	nextLetter := we.letterAt(1)
	if nextLetter != utf8.RuneError && we.allowedOffset >= 2 {
		we.FuzzyFirstLetter[1] = nextLetter

		nextLetter := we.letterAt(2)
		if nextLetter != utf8.RuneError && we.allowedOffset >= 3 {
			we.FuzzyFirstLetter[2] = nextLetter
		}
	}
}

type pathToWord struct {
	Letter             rune
	Sentence           int
	Word               int
	WordOffset         int
	MustRemainingChars int
}

type sentenceT struct {
	Words                []wordEntry
	IdxInNewMatcherInput int

	// the fields below are generated with the (*sentence).complete() method
	Paths       []pathToWord
	IndexSum    uint64
	SentenceLen int

	// Used in the matching process
	MatchIndexSum uint64
}

func (s *sentenceT) complete() {
	s.Paths = []pathToWord{}
	for wordIdx, word := range s.Words {
		for offset, letter := range word.FuzzyFirstLetter {
			if letter == 0 {
				break
			}
			s.Paths = append(s.Paths, pathToWord{
				Letter:             letter,
				Sentence:           -1,
				Word:               wordIdx,
				WordOffset:         offset,
				MustRemainingChars: word.len - word.allowedOffset - 1,
			})
		}
		s.IndexSum |= word.WordIdx
		s.SentenceLen += word.len
		if wordIdx != len(s.Words)-1 {
			// Also add a space character for the
			s.SentenceLen++
		}
	}
}

// Matcher is used to match sentences
type Matcher struct {
	Sentences []sentenceT

	// the fields below are generated with the (*Matcher).complete() method
	Paths                []pathToWord
	HasPathsWithRuneSelf bool                        // basicly tells if there are complex utf8 chars
	PathByLetterMap      map[rune][]pathToWord       // Use if HasPathsWithRuneSelf == true
	PathByLetterList     [utf8.RuneSelf][]pathToWord // Use if HasPathsWithRuneSelf == false

	// Zero alloc cache
	UTF8RuneCreation  []byte
	InProgressMatches []inProgressMatch
}

func (m *Matcher) complete() {
	m.Paths = []pathToWord{}
	m.PathByLetterMap = map[rune][]pathToWord{}
	m.PathByLetterList = [utf8.RuneSelf][]pathToWord{}

	for idx, sentence := range m.Sentences {
		for _, path := range sentence.Paths {
			path.Sentence = idx

			letter := path.Letter

			m.Paths = append(m.Paths, path)
			list, ok := m.PathByLetterMap[letter]
			// Add the path to a specific paths list or create a new paths list
			if !ok {
				m.PathByLetterMap[letter] = []pathToWord{path}
			} else {
				m.PathByLetterMap[letter] = append(list, path)
			}

			if letter < utf8.RuneSelf {
				m.PathByLetterList[letter] = append(m.PathByLetterList[letter], path)
			} else {
				m.HasPathsWithRuneSelf = true
			}
		}
	}
}

const upperToLowerCaseOffset = 'a' - 'A'

// NewMatcher creates a new instance of the matcher
// This function takes relatively long to execute so do this once, and use the returned matcher to match it against lots of entries
func NewMatcher(sentences ...string) *Matcher {
	res := Matcher{
		Sentences:         []sentenceT{},
		UTF8RuneCreation:  []byte{},
		InProgressMatches: []inProgressMatch{},
	}

	for sentenceIdx, sentence := range sentences {
		parsedSentence := sentenceT{
			Words:                []wordEntry{},
			IdxInNewMatcherInput: sentenceIdx,
		}

		word := wordEntry{WordIdx: 1}
		commitWord := func() {
			if len(word.Letters) <= 1 {
				// Just reset the current word
				word = wordEntry{WordIdx: word.WordIdx}
			} else {
				word.calculateFuzzyLetterOrder()
				parsedSentence.Words = append(parsedSentence.Words, word)
				word = wordEntry{WordIdx: word.WordIdx << 1}
			}
		}

		for _, c := range []rune(sentence) {
			if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
				word.Letters = append(word.Letters, c)
			} else if c >= 'A' && c <= 'Z' {
				word.Letters = append(word.Letters, c+upperToLowerCaseOffset)
			} else if c >= utf8.RuneSelf {
				newC, ok := checkAndCorredUnicodeChar(c)
				if ok {
					word.Letters = append(word.Letters, newC)
				}
			} else {
				commitWord()
			}
		}
		commitWord()

		if len(parsedSentence.Words) == 0 {
			continue
		}

		parsedSentence.complete()
		res.Sentences = append(res.Sentences, parsedSentence)
	}

	res.complete()
	return &res
}

type inProgressMatch struct {
	PathToWord   pathToWord
	Word         *wordEntry
	Sentence     *sentenceT
	WordOffset   int
	SkippedChars int
}

func (e *inProgressMatch) addWordIdxToSentence() int {
	e.Sentence.MatchIndexSum |= e.Word.WordIdx
	if e.Sentence.MatchIndexSum == e.Sentence.IndexSum {
		return e.Sentence.IdxInNewMatcherInput
	}
	return -1
}

// Match matches a sentence to the matchers input
// Returns the index of the matched sentenced
// If nothing found returns -1
func (m *Matcher) Match(sentence string) int {
	// Reset the matching index sums and zero alloc cache
	for idx := range m.Sentences {
		m.Sentences[idx].MatchIndexSum = 0
	}
	m.InProgressMatches = m.InProgressMatches[:0]

	sentenceLen := len(sentence)
	var rLetter rune

	beginWord := true
	for i := 0; i < sentenceLen; i++ {
		letter := sentence[i]
		if letter == 0 {
			continue
		}

		if letter >= utf8.RuneSelf {
			if beginWord && !m.HasPathsWithRuneSelf {
				// There are not even words starting with this letter, lets skip this one
				beginWord = false
				continue
			}
			if !beginWord && len(m.InProgressMatches) == 0 {
				// We are matching nothing on the current word, no need to execute heavy instructions
				continue
			}

			switch i {
			case sentenceLen - 1:
				// This is a invalid character
				rLetter = utf8.RuneError
			case sentenceLen - 2:
				r, size := utf8.DecodeRune(append(m.UTF8RuneCreation[:0], letter, sentence[i+1]))
				i += size - 1
				rLetter, _ = checkAndCorredUnicodeChar(r)
			case sentenceLen - 3:
				r, size := utf8.DecodeRune(append(m.UTF8RuneCreation[:0], letter, sentence[i+1], sentence[i+2]))
				i += size - 1
				rLetter, _ = checkAndCorredUnicodeChar(r)
			default:
				r, size := utf8.DecodeRune(append(m.UTF8RuneCreation[:0], letter, sentence[i+1], sentence[i+2], sentence[i+3]))
				i += size - 1
				rLetter, _ = checkAndCorredUnicodeChar(r)
			}
			if rLetter == utf8.RuneError {
				continue
			}
		} else {
			rLetter = rune(sentence[i])

			if (rLetter >= 'a' && rLetter <= 'z') || (rLetter >= '0' && rLetter <= '9') {
				// Do nothing
			} else if rLetter >= 'A' && rLetter <= 'Z' {
				rLetter += upperToLowerCaseOffset
			} else {
				// go to next word

				// Firstly lets check if there where any matches from the last word
				for _, entry := range m.InProgressMatches {
					// Check if we mis the last chars
					// If so this entry is oke
					// Makes sure "banan" can match "banana"
					if len(entry.Word.FuzzyLettersOrder)-entry.WordOffset-1 <= entry.Word.allowedOffset-entry.SkippedChars {
						res := entry.addWordIdxToSentence()
						if res != -1 {
							return res
						}
					}
				}

				// Reset the m.InProgressMatches so we can scan for new words
				m.InProgressMatches = m.InProgressMatches[:0]
				beginWord = true
				continue
			}
		}

		if beginWord {
			var paths []pathToWord
			if !m.HasPathsWithRuneSelf || rLetter < utf8.RuneSelf {
				paths = m.PathByLetterList[rLetter]
			} else {
				paths = m.PathByLetterMap[rLetter]
			}

			for _, path := range paths {
				if sentenceLen-i-1 >= path.MustRemainingChars {
					sentence := &m.Sentences[path.Sentence]
					word := &sentence.Words[path.Word]

					if sentence.MatchIndexSum&word.WordIdx != 0 {
						// This word was earlier already matched
						continue
					}

					m.InProgressMatches = append(m.InProgressMatches, inProgressMatch{
						PathToWord:   path,
						Word:         word,
						Sentence:     sentence,
						WordOffset:   path.WordOffset,
						SkippedChars: path.WordOffset,
					})
				}
			}

			beginWord = false
			continue
		}

	outer:
		for i := len(m.InProgressMatches) - 1; i >= 0; i-- {
			entry := m.InProgressMatches[i]
			for offset, c := range entry.Word.FuzzyLettersOrder[entry.WordOffset] {
				if c == rLetter {
					if offset > 0 && offset > entry.Word.allowedOffset-entry.SkippedChars {
						continue
					}

					entry.WordOffset += offset + 1
					entry.SkippedChars += offset
					if entry.WordOffset == len(entry.Word.FuzzyLettersOrder) {
						// Completed matching this word
						res := entry.addWordIdxToSentence()
						if res != -1 {
							return res
						}
						m.InProgressMatches = append(m.InProgressMatches[:i], m.InProgressMatches[i+1:]...)
					} else {
						m.InProgressMatches[i] = entry
					}
					continue outer
				}

				if c == 0 {
					break
				}
			}

			if entry.SkippedChars < entry.Word.allowedOffset {
				entry.SkippedChars++
				m.InProgressMatches[i] = entry
			} else {
				m.InProgressMatches = append(m.InProgressMatches[:i], m.InProgressMatches[i+1:]...)
			}
		}
	}

	for _, entry := range m.InProgressMatches {
		// Check if we mis the last chars
		// If so this entry is oke
		// Makes sure "banan" can match "banana"
		if len(entry.Word.FuzzyLettersOrder)-entry.WordOffset-1 <= entry.Word.allowedOffset-entry.SkippedChars {
			res := entry.addWordIdxToSentence()
			if res != -1 {
				return res
			}
		}
	}

	return -1
}

func checkAndCorredUnicodeChar(c rune) (rune, bool) {
	switch c {
	case 'à', 'À', 'á', 'Á', 'â', 'Â', 'ã', 'Ã', 'ä', 'Ä', 'å', 'Å', 'æ', 'Æ':
		return 'a', true
	case 'è', 'È', 'é', 'É', 'ê', 'Ê', 'ë', 'Ë':
		return 'e', true
	case 'ì', 'Ì', 'í', 'Í', 'î', 'Î', 'ï', 'Ï':
		return 'i', true
	case 'ò', 'Ò', 'ó', 'Ó', 'ô', 'Ô', 'õ', 'Õ', 'ö', 'Ö', 'ð', 'Ð', 'ø', 'Ø':
		return 'o', true
	case 'ù', 'Ù', 'ú', 'Ú', 'û', 'Û', 'ü', 'Ü':
		return 'u', true
	case 'ß':
		return 's', true
	case 'ñ', 'Ñ':
		return 'n', true
	case 'ý', 'Ý', 'ÿ', 'Ÿ':
		return 'y', true
	case 'ç', 'Ç':
		return 'c', true
	case '©':
		return 'c', true
	case '®':
		return 'r', true
	case 768, // accent of: à
		769, // accent of: á
		770, // accent of: â
		771, // accent of: ã
		776, // accent of: ä
		778, // accent of: å
		'¿',
		'¡',
		0x2002, // En space
		0x2003, // Em space
		0x2004, // Three-per-em space
		0x2005, // Four-per-em space
		0x2006, // Six-per-em space
		0x2007, // Figure space
		0x2008, // Punctuation space
		0x2009, // Thin space
		0x200A, // Hair space
		0x200B, // Zero width space
		0x202F, // Narrow no-break space
		0x205F, // Medium mathematical space
		0x3000, // Ideographic space
		'“',
		'”',
		'’',
		'‵',
		'‹',
		'›',
		'»',
		'«',
		utf8.RuneError:
		return utf8.RuneError, false
	default:
		return c, true
	}
}
