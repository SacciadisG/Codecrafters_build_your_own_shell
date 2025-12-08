package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {

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
	commandString := strings.TrimSpace(bufferedString)
	switch commandString {
	default:
		fmt.Printf("%s: command not found\n", commandString)
	}
}
