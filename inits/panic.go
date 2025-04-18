package inits

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func HandlePanic() {
	if r := recover(); r != nil {
		log.Errorf("Unexpected Error: %v", r)
		fmt.Println("Press Enter to Exit..")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}
}
