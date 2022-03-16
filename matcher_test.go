package fuzzymatcher

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

const lordemIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Id faucibus nisl tincidunt eget nullam non. Bibendum neque egestas congue quisque egestas. Quis hendrerit dolor magna eget. Adipiscing elit duis tristique sollicitudin nibh. Venenatis tellus in metus vulputate eu scelerisque felis. Neque volutpat ac tincidunt vitae semper quis lectus nulla at. Donec et odio pellentesque diam volutpat commodo sed egestas egestas. Facilisis sed odio morbi quis commodo odio. Viverra ipsum nunc aliquet bibendum enim facilisis. Molestie ac feugiat sed lectus vestibulum mattis. Quisque id diam vel quam elementum pulvinar etiam. Sed felis eget velit aliquet sagittis id consectetur purus ut. In egestas erat imperdiet sed euismod nisi. Mollis nunc sed id semper risus. Nec feugiat nisl pretium fusce id velit. Amet luctus venenatis lectus magna fringilla urna.

Ullamcorper a lacus vestibulum sed arcu non. Tellus mauris a diam maecenas. Quis auctor elit sed vulputate. Ullamcorper malesuada proin libero nunc consequat interdum varius sit. Etiam erat velit scelerisque in dictum. Auctor urna nunc id cursus metus aliquam eleifend. Aliquam id diam maecenas ultricies mi. Arcu bibendum at varius vel pharetra vel turpis. Ornare arcu dui vivamus arcu felis bibendum ut. Vel risus commodo viverra maecenas. Duis convallis convallis tellus id. Dictumst vestibulum rhoncus est pellentesque elit. Faucibus pulvinar elementum integer enim neque volutpat. At consectetur lorem donec massa sapien faucibus et. Ultrices in iaculis nunc sed.

Erat pellentesque adipiscing commodo elit. Massa sapien faucibus et molestie ac feugiat. Tempus urna et pharetra pharetra massa massa. Cras sed felis eget velit. Nec nam aliquam sem et. Sed blandit libero volutpat sed cras ornare arcu dui vivamus. Aliquam vestibulum morbi blandit cursus. Morbi tincidunt augue interdum velit euismod in. Gravida cum sociis natoque penatibus et. Sollicitudin tempor id eu nisl. Porttitor eget dolor morbi non. Turpis cursus in hac habitasse platea dictumst quisque. Pharetra et ultrices neque ornare aenean euismod elementum nisi. In fermentum posuere urna nec tincidunt praesent. Ornare lectus sit amet est placerat in egestas erat imperdiet. Ultricies mi quis hendrerit dolor magna eget est lorem.

Ac tortor vitae purus faucibus ornare. Non diam phasellus vestibulum lorem sed risus. Enim neque volutpat ac tincidunt vitae semper quis. Quis blandit turpis cursus in hac habitasse platea dictumst quisque. Urna et pharetra pharetra massa massa ultricies mi quis. Phasellus vestibulum lorem sed risus ultricies tristique nulla aliquet. Fames ac turpis egestas integer eget aliquet. Tincidunt eget nullam non nisi est sit amet. Velit aliquet sagittis id consectetur. Sapien faucibus et molestie ac feugiat sed. Elit sed vulputate mi sit amet mauris commodo quis imperdiet. Vitae nunc sed velit dignissim sodales.

Blandit volutpat maecenas volutpat blandit aliquam etiam erat velit scelerisque. In dictum non consectetur a. Pulvinar sapien et ligula ullamcorper. Arcu odio ut sem nulla pharetra diam. Vel risus commodo viverra maecenas accumsan lacus vel facilisis volutpat. Nunc vel risus commodo viverra maecenas accumsan lacus vel. Mattis nunc sed blandit libero volutpat sed cras. Porttitor lacus luctus accumsan tortor posuere ac. In hendrerit gravida rutrum quisque. Feugiat scelerisque varius morbi enim nunc. Tempor orci eu lobortis elementum nibh. Mattis enim ut tellus elementum sagittis vitae. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper velit sed. Ultricies integer quis auctor elit sed.

Eget arcu dictum varius duis at consectetur. Vitae elementum curabitur vitae nunc sed velit dignissim sodales ut. Odio aenean sed adipiscing diam donec adipiscing. Diam ut venenatis tellus in metus vulputate eu. Pharetra et ultrices neque ornare aenean euismod elementum. In dictum non consectetur a erat nam at lectus urna. Eget egestas purus viverra accumsan in nisl nisi. Felis donec et odio pellentesque diam volutpat commodo. Nulla facilisi morbi tempus iaculis. Orci eu lobortis elementum nibh. Massa tincidunt dui ut ornare lectus. Tincidunt praesent semper feugiat nibh sed pulvinar proin gravida. Nunc eget lorem dolor sed viverra. Tempor orci eu lobortis elementum nibh tellus molestie nunc non. Aliquam etiam erat velit scelerisque in dictum non. Non tellus orci ac auctor augue mauris augue neque. Sit amet massa vitae tortor condimentum lacinia quis vel eros.

Ac turpis egestas maecenas pharetra convallis posuere. Porttitor leo a diam sollicitudin tempor id. Nisl pretium fusce id velit ut tortor pretium. Purus non enim praesent elementum facilisis leo vel fringilla. Scelerisque mauris pellentesque pulvinar pellentesque habitant morbi tristique senectus. Cras sed felis eget velit aliquet sagittis id consectetur. Diam sollicitudin tempor id eu nisl nunc mi ipsum faucibus. Egestas tellus rutrum tellus pellentesque eu tincidunt tortor aliquam nulla. Auctor neque vitae tempus quam pellentesque nec nam. Vel eros donec ac odio tempor orci. Lorem ipsum dolor sit amet consectetur adipiscing elit pellentesque. Semper risus in hendrerit gravida. Lacus luctus accumsan tortor posuere ac. Nulla facilisi morbi tempus iaculis urna id volutpat lacus laoreet. Velit euismod in pellentesque massa placerat duis ultricies. Ut morbi tincidunt augue interdum velit euismod in pellentesque. Felis imperdiet proin fermentum leo vel orci porta non pulvinar. Id donec ultrices tincidunt arcu non. Morbi leo urna molestie at elementum eu facilisis sed odio.

A cras semper auctor neque vitae tempus quam pellentesque nec. Ac feugiat sed lectus vestibulum mattis ullamcorper velit. Dolor sed viverra ipsum nunc aliquet. Dolor morbi non arcu risus quis varius. Purus gravida quis blandit turpis cursus in hac habitasse platea. Proin nibh nisl condimentum id venenatis a condimentum. Tortor id aliquet lectus proin nibh. Est lorem ipsum dolor sit. At tellus at urna condimentum. Nec dui nunc mattis enim ut tellus elementum sagittis vitae. Cursus vitae congue mauris rhoncus aenean. Nisl suscipit adipiscing bibendum est ultricies integer. Adipiscing tristique risus nec feugiat in fermentum posuere urna nec. Volutpat sed cras ornare arcu. Pellentesque nec nam aliquam sem et tortor consequat id. Diam volutpat commodo sed egestas egestas fringilla phasellus faucibus scelerisque.

Leo in vitae turpis massa sed elementum tempus. Et egestas quis ipsum suspendisse ultrices gravida. Ipsum a arcu cursus vitae congue mauris rhoncus aenean vel. Volutpat lacus laoreet non curabitur gravida. Sed nisi lacus sed viverra. Purus in massa tempor nec feugiat nisl. Nunc vel risus commodo viverra. Aliquam faucibus purus in massa. Donec ultrices tincidunt arcu non. Pulvinar pellentesque habitant morbi tristique senectus et netus et malesuada.

Aliquet sagittis id consectetur purus ut faucibus pulvinar. Vitae suscipit tellus mauris a diam maecenas sed enim. Duis tristique sollicitudin nibh sit amet commodo. Arcu dictum varius duis at consectetur lorem donec massa. Ut tellus elementum sagittis vitae et leo duis. Risus nullam eget felis eget nunc lobortis mattis aliquam. Ut morbi tincidunt augue interdum. Venenatis urna cursus eget nunc scelerisque viverra mauris in. Quam viverra orci sagittis eu volutpat odio facilisis mauris sit. Sed cras ornare arcu dui vivamus arcu felis. Bibendum est ultricies integer quis auctor elit sed vulputate. Aliquam etiam erat velit scelerisque in dictum non consectetur a. Consequat mauris nunc congue nisi. Eget mauris pharetra et ultrices.`

func TestNewMatcher(t *testing.T) {
	m := NewMatcher("foo")
	a.Len(t, m.Paths, 1)

	a.Len(t, m.Sentences, 1)
	sentence := m.Sentences[0]
	a.Len(t, sentence.Words, 1)
	a.Equal(t, uint64(1), sentence.IndexSum)
	a.Len(t, sentence.Paths, 1)
	a.NotEqual(t, 0, sentence.SentenceLen)

	m = NewMatcher("foo bar   fooBar")
	a.Len(t, m.Paths, 3+1) // 3 words + 1 extra letter for "fooBar"

	a.Len(t, m.Sentences, 1)
	sentence = m.Sentences[0]
	a.Len(t, sentence.Words, 3)
	a.Equal(t, uint64(1+1<<1+1<<2), sentence.IndexSum)
	a.Len(t, sentence.Paths, 3+1) // 3 words + 1 extra letter for "fooBar"
	a.NotEqual(t, 0, sentence.SentenceLen)

	m = NewMatcher("foo", "bar")
	a.Len(t, m.Paths, 2)

	NewMatcher("banana", "i like peers", "foo bar  baz", "another entry that is somwhat long")
	NewMatcher(lordemIpsum)
}

func TestSimpleMatch(t *testing.T) {
	a.Equal(t, 0, NewMatcher("foo").Match("foo"))
	a.Equal(t, -1, NewMatcher("foo").Match("bar"))
}

func TestSimpleMultiMatch(t *testing.T) {
	m := NewMatcher(
		"I love trees",
		"bananas are the best fruit",
		"banana",
	)
	a.Equal(t, 2, m.Match("banana"))
	a.Equal(t, 2, m.Match("bananas are the best fruit"))
	a.Equal(t, 2, m.Match("are the best fruit bananas?"))
}

func TestMatching(t *testing.T) {
	testCases := []struct {
		input       string
		matchWith   string
		shouldMatch bool
	}{
		{"banana", "banana", true},
		{"banana", "banan", true},
		{"banana", "banaana", true},
		{"banana", "bananas", true},
		{"coördinator", "coordinator", true},
		{"coordinator", "coördinator", true},
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
		{"Metselaar", "slijterij", false},
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
		{2, "banana"},
		{2, "bananas are the best fruit"},
		{2, "are the best fruit bananas?"},
		{0, "do you also love trees? i do."},
		{2, "on a sunday afternoon i like to eat a bänanã"},
	}

	for _, testCase := range matchesWith {
		a.Equal(t, testCase.matches, matcher.Match(testCase.input), testCase.input)
	}
}
