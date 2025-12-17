package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (can remove later)
var _ = fmt.Fprint
var _ = os.Stdout

var builtinCommands []string = []string{"exit", "echo", "type", "pwd", "cd"}

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		buffer := make([]byte, 1024)
		numBytesRead, err := os.Stdin.Read(buffer)

		if !(err == nil || err == io.EOF) {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return
		}

		bufferedString := string(buffer[:numBytesRead])
		argsSplicedFromBuffer := ParseStdinBufferString(bufferedString)
		if len(argsSplicedFromBuffer) == 0 {
			// No command entered, continue to next loop iteration
			continue
		}
		commandString, inputArgs := argsSplicedFromBuffer[0], argsSplicedFromBuffer[1:]

		switch commandString {

		case "exit":
			return

		case "echo":
			fmt.Println(strings.Join(inputArgs, " "))

		case "pwd":
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				continue
			}
			fmt.Println(currentDir)

		case "cd":
			targetDir := inputArgs[0]

			// Special case: if the targetDir is "~", change to HOME directory
			if targetDir == "~" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
					continue
				}
				targetDir = homeDir
			}

			isDir, err := IsDirectory(targetDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking directory: %v\n", err)
				continue
			} else if !isDir {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", targetDir)
				continue
			}

			err = os.Chdir(targetDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error changing directory: %v\n", err)
				continue
			}

		case "type":
			firstArgument := inputArgs[0]
			if slices.Contains(builtinCommands, firstArgument) {
				fmt.Printf("%s is a shell builtin\n", firstArgument)
			} else {
				executablePath := FindPathOfGivenExecutable(firstArgument)
				if executablePath != "" {
					fmt.Printf("%s is %s\n", firstArgument, executablePath)
				} else {
					fmt.Printf("%s not found\n", firstArgument)
				}
			}

		default:
			executablePath := FindPathOfGivenExecutable(commandString)
			if executablePath == "" {
				fmt.Printf("%s: command not found\n", commandString)
				continue
			}
			cmd := exec.Command(commandString, inputArgs...)
			output, _ := cmd.Output()
			fmt.Print(string(output))
		}
	}
}
