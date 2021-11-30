package fuzzymatcher

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	testCases := []struct {
		input       string
		matchWith   string
		shouldMatch bool
	}{
		{"banana", "banana", true},
		{"banana", "banan", true},
		{"banana", "banaana", true},
		{"banana", "bananas", true},
		{"banana", "i want a banana", true},
		{"some thing", "thing some", true},
		{"says pet", "i love food says the pet", true},
		{"love i", "i love food says the pet", true},
		{"bananen lekker", "ik vind bananen erg lekker", true},
		{
			"this is a very long sentence",
			"another sentence that contains the other sentence \"this is a very long sentence\" so there should be a match",
			true,
		},

		{"somewhere over the rainbow", "somewhere", false},
		{"banana", "apple", false},
		{"bananen lekker", "bananen zijn vies", false},
	}

	for _, testCase := range testCases {
		parsedInput := NewMatcher(testCase.input)
		if testCase.shouldMatch {
			a.Equal(t, 0, parsedInput.Match(testCase.matchWith), `Expected "%s" to match "%v"`, testCase.input, testCase.matchWith)
		} else {
			a.Equal(t, -1, parsedInput.Match(testCase.matchWith), `Expected "%s" to not match "%v"`, testCase.input, testCase.matchWith)
		}
	}

	matcher := NewMatcher(
		"I love trees",
		"bananas are the best fruit",
		"banana",
	)

	matchesWith := []struct {
		matches int
		input   string
	}{
		{-1, "nothing"},
		{0, "i love trees"},
		{2, "bananas are the best fruit"},
		{2, "banana"},
		{0, "do you also love trees? i do."},
		{2, "on a sunday afternoon i like to eat a banana"},
	}

	for _, testCase := range matchesWith {
		a.Equal(t, testCase.matches, matcher.Match(testCase.input), testCase.input)
	}
}

func BenchmarkMatchWithProfile(b *testing.B) {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.Fatal(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	BenchmarkMatch(b)

	pprof.StopCPUProfile()
	f.Close()
}

func BenchmarkMatch(b *testing.B) {
	// BenchmarkMatch-12    	 1897522	       608.7 ns/op	       0 B/op	       0 allocs/op

	matcher := NewMatcher(
		"I love trees",
		"bananas are the best fruit",
		"banana",
	)

	matchesWith := []string{
		"nothing",
		"i love trees",
		"bananas are the best fruit",
		"banana",
		"do you also love trees? i do.",
		"on a sunday afternoon i like to eat a banana",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range matchesWith {
			matcher.Match(v)
		}
	}
}
