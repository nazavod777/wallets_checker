package utils

import (
	"os"
)

func AppendFile(filePath string, fileContent string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString(fileContent); err != nil {
		panic(err)
	}
}
