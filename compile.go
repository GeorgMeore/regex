package main

type state struct {
	c    rune   // required character (0 for no character)
	next *state
	alt  *state // used for branching
}

type nfa struct {
	in  *state
	out *state
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
	t := compile(s.term)
	sout := &state{c: 0, next: nil}
	sin := &state{c: 0, next: t.in, alt: sout}
	t.out.next = sin
	return nfa{in: sin, out: sout}
}

// -[0]-->--[0]->
//   `>--t---^
func compileOpt(o *optional) nfa {
	t := compile(o.term)
	oout := &state{c: 0, next: nil}
	oin := &state{c: 0, next: t.in, alt: oout}
	t.out.next = oout
	return nfa{in: oin, out: oout}
}

// -> l -> r ->
func compileConcat(c *concatenation) nfa {
	l, r := compile(c.left), compile(c.right)
	l.out.next = r.in
	return nfa{in: l.in, out: r.out}
}

//   .->--r-->.
// -[0]      [0]->
//   `->--l-->'
func compileAlt(a *alternative) nfa {
	l, r := compile(a.left), compile(a.right)
	aout := &state{c: 0, next: nil}
	l.out.next = aout
	r.out.next = aout
	ain := &state{c: 0, next: l.in, alt: r.in}
	return nfa{in: ain, out: aout}
}

func compile(regex node) nfa {
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
