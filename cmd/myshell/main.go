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
		if cmd[0] == "exit" {
			os.Exit(0)
		}
		fmt.Printf("%s: command not found\n", cmd[0])
	}
}

func waitForUserInput() string {
	value, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("command error %q", err)
	}
	return value
}
