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
	if regex.lastid == 0 {
		// handle overflow (highly unlikely, just for the sake of correctness)
		resetRegex(regex)
	}
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

func resetRegex(regex *nfa) {
	regex.lastid = 1
	resetStates(regex.in)
}

func resetStates(s *state) {
	// if s.id is 0 then s was either already reset,
	// or it was never visited during matching in the first place
	if s == nil || s.id == 0 {
		return
	}
	s.id = 0
	resetStates(s.next)
	resetStates(s.alt)
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
