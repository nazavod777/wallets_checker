package inits

import (
	"bufio"
	"os"
	"strings"
)

func ReadUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
