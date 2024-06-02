package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		value := waitForUserInput()
		value = strings.Trim(value, "\n")
		cmd := strings.Split(value, " ")
		shellFn, ok := shellBuiltins[cmd[0]]
		if ok {
			shellFn(cmd, shellBuiltins)
			continue
		}
		shellExec(cmd)
	}
}

type shellFunction func([]string, shellBuiltinsType)

type shellBuiltinsType map[string]shellFunction

var shellBuiltins = shellBuiltinsType{
	"exit": checkArgs(shellExit, -1),
	"echo": checkArgs(shellEcho, -1),
	"type": checkArgs(shellType, -1),
}

func shellEcho(cmd []string, _builtins shellBuiltinsType) {
	fmt.Println(strings.Join(cmd[1:], " "))
}

func shellExit(_cmd []string, _builtins shellBuiltinsType) {
	os.Exit(0)
}

func shellType(cmd []string, builtins shellBuiltinsType) {
	_, ok := builtins[cmd[1]]
	if ok {
		fmt.Printf("%s is a shell builtin\n", cmd[1])
		return
	}

	dir, err := findInPath(cmd[1])
	if err != nil {
		fmt.Printf("%s %s\n", cmd[1], err)
		return
	}

	fmt.Printf("%s is %s/%s\n", cmd[1], dir, cmd[1])
}

func shellExec(cmd []string) {
	_exec := func(executable string, args ...string) (string, error) {
		_cmd := exec.Command(executable, args...)
		output, err := _cmd.Output()
		if err != nil {
			return "", fmt.Errorf("%s: error=%q\n", cmd[0], err)
		}
		return string(output), nil
	}

	output, err := _exec(cmd[0], cmd[1:]...)
	if err == nil {
		fmt.Print(output)
		return
	} // ignore when exe not found

	dir, err := findInPath(cmd[0])
	if err != nil {
		fmt.Printf("%s: command %s\n", cmd[0], err)
		return
	}
	executable := fmt.Sprintf("%s/%s", dir, cmd[0])
	args := cmd[1:]

	output, err = _exec(executable, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", output)
}

func waitForUserInput() string {
	value, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("command error %q", err)
	}
	return value
}

func checkArgs(fn shellFunction, _n int) shellFunction {
	return func(cmd []string, builtins shellBuiltinsType) {
		// TO DO
		fn(cmd, builtins)
	}
}

func findInPath(cmd string) (string, error) {
	osPath := os.Getenv("PATH")
	paths := strings.Split(osPath, ":")

	for _, dir := range paths {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == cmd {
				return dir, nil
			}
		}
	}

	return "", errors.New("not found")
}
