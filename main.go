package main

import (
	"bufio"
	"debank_checker_v3/core"
	"debank_checker_v3/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func inputUser() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func processAccounts(accounts []string, threads int, userAction int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, threads)

	for _, account := range accounts {
		wg.Add(1)
		sem <- struct{}{}

		go func(acc string) {
			defer wg.Done()

			if userAction == 1 {
				core.ParseDebankAccount(acc)
			} else if userAction == 2 {
				core.ParseRabbyAccount(acc)
			}
			<-sem
		}(account)
	}

	wg.Wait()
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("Unexpected Error: %v", r)
		fmt.Println("Press Enter to Exit..")
		_, err := fmt.Scanln()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("WebSite - nazavod.dev\nAntiDrain - antidrain.me\nTG - t.me/n4z4v0d\n\n")
	defer handlePanic()

	err := utils.InitProxies()

	if err != nil {
		log.Panicf("%s", err)
	}

	err = ensureDir("./results")

	if err != nil {
		log.Panicf("Error When Creating Results Directory: %v", err)
	}

	accountsList, err := utils.ReadFileByRows("data/accounts.txt")

	if err != nil {
		log.Panicf("Error When Reading Accounts File: %v", err)
	}

	fmt.Printf("Successfully Loaded %d Accounts // %d Proxies\n\n", len(accountsList), len(utils.Proxies))
	fmt.Print("Threads: ")

	threads, err := strconv.Atoi(inputUser())

	if err != nil {
		log.Panicf("Wrong Threads Number: %v", err)
	}

	fmt.Print("1. Debank Checker\n2. Rabby Checker\nEnter Your Choice: ")

	userAction, err := strconv.Atoi(inputUser())

	if err != nil || (userAction != 1 && userAction != 2) {
		log.Panicf("Wrong User Action Number: %v", err)
	}

	fmt.Println()

	processAccounts(accountsList, threads, userAction)

	fmt.Print("The Work Has Been Successfully Finished..\n\n")
	fmt.Print("Press Enter to Exit..")
	inputUser()
}
