package main

import (
	"errors"
	"fmt"
)

// input iterator
type iterator struct {
	chars []rune
	pos   int
}

func (i *iterator) peek() rune {
	if i.pos == len(i.chars) {
		return 0
	}
	return i.chars[i.pos]
}

func (i *iterator) next() rune {
	if i.pos == len(i.chars) {
		return 0
	}
	i.pos += 1
	return i.chars[i.pos-1]
}

// ast node
type node interface {
	astNode()
	String() string
}

type char rune

func (*char) astNode() {}
func (c *char) String() string {
	return string(*c)
}

type star struct {
	term node
}

func (*star) astNode() {}
func (s *star) String() string {
	return "star{" + s.term.String() + "}"
}

type optional struct {
	term node
}

func (*optional) astNode() {}
func (o *optional) String() string {
	return "opt{" + o.term.String() + "}"
}

type concatenation struct {
	left  node
	right node
}

func (*concatenation) astNode() {}
func (c *concatenation) String() string {
	return "concat{" + c.left.String() + "," + c.right.String() + "}"
}

type alternative struct {
	left  node
	right node
}

func (*alternative) astNode() {}
func (a *alternative) String() string {
	return "alt{" + a.left.String() + "," + a.right.String() + "}"
}

// TERM ::= RUNE | '(' ALTERNATIVE ')'
func parseTerm(input *iterator) (node, error) {
	if input.peek() == 0 {
		return nil, errors.New("unexpected EOF")
	}
	if input.peek() == '(' {
		input.next()
		alt, err := parseAlternative(input)
		if err != nil {
			return nil, err
		}
		if input.next() != ')' {
			return nil, errors.New("')' expected")
		}
		return alt, nil
	}
	char := char(input.next())
	return &char, nil
}

// QUANTIFIED ::= TERM | TERM '*' | TERM '?'
func parseQuantified(input *iterator) (node, error) {
	term, err := parseTerm(input)
	if err != nil {
		return nil, err
	}
	if input.peek() == '*' {
		input.next()
		return &star{term}, nil
	}
	if input.peek() == '?' {
		input.next()
		return &optional{term}, nil
	}
	return term, nil
}

// CONCATENATION ::= QUANTIFIED | QUANTIFIED CONCATENATION
func parseConcatenation(input *iterator) (node, error) {
	left, err := parseQuantified(input)
	if err != nil {
		return nil, err
	}
	if input.peek() != '|' && input.peek() != ')' && input.peek() != 0 {
		right, err := parseConcatenation(input)
		if err != nil {
			return nil, err
		}
		left = &concatenation{left: left, right: right}
	}
	return left, nil
}

// ALTERNATIVE ::= CONCATENATION | CONCATENATION '|' ALTERNATIVE
func parseAlternative(input *iterator) (node, error) {
	left, err := parseConcatenation(input)
	if err != nil {
		return nil, err
	}
	if input.peek() == '|' {
		input.next()
		right, err := parseAlternative(input)
		if err != nil {
			return nil, err
		}
		left = &alternative{left: left, right: right}
	}
	return left, nil
}

func parse(regex string) (node, error) {
	return parseAlternative(&iterator{chars: []rune(regex), pos: 0})
}

func main() {
	a, e := parse("a(b|c)*d?")
	fmt.Println(a, e)
	a, e = parse("a(b?|c?).d")
	fmt.Println(a, e)
}