package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// AppendToFile appends the given content to the specified file path.
// If the file does not exist, it will be created.
func AppendToFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("failed to open file %q for appending: %v", filePath, err)
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Warnf("failed to close file %q: %v", filePath, cerr)
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		log.Errorf("failed to write to file %q: %v", filePath, err)
	}
}
