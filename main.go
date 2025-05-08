package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type KeyValue struct {
	Key     string
	Value   string
	LineNum int
}

func main() {
	// Parse command-line flags
	var outputFile string
	var inputFile string
	var cleanFile string
	var verbose bool

	flag.StringVar(&outputFile, "o", "", "Output file for results (optional)")
	flag.StringVar(&inputFile, "f", "Localizable.strings", "Input localization file (default: Localizable.strings)")
	flag.StringVar(&cleanFile, "clean", "", "Create a cleaned version (without duplicates) at the specified path")
	flag.BoolVar(&verbose, "v", false, "Verbose output - include details in terminal output")
	flag.Parse()

	// Set up output
	var output *os.File
	var err error
	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// Analyze the file
	duplicateKeys, uniqueEntries, rawLines, err := analyzeLocalizationFile(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Report duplicate keys
	if len(duplicateKeys) > 0 {
		fmt.Fprintf(output, "Duplicate keys found: %d\n", len(duplicateKeys))
		fmt.Fprintf(output, "====================\n")

		// Sort keys for consistent output
		var keys []string
		for key := range duplicateKeys {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			entries := duplicateKeys[key]
			fmt.Fprintf(output, "Key: \"%s\" appears %d times:\n", key, len(entries))

			// Are all values the same?
			allSame := true
			firstValue := entries[0].Value
			for _, entry := range entries[1:] {
				if entry.Value != firstValue {
					allSame = false
					break
				}
			}

			if allSame {
				fmt.Fprintf(output, "  All entries have the same value: \"%s\"\n", firstValue)
			} else {
				fmt.Fprintf(output, "  WARNING: Key has different values (localization conflict)!\n")
			}

			fmt.Fprintf(output, "  Found at lines:\n")
			for _, entry := range entries {
				if !allSame {
					fmt.Fprintf(output, "    Line %d: \"%s\"\n", entry.LineNum, entry.Value)
				} else {
					fmt.Fprintf(output, "    Line %d\n", entry.LineNum)
				}
			}
			fmt.Fprintf(output, "\n")
		}
	} else {
		fmt.Fprintf(output, "No duplicate keys found.\n")
	}

	// Create a cleaned file if requested
	if cleanFile != "" {
		// Make sure we're not overwriting the input file
		if filepath.Clean(cleanFile) == filepath.Clean(inputFile) {
			// Suggest a different name based on the input file
			suggestedName := createUniqueFilename(inputFile)
			fmt.Printf("Error: Clean file cannot be the same as input file.\n")
			fmt.Printf("Please use a different filename, e.g., '%s'\n", suggestedName)
			os.Exit(1)
		}

		err := createCleanFile(cleanFile, uniqueEntries, rawLines)
		if err != nil {
			fmt.Printf("Error creating clean file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created cleaned file at %s\n", cleanFile)
		fmt.Printf("Removed %d duplicate key entries.\n", countDuplicates(duplicateKeys))
	}

	// Print summary if outputting to file or in verbose mode
	if outputFile != "" || verbose {
		if len(duplicateKeys) > 0 {
			fmt.Printf("Analysis complete. Found %d duplicate keys with %d total duplicated entries.\n",
				len(duplicateKeys), countDuplicates(duplicateKeys))

			if outputFile != "" {
				fmt.Printf("Results written to %s\n", outputFile)
			}

			if cleanFile == "" {
				fmt.Println("Use -clean=filename.strings to create a cleaned version with duplicates removed.")
			}
		} else if verbose {
			fmt.Println("No duplicate keys found.")
		}
	}
}

func createUniqueFilename(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	return filepath.Join(dir, nameWithoutExt+"-cleaned"+ext)
}

func countDuplicates(duplicateKeys map[string][]KeyValue) int {
	count := 0
	for _, entries := range duplicateKeys {
		count += len(entries) - 1 // Count all occurrences beyond the first one
	}
	return count
}

func createCleanFile(filename string, uniqueEntries map[string]KeyValue, rawLines []string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	cleanFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create clean file: %w", err)
	}
	defer cleanFile.Close()

	// First, write all non-key-value lines (comments, empty lines)
	// and the first occurrence of each key
	writtenKeys := make(map[string]bool)

	for _, line := range rawLines {
		trimmedLine := strings.TrimSpace(line)

		// Write comments and empty lines as-is
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			fmt.Fprintln(cleanFile, line)
			continue
		}

		// Extract key if this is a key-value line
		kvPattern := regexp.MustCompile(`"([^"]+)"\s*=\s*"([^"]+)"\s*;`)
		matches := kvPattern.FindStringSubmatch(line)

		if len(matches) == 3 {
			key := matches[1]

			// If we haven't written this key yet, write it
			if !writtenKeys[key] {
				fmt.Fprintln(cleanFile, line)
				writtenKeys[key] = true
			}
			// Otherwise, skip this duplicate
		} else {
			// Write non-matching lines (not key-value format) as-is
			fmt.Fprintln(cleanFile, line)
		}
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func analyzeLocalizationFile(filename string) (map[string][]KeyValue, map[string]KeyValue, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Map to track keys and all their occurrences
	keyEntries := make(map[string][]KeyValue)

	// Map to track duplicate keys (keys with multiple entries)
	duplicateKeys := make(map[string][]KeyValue)

	// Map to store unique entries (first occurrence of each key)
	uniqueEntries := make(map[string]KeyValue)

	// Store all raw lines for recreating the file
	var rawLines []string

	// Regular expression to extract key-value pairs
	// This pattern matches: "key" = "value";
	kvPattern := regexp.MustCompile(`"([^"]+)"\s*=\s*"([^"]+)"\s*;`)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		rawLines = append(rawLines, line)

		// Skip comment lines or empty lines for key analysis
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		matches := kvPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			value := matches[2]

			// Add this entry to keyEntries
			entry := KeyValue{
				Key:     key,
				Value:   value,
				LineNum: lineNum,
			}

			// Store first occurrence in uniqueEntries
			if _, exists := uniqueEntries[key]; !exists {
				uniqueEntries[key] = entry
			}

			keyEntries[key] = append(keyEntries[key], entry)

			// If we now have more than one entry for this key, it's a duplicate
			if len(keyEntries[key]) > 1 {
				duplicateKeys[key] = keyEntries[key]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("error scanning file: %w", err)
	}

	return duplicateKeys, uniqueEntries, rawLines, nil
}
