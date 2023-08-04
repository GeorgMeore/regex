package main

func build(regex string) (nfa, error) {
	ast, err := parse(regex)
	if err != nil {
		return nfa{}, err
	}
	return compile(ast), nil
}

func match(regex *nfa, input string) bool {
	states := []*state{}
	regex.lastid += 1
	addState(&states, regex.in, regex.lastid)
	for _, c := range input {
		regex.lastid += 1
		states = step(states, c, regex.lastid)
		if len(states) == 0 {
			return false
		}
	}
	for _, s := range states {
		if s == regex.out {
			return true
		}
	}
	return false
}

func addState(l *[]*state, s *state, id int) {
	if s == nil || s.id == id {
		return
	}
	s.id = id
	if s.c == 0 {
		addState(l, s.next, id)
		addState(l, s.alt, id)
		return
	}
	*l = append(*l, s)
}

func step(states []*state, c rune, id int) []*state {
	next := []*state{}
	for _, s := range states {
		if s.c == 0 || s.c == c {
			addState(&next, s.next, id)
			addState(&next, s.alt, id)
		}
	}
	return next
}
