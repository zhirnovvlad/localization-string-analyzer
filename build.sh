#!/bin/bash

# Build script for Localization String Analyzer

echo "Building Localization String Analyzer..."

# Make sure Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Build the binary
go build -o localization-analyzer main.go

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created as 'localization-analyzer'"
    echo ""
    echo "Usage:"
    echo "  ./localization-analyzer                  # Analyze Localizable.strings in current directory"
    echo "  ./localization-analyzer -f FILE          # Analyze specified .strings file"
    echo "  ./localization-analyzer -o OUTPUT        # Save analysis to output file"
    echo "  ./localization-analyzer -clean=NEWFILE   # Create cleaned version with duplicates removed"
    echo "  ./localization-analyzer -v               # Verbose output mode"
    echo ""
    echo "Example with multiple options:"
    echo "  ./localization-analyzer -f input.strings -o report.txt -clean=cleaned.strings -v"
    echo ""
    echo "Note: The cleaned file must be different from the input file."
else
    echo "Build failed."
    exit 1
fi 