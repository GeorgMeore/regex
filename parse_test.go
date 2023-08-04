package main

import "testing"

type testCase struct {
	in  string
	out string
	err string
}

var cases = []testCase{
	{in: "a(b|c)*d?",  out: "concat{a,concat{star{alt{b,c}},opt{d}}}"},
	{in: "a(b?|c?)cd", out: "concat{a,concat{alt{opt{b},opt{c}},concat{c,d}}}"},
	{in: ")abc",       err: "wanted: a letter or '(' or '\\', got: ')'"},
	{in: "",           err: "wanted: a letter or '(' or '\\', got: EOF"},
	{in: "?",          err: "wanted: a letter or '(' or '\\', got: '?'"},
	{in: "(*)",        err: "wanted: a letter or '(' or '\\', got: '*'"},
	{in: "a(bc",       err: "wanted: ')', got: EOF"},
	{in: "((a)))",     err: "wanted: EOF, got: ')'"},
	{in: "a\\",        err: "wanted: a character, got: EOF"},
}

func TestParse(t *testing.T) {
	for _, c := range cases {
		res, err := parse(c.in)
		if err != nil && err.Error() != c.err {
			t.Errorf("bad error message: expected %s, got %s", c.err, err.Error())
		}
		if res != nil && res.String() != c.out {
			t.Errorf("bad value: expected %s, got %s", c.out, res.String())
		}
	}
}
