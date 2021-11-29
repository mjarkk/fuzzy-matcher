package main

import "fmt"

func main() {
	matcher := NewMatcher("bananen lekker")
	fmt.Println("should match")
	fmt.Println(matcher.Match("bananen lekker"))             // full match
	fmt.Println(matcher.Match("bananen lekkere"))            // with spell error
	fmt.Println(matcher.Match("ik vind bananen erg lekker")) // extra words

	fmt.Println("should not match")
	fmt.Println(matcher.Match("bananen"))
	fmt.Println(matcher.Match("lekker"))
	fmt.Println(matcher.Match("bananen zijn vies"))
	fmt.Println(matcher.Match("compleet andere string"))
}
