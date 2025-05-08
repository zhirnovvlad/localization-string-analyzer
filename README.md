# Localization String Analyzer

A Go utility to analyze iOS/macOS `.strings` localization files for duplicate keys and optionally clean them.

## Features

- Detects duplicate keys in localization files
- Identifies whether duplicated keys have the same or different values
- Reports line numbers of all occurrences
- Can automatically create a cleaned version of the file with duplicates removed
- Preserves the original file by saving cleaned output to a separate file
- Supports output to file and custom input files
- Includes additional utility tools for specific tasks

## Why Use This Tool

In iOS/macOS localization files, duplicate keys can cause issues:

- When the same key appears multiple times with the same value, it unnecessarily increases the file size
- When the same key appears with different values, it creates a localization conflict, as the system might use any of the values unpredictably
- Finding and resolving these issues improves the reliability of your application's localization

### Real-World Performance

In our testing with a sample Localizable.strings file:

- Original file contained 1611 unique keys but 1788 total entries (177 duplicates)
- File size was reduced from 139KB to 133KB (4.3% reduction)
- All functionality was preserved while eliminating potential ambiguity

## Usage

This tool requires Go to be installed on your system. To use it, follow these steps:

1. Clone this repository or download the source code
2. Either run the tool directly with Go:

```bash
# Basic usage (looks for Localizable.strings in current directory)
go run main.go

# Specify a different input file
go run main.go -f path/to/your/Localizable.strings

# Save the output to a file
go run main.go -o output.txt

# Create a cleaned file with duplicates removed
go run main.go -clean=cleaned.strings

# For Localizable.strings, suggest a descriptive name for the cleaned file
go run main.go -clean=Localizable-cleaned.strings

# Combine options
go run main.go -f path/to/your/Localizable.strings -o output.txt -clean=cleaned.strings -v
```

Or build and run the binary:

```bash
# Build the binary using the provided script
./build.sh

# Run the compiled binary
./localization-analyzer -clean=cleaned.strings
```

### Command-line Options

- `-f` : Specify the input localization file (default: Localizable.strings)
- `-o` : Write analysis results to the specified output file instead of stdout
- `-clean` : Create a cleaned version of the file at the specified path (must be different from input file)
- `-v` : Verbose mode - show more details in terminal output

## Additional Utility Tools

In addition to the main analyzer, this repository includes two useful utility tools for specific localization tasks:

### 1. Key Counter (count_keys.go)

A simple utility that counts the total number of keys and unique keys in a .strings file.

```bash
# Count keys in the default Localizable.strings file
go run count_keys.go

# Count keys in a specific file
go run count_keys.go -f path/to/your/Localizable.strings
```

Output example:
```
File: Localizable.strings
Total Entries: 1788
Unique Keys: 1611
Duplicate Entries: 177 (9.9%)
```

### 2. Key Checker (check_keys.go)

A utility to check if a specific key exists in a .strings file and displays its value(s).

```bash
# Check if a key exists in the default Localizable.strings file
go run check_keys.go "YourKeyToCheck"

# Check a key in a specific file
go run check_keys.go -f path/to/your/Localizable.strings "YourKeyToCheck"
```

Output examples:

When a key is found once:
```
Key "Cancel" found in Localizable.strings (1 occurrence):
  Line 45: "Cancel"
```

When a key has multiple occurrences:
```
Key "OK" found in Localizable.strings (2 occurrences):
  Line 15: "OK"
  Line 225: "OK"
All occurrences have the same value.
```

When a key has conflicts:
```
Key "Hello World" found in Localizable.strings (2 occurrences):
  Line 10: "Hello World"
  Line 42: "Hola Mundo"
WARNING: Key has different values in different occurrences (localization conflict)!
```

## Sample Output

When duplicate keys with the same value are found:

```
Duplicate keys found: 2
====================
Key: "Cancel" appears 3 times:
  All entries have the same value: "Cancel"
  Found at lines:
    Line 45
    Line 120
    Line 301

Key: "OK" appears 2 times:
  All entries have the same value: "OK"
  Found at lines:
    Line 15
    Line 225
```

When duplicate keys with different values are found (localization conflict):

```
Key: "Hello World" appears 2 times:
  WARNING: Key has different values (localization conflict)!
  Found at lines:
    Line 10: "Hello World"
    Line 42: "Hola Mundo"
```

## Cleaning Behavior

When using the `-clean` option:

1. The tool creates a new file at the specified path with all duplicate keys removed
2. Only the first occurrence of each key is kept in the cleaned file
3. Comments and empty lines are preserved
4. The original input file is never modified
5. A summary shows how many duplicate entries were removed
6. If you try to use the same filename for input and output, the tool will suggest an alternative

## Localization File Format

This tool is designed to work with standard iOS/macOS `.strings` files that follow this format:

```
"key" = "value";
```

Comments (lines starting with `//`) are automatically ignored.

## Building From Source

```bash
git clone https://github.com/zhirnovvlad/localization-string-analyzer.git
cd localization-string-analyzer
./build.sh
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 