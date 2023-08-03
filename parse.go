package main

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

// parsing error
type parsingError struct {
	wanted string
	got    rune
}

func (e parsingError) Error() string {
	if e.got == 0 {
		return "wanted: " + e.wanted + ", got: EOF"
	}
	return "wanted: " + e.wanted + ", got: '" + string(e.got) + "'"
}

// term ::= rune | '(' alternative ')'
func parseTerm(input *iterator) (node, error) {
	if !isTermChar(input.peek()) {
		return nil, parsingError{"a letter or '('", input.peek()}
	}
	if input.peek() == '(' {
		input.next()
		alt, err := parseAlt(input)
		if err != nil {
			return nil, err
		}
		if input.peek() != ')' {
			return nil, parsingError{"')'", input.peek()}
		}
		input.next()
		return alt, nil
	}
	char := char(input.next())
	return &char, nil
}

func isTermChar(c rune) bool {
	switch c {
	case '|', ')', 0, '*', '?':
		return false
	}
	return true
}

// quantified ::= term | term '*' | term '?'
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

// concatenation ::= quantified | quantified concatenation
func parseConcat(input *iterator) (node, error) {
	left, err := parseQuantified(input)
	if err != nil {
		return nil, err
	}
	if isTermChar(input.peek()) {
		right, err := parseConcat(input)
		if err != nil {
			return nil, err
		}
		left = &concatenation{left: left, right: right}
	}
	return left, nil
}

// alternative ::= concatenation | concatenation '|' alternative
func parseAlt(input *iterator) (node, error) {
	left, err := parseConcat(input)
	if err != nil {
		return nil, err
	}
	if input.peek() == '|' {
		input.next()
		right, err := parseAlt(input)
		if err != nil {
			return nil, err
		}
		left = &alternative{left: left, right: right}
	}
	return left, nil
}

// regex ::= alternative '\0'
func parse(regex string) (node, error) {
	input := &iterator{chars: []rune(regex), pos: 0}
	res, err := parseAlt(input)
	if err != nil {
		return nil, err
	}
	if input.peek() != 0 {
		return nil, parsingError{"EOF", input.peek()}
	}
	return res, nil
}
