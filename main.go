package main

import (
	"os"
	"fmt"
	"bufio"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: match REGEX")
		os.Exit(1)
	}
	ast, err := parse(os.Args[1])
	if err != nil {
		fmt.Println("error: regex parsing failed:", err)
		os.Exit(1)
	}
	regex := compile(ast)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if match(&regex, scanner.Text()) {
			fmt.Println(scanner.Text())
		}
	}
}
