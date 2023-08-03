package main

type state struct {
	c    rune   // required character (0 for no character)
	next *state
	alt  *state
	id   int    // marks visited states during execution
}

type nfa struct {
	in     *state
	out    *state
	lastid int // last mark id
}

// -[c]->
func compileChar(c *char) nfa {
	s := &state{c: rune(*c), next: nil}
	return nfa{in: s, out: s}
}

// -[0]->-[0]->
//   |`->--.
//   `<-t--'
func compileStar(s *star) nfa {
	t := compileNode(s.term)
	sout := &state{c: 0, next: nil}
	sin := &state{c: 0, next: t.in, alt: sout}
	t.out.next = sin
	return nfa{in: sin, out: sout}
}

// -[0]-->--[0]->
//   `>--t---^
func compileOpt(o *optional) nfa {
	t := compileNode(o.term)
	oout := &state{c: 0, next: nil}
	oin := &state{c: 0, next: t.in, alt: oout}
	t.out.next = oout
	return nfa{in: oin, out: oout}
}

// -> l -> r ->
func compileConcat(c *concatenation) nfa {
	l, r := compileNode(c.left), compileNode(c.right)
	l.out.next = r.in
	return nfa{in: l.in, out: r.out}
}

//   .->--r-->.
// -[0]      [0]->
//   `->--l-->'
func compileAlt(a *alternative) nfa {
	l, r := compileNode(a.left), compileNode(a.right)
	aout := &state{c: 0, next: nil}
	l.out.next = aout
	r.out.next = aout
	ain := &state{c: 0, next: l.in, alt: r.in}
	return nfa{in: ain, out: aout}
}

func compileNode(regex node) nfa {
	switch regex.(type) {
	case *char:
		return compileChar(regex.(*char))
	case *star:
		return compileStar(regex.(*star))
	case *optional:
		return compileOpt(regex.(*optional))
	case *concatenation:
		return compileConcat(regex.(*concatenation))
	case *alternative:
		return compileAlt(regex.(*alternative))
	default:
		panic("unhandled expression type")
	}
}

func compile(regex node) nfa {
	r := compileNode(regex)
	success := &state{c: -1}
	r.out.next = success
	return nfa{in: r.in, out: success}
}
