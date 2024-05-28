package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Exit = "exit"
	Echo = "echo"
	Type = "type"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("reading line: %v\n", err)
			return
		}

		resp := evaluate(strings.TrimSpace(line))
		fmt.Println(resp)
	}
}

func evaluate(line string) string {
	switch {
	case strings.HasPrefix(line, Exit):
		status := 0
		if p := strings.Split(line, " "); len(p) == 2 {
			var err error
			status, err = strconv.Atoi(p[1])
			if err != nil || status < 0 || status > 255 {
				return fmt.Sprintf("invalid exit status: %s", p[1])
			}
		}
		os.Exit(status)

	case strings.HasPrefix(line, Echo):
		return strings.TrimLeft(line, Echo+" ")

	case strings.HasPrefix(line, Type):
		p := strings.Split(line, " ")
		if len(p) != 2 {
			return fmt.Sprintf("invalid number of arguments")
		}

		if isBuiltIn(p[1]) {
			return fmt.Sprintf("%s is a shell builtin", p[1])
		}

		path, err := isInPath(p[1])
		if err != nil {
			return err.Error()
		}
		if path != "" {
			return path
		}

		return fmt.Sprintf("%s: command not found", p[1])

	default:
		return fmt.Sprintf("%s: command not found", line)
	}
	return ""
}

func isBuiltIn(cmd string) bool {
	switch cmd {
	case "exit":
		return true
	case "echo":
		return true
	case "type":
		return true
	default:
		return false
	}
}

func isInPath(cmd string) (string, error) {
	path, ok := os.LookupEnv("PATH")
	if !ok {
		return "", fmt.Errorf("PATH not set")
	}

	for _, dir := range strings.Split(path, ":") {
		contents, err := os.ReadDir(dir)
		if err != nil {
			return "", fmt.Errorf("reading path directory: %w", err)
		}

		for _, file := range contents {
			if file.IsDir() {
				continue
			}
			if file.Name() == cmd {
				return fmt.Sprintf("%s is %s", cmd, filepath.Join(dir, file.Name())), nil
			}
		}
	}
	return "", nil
}
