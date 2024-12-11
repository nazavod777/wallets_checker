package utils

import "strings"

func RemoveHexPrefix(key string) string {
	if strings.HasPrefix(key, "0x") {
		return key[2:] // Убираем первые два символа
	}
	return key
}
