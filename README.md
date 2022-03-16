# `fuzzy-matcher`

A fast golang fuzzy string matcher

```go
// New matcher creates a new matcher to be used for matching
// Note that this operation takes the most time
matcher := fuzzymatcher.NewMatcher(
    "I love trees",
    "bananas are the best fruit",
    "peer",
)

// match returns the best match for the given string or -1 if no match was found
fmt.Println(matcher.Match("do i love the trees") == 0)
```
