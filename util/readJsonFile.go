package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// ReadJSONFile reads a JSON file into the provided struct reference.
// It returns an error if the file cannot be opened, read, or parsed.
func ReadJSONFile(fileName string, out interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", fileName, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Warnf("failed to close file %q: %v", fileName, cerr)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", fileName, err)
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to parse JSON from file %q: %w", fileName, err)
	}

	return nil
}
