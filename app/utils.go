package main

import (
	"os"
	"path/filepath"
	"slices"
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

// ParseStdinBufferString splits a string into arguments using the following rules:
//  1. Whitespace is a delimiter only outside of quotes.
//  2. Inside ” or "", all characters (including spaces) are treated literally.
//  3. Quotes of one type lose their special meaning inside quotes of the other type.
//  4. Adjacent segments (quoted or unquoted) are concatenated into a single argument.
func ParseStdinBufferString(s string) []string {
	var result []string
	var current strings.Builder
	var activeQuote rune                     // Tracks ' or ". 0 means no active quote
	escapedSpecialChars := []rune{'"', '\\'} // Characters that can be escaped within double quotes
	hasContent := false

	for i := 0; i < len(s); i++ {
		char := rune(s[i])
		switch {
		// In quote blocks
		case activeQuote != 0:
			switch {
			case char == '\\' && activeQuote == '"':
				if i+1 < len(s) && slices.Contains(escapedSpecialChars, rune(s[i+1])) {
					current.WriteRune(rune(s[i+1]))
					i++ // Skip the next character as it's been escaped
				} else {
					current.WriteRune(char)
				}

			case char == activeQuote:
				activeQuote = 0   // Close the quote block
				hasContent = true // Mark that we have an argument in progress (concatenation)

			default:
				current.WriteRune(char)
			}

		// Outside quote blocks
		case char == '\\':
			if i+1 < len(s) {
				current.WriteRune(rune(s[i+1]))
				i++ // Skip the next character as it's been escaped
			} else {
				current.WriteRune(char)
			}
			hasContent = true

		case char == '\'' || char == '"':
			activeQuote = char
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
