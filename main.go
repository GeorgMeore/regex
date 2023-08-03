package main

import "fmt"

func test(regex string) {
	n, e := parse(regex)
	if e == nil {
		fmt.Println(n, compile(n))
	} else {
		fmt.Println(e)
	}
}

func main() {
	test("a(b|c)*d?")
	test("a(b?|c?).d")
	test("a")
	test("")
	test(")a")
	test("a(bc))")
}
