package main

import (
	"debank_checker_v3/core"
	"debank_checker_v3/core/debank"
	debankl2 "debank_checker_v3/core/debankL2"
	"debank_checker_v3/global"
	"debank_checker_v3/inits"
	"debank_checker_v3/util"
	"fmt"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func main() {
	defer inits.HandlePanic()

	logFile := inits.InitLog()
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Errorf("Error closing log file: %s", err)
		}
	}()

	fmt.Println("Website: nazavod.dev")
	fmt.Println("AntiDrain: antidrain.me")
	fmt.Println("Telegram: t.me/n4z4v0d")
	fmt.Println()

	// Init proxies
	if err := util.InitProxies(filepath.Join("config", "proxies.txt")); err != nil {
		log.Panicf("Failed to initialize proxies: %s", err)
	}

	// Init clients
	inits.InitClients()

	// Read config
	if err := util.ReadJSONFile("./config/config.json", &global.ConfigFile); err != nil {
		log.Panicf("Failed to read config file: %s", err)
	}

	// Ensure results dir
	if err := inits.EnsureDir("./results"); err != nil {
		log.Panicf("Failed to create results directory: %v", err)
	}

	// Read accounts
	accounts, err := util.ReadFileByRows("config/accounts.txt")
	if err != nil {
		log.Panicf("Failed to read accounts file: %v", err)
	}

	global.TargetProgress = len(accounts)
	global.CurrentProgress = 1

	fmt.Printf("Successfully loaded %d accounts // %d proxies\n\n", len(accounts), len(util.Proxies))
	fmt.Print("Threads: ")

	threads, err := strconv.Atoi(inits.ReadUserInput())
	if err != nil || threads <= 0 {
		log.Panicf("Invalid number of threads: %v", err)
	}

	fmt.Print("1. Debank Checker\n2. Debank L2 Balance Parser\nEnter your choice: ")

	action, err := strconv.Atoi(inits.ReadUserInput())
	if err != nil || (action < 1 || action > 2) {
		log.Panicf("Invalid action number: %v", err)
	}

	var parser core.AccountParser

	switch action {
	case 1:
		parser = debank.ParseDebank{}
	case 2:
		parser = debankl2.L2ParserDebank{}
	default:
		log.Panicf("Unknown parser selected")
	}

	fmt.Println()
	inits.ProcessAccounts(accounts, threads, parser)

	fmt.Println("The work has been successfully completed.")
	fmt.Println()
	fmt.Print("Press Enter to exit...")
	_ = inits.ReadUserInput()
}
