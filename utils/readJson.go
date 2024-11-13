package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJson(fileName string, fileStruct interface{}) error {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	err = json.Unmarshal(byteValue, &fileStruct)
	if err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return nil
}
