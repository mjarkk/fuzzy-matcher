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

## Benchmark

```sh
go test -benchmem -bench "^BenchmarkMatch$"
```

```
goos: darwin
goarch: amd64
pkg: github.com/mjarkk/fuzzy-matcher
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkMatch-12    	  212745	      5408 ns/op	       0 B/op	       0 allocs/op
```

Test code can be found in [matcher_test.go](https://github.com/mjarkk/fuzzy-matcher/blob/main/matcher_test.go)
