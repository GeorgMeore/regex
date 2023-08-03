package main

import "fmt"

func test(regex string, inputs ...string) {
	n, e := parse(regex)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(n)
	comp := compile(n)
	for _, input := range inputs {
		fmt.Println(input, match(&comp, input))
	}
}

func main() {
	test("a(b|c)*d?", "abd", "abcd", "abcbc", "ad", "a", "aa")
	test("a(b?|c?).d")
	test("a", "a", "aa", "")
	test("")
	test(")a")
	test("a(bc))")
}
