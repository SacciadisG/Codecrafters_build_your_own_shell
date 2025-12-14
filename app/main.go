package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var builtinCommands []string = []string{"exit", "echo", "type"}

// Checks if any execute permissions (by Owner, Group, or Others) are set on the given file mode
func IsExecByAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

// Returns the full path of the given executable if found anywhere in the PATH environment variable.
// If not found, returns an empty string.
func FindPathOfGivenExecutable(executableName string) string {
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv) // Splitting is done in an OS-agnostic way

	for _, dirPath := range paths {
		fileFullPath := filepath.Join(dirPath, executableName)
		fileInfo, err := os.Stat(fileFullPath)
		if err != nil {
			// File doesn't exist in this dir, continue to next dir
			continue
		}
		if IsExecByAny(fileInfo.Mode()) {
			return fileFullPath
		} else {
			// File exists but isn't executable, continue to next dir
			continue
		}
	}
	return ""
}

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

		// TODO: Assuming that inputArgs always has at least one element for now. Handle edge cases later.

		switch commandString {

		case "exit":
			return

		case "echo":
			fmt.Println(strings.Join(inputArgs, " "))

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
			fmt.Println(string(output))
		}
	}
}
