package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("reading line: %v\n", err)
		return
	}
	line = strings.TrimSpace(line)

	fmt.Printf("%s: command not found\n", line)
}
