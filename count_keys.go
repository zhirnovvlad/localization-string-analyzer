package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Simple utility to count the number of unique keys in a .strings file
func main() {
	// Parse command-line flags
	var inputFile string
	flag.StringVar(&inputFile, "f", "Localizable.strings", "Input localization file (default: Localizable.strings)")
	flag.Parse()

	// Check if the file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist\n", inputFile)
		os.Exit(1)
	}

	// Count unique keys
	keyCount, totalEntries, err := countKeys(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Report results
	fmt.Printf("File: %s\n", inputFile)
	fmt.Printf("Total Entries: %d\n", totalEntries)
	fmt.Printf("Unique Keys: %d\n", keyCount)

	if totalEntries > keyCount {
		duplicates := totalEntries - keyCount
		duplicatePercentage := float64(duplicates) / float64(totalEntries) * 100
		fmt.Printf("Duplicate Entries: %d (%.1f%%)\n", duplicates, duplicatePercentage)
	} else {
		fmt.Println("No duplicate keys found.")
	}
}

func countKeys(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Map to track unique keys
	uniqueKeys := make(map[string]bool)

	// Regular expression to extract key-value pairs
	// This pattern matches: "key" = "value";
	kvPattern := regexp.MustCompile(`"([^"]+)"\s*=\s*"([^"]+)"\s*;`)

	scanner := bufio.NewScanner(file)
	totalEntries := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comment lines or empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		matches := kvPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			uniqueKeys[key] = true
			totalEntries++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("error scanning file: %w", err)
	}

	return len(uniqueKeys), totalEntries, nil
}
