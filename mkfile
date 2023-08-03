main: main.go parse.go compile.go match.go
	go build -o $target $prereq

clean:V:
	rm main
