package main

import (
	"os"
	"path/filepath"
	"strings"
)

// Checks if any execute permissions (by Owner, Group, or Others) are set on the given file mode
func IsExecByAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

// Checks if the given string is a path to an existing directory.
// Returns (false, nil) if the path does not exist.
func IsDirectory(dirPath string) (bool, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
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

// Parses the given stdin buffer string into a slice of arguments.
// Handles single quotes to allow spaces within arguments.
func ParseStdinBufferString(s string) []string {
	var result []string
	var current strings.Builder
	inQuotes, hasContent := false, false

	for _, char := range s {
		switch {
		case char == '\'':
			inQuotes = !inQuotes // Toggle mode
			hasContent = true    // Even an empty quote '' counts as content
		case inQuotes:
			current.WriteRune(char)
			hasContent = true
		case isWhitespace(char):
			if hasContent {
				result = append(result, current.String())
				current.Reset()
				hasContent = false
			}
		default:
			current.WriteRune(char)
			hasContent = true
		}
	}

	if hasContent {
		result = append(result, current.String())
	}
	return result
}

// Checks if the given rune is a (standard) whitespace character
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
