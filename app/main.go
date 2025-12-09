package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var builtinCommands []string = []string{"exit", "echo", "type"}

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		buffer := make([]byte, 1024)
		numBytesRead, err := os.Stdin.Read(buffer)

		if !(err == nil || err == io.EOF) {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return
		}

		// Switch on the command string (trimming the newline & return characters)
		// For now, all commands are considered as 'invalid'
		bufferedString := string(buffer[:numBytesRead])
		argsSplicedFromBuffer := strings.Fields(bufferedString)
		commandString, inputArgs := argsSplicedFromBuffer[0], argsSplicedFromBuffer[1:]

		switch commandString {
		case "exit":
			return
		case "echo":
			fmt.Println(strings.Join(inputArgs, " "))
		case "type":
			if len(inputArgs) > 2 || (!slices.Contains(builtinCommands, inputArgs[0])) {
				fmt.Printf("%s: not found\n", inputArgs[0])
			} else {
				fmt.Printf("%s is a shell builtin\n", inputArgs[0])
			}
		default:
			fmt.Printf("%s: command not found\n", commandString)
		}
	}
}
