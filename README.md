# A simple NFA-based regular expressions implementation.

This implementation is inspired by the Russ Cox's article
[Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html).

## What can it do?

Not much, this is a very rudimentary implementation done for purely
educational purposes.

You can parse and compile a very basic version of regular
expressions to NFAs and use those NFAs to check if a given string
matches a given expression.

## Supported syntax

There are characters that have special meaning: `(`, `)`, `?`, `*`, `|`

1. UTF8 runes that have no special meaning match themselves.
For example `a` matches the character `a`.

1. You can use quantifiers: `?` (one or zero) and `*` (zero or more).
For example `ab*` matches `a`, `ab`, `abb`, ... .

1. You can concatenate expressions by writing them one after another.
If you take two expressions `E1` and `E2` and write them together: `E1E2`
the resulting expression will match any string that is a concatenation of a string
that matches `E1` and a string that matches `E2`.
For example `ab` matches the string `ab`.

1. You can create union expressions using the alternation symbol `|`.
If you take two expressions `E1` and `E2` and join them with `|`: `E1|E2`
the resulting expression will match any string that matches either `E1` or `E2`.
For example `ab|bc` matches strings `ab` and `bc`.

Priorities from higher to lower go as follows: quantifiers -> concatenation -> alternation.
You can use parentheses (`()`) for grouping.

## Example usage

This function takes a regexp string and a list of strings
and returns the array of matching results (or the error if the regexp is malformed).

```
func patternMatch(pat string, in []string) ([]bool, error) {
	ast, err := parse(pat)
	if err != nil {
		return nil, err
	}
	regex := compile(ast)
	res := make([]bool, len(in))
	for i, s := range in {
		res[i] = match(&regex, s)
	}
	return res, nil
}
```
