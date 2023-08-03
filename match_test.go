package main

import "testing"

type matchResult struct {
	in  string
	res bool
}

type testCase struct {
	regex  string
	inputs []matchResult
}

var cases = []testCase{
	{"a(b|c)*d?", []matchResult{
		{"abcd", true},
		{"ad", true},
		{"a", true},
		{"bcd", false},
	}},
	{"a?a?aa", []matchResult{
		{"aa", true},
		{"aaa", true},
		{"aaaa", true},
		{"a", false},
		{"aaaaa", false},
	}},
	{"(aa|ab|ba|aaa)*", []matchResult{
		{"aa", true},
		{"aba", false},
		{"aaba", true},
		{"aabaa", false},
		{"aabaaa", true},
		{"aaabaa", true},
		{"aa", true},
		{"baaaaaa", true},
	}},
}

func TestMatch(t *testing.T) {
	for _, c := range cases {
		ast, err := parse(c.regex)
		if err != nil {
			t.Errorf("parsing failed: %s", err)
		}
		regex := compile(ast)
		for _, i := range c.inputs {
			if match(&regex, i.in) != i.res {
				t.Errorf("wrong results: %s, %s", c.regex, i.in)
			}
		}
	}
}
