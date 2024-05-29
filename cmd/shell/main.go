package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Exit = "exit"
	Echo = "echo"
	Type = "type"
	CD   = "cd"
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
		if resp != "" {
			fmt.Println(resp)
		}
	}
}

func evaluate(line string) string {
	if line == "" {
		return ""
	}
	parts := strings.Split(line, " ")
	app, err := isInPath(parts[0])
	if err != nil {
		return err.Error()
	}
	if app != "" {
		command := exec.Command(app, parts[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return err.Error()
		}
		return ""
	}

	switch parts[0] {
	case Exit:
		status := 0
		if len(parts) == 2 {
			status, err = strconv.Atoi(parts[1])
			if err != nil || status < 0 || status > 255 {
				return fmt.Sprintf("invalid exit status: %s", parts[1])
			}
		}
		os.Exit(status)

	case Echo:
		echo, _ := strings.CutPrefix(line, Echo+" ")
		return echo

	case Type:
		if len(parts) != 2 {
			return fmt.Sprintf("invalid number of arguments")
		}

		if isBuiltIn(parts[1]) {
			return fmt.Sprintf("%s is a shell builtin", parts[1])
		}

		path, err := isInPath(parts[1])
		if err != nil {
			return err.Error()
		}
		if path != "" {
			return fmt.Sprintf("%s is %s", parts[1], path)
		}

		return fmt.Sprintf("%s: not found", parts[1])

	case CD:
		if len(parts) != 2 {
			return fmt.Sprintf("invalid number of arguments")
		}
		if err = os.Chdir(parts[1]); err != nil {
			return fmt.Sprintf("%s: No such file or directory", parts[1])
		}
		return ""
	}

	return fmt.Sprintf("%s: command not found", line)
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
			continue
		}

		for _, file := range contents {
			if file.IsDir() {
				continue
			}
			if file.Name() == cmd {
				return filepath.Join(dir, file.Name()), nil
			}
		}
	}
	return "", nil
}
