package util

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// ReadFileByRows reads a file line by line and returns all lines as a slice of strings.
// Returns an error if the file cannot be opened or read.
func ReadFileByRows(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Warnf("failed to close file %q: %v", filename, cerr)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file %q: %w", filename, err)
	}

	return lines, nil
}
