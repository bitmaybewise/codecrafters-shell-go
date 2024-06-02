package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		value := waitForUserInput()
		value = strings.Trim(value, "\n")
		cmd := strings.Split(value, " ")
		shellFn, ok := shellBuiltins[cmd[0]]
		if !ok {
			fmt.Printf("%s: command not found\n", cmd[0])
			continue
		}
		shellFn(cmd, shellBuiltins)
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
	osPath := os.Getenv("PATH")
	paths := strings.Split(osPath, ":")
	for _, dir := range paths {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == cmd[1] {
				fmt.Printf("%s is %s/%s\n", cmd[1], dir, cmd[1])
				return
			}
		}
	}

	fmt.Printf("%s not found\n", cmd[1])
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
