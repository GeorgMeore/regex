main: main.go compile.go syntax.go
	go build -o $target $prereq

clean:V:
	rm main
