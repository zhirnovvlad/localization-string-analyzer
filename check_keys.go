package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// This utility checks if a key exists in a localization file and returns its value
// Useful for quickly verifying the presence and value of a specific key
func main() {
	// Parse command-line flags
	var inputFile string
	flag.StringVar(&inputFile, "f", "Localizable.strings", "Input localization file (default: Localizable.strings)")
	flag.Parse()

	// Get the key to check
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No key specified")
		fmt.Println("Usage: go run check_keys.go [-f filename.strings] \"key_to_check\"")
		os.Exit(1)
	}

	keyToCheck := args[0]

	// Check if the file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist\n", inputFile)
		os.Exit(1)
	}

	// Look for the key
	occurrences, err := findKeyOccurrences(inputFile, keyToCheck)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Report findings
	if len(occurrences) == 0 {
		fmt.Printf("Key \"%s\" not found in %s\n", keyToCheck, inputFile)
	} else {
		fmt.Printf("Key \"%s\" found in %s (%d occurrences):\n", keyToCheck, inputFile, len(occurrences))

		for _, occurrence := range occurrences {
			fmt.Printf("  Line %d: \"%s\"\n", occurrence.LineNum, occurrence.Value)
		}

		if len(occurrences) > 1 {
			// Check if all values are the same
			allSame := true
			firstValue := occurrences[0].Value
			for _, occ := range occurrences[1:] {
				if occ.Value != firstValue {
					allSame = false
					break
				}
			}

			if allSame {
				fmt.Println("All occurrences have the same value.")
			} else {
				fmt.Println("WARNING: Key has different values in different occurrences (localization conflict)!")
			}
		}
	}
}

type KeyOccurrence struct {
	Value   string
	LineNum int
}

func findKeyOccurrences(filename, keyToFind string) ([]KeyOccurrence, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var occurrences []KeyOccurrence

	// Regular expression to extract key-value pairs
	// This pattern matches: "key" = "value";
	kvPattern := regexp.MustCompile(`"([^"]+)"\s*=\s*"([^"]+)"\s*;`)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comment lines or empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		matches := kvPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			value := matches[2]

			if key == keyToFind {
				occurrences = append(occurrences, KeyOccurrence{
					Value:   value,
					LineNum: lineNum,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return occurrences, nil
}
