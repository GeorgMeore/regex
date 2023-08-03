match: main.go parse.go compile.go match.go
	go build -o $target $prereq

test-%:V: %.go %_test.go
	go test $prereq

test-match: compile.go parse.go

test:V: test-parse test-match

clean:V:
	rm -f match
