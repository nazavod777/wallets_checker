package inits

import (
	"debank_checker_v3/core"
	"debank_checker_v3/global"
	"sync"
	"sync/atomic"
)

func ProcessAccounts(accounts []string, threads int, parser core.AccountParser) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, threads)

	for _, account := range accounts {
		wg.Add(1)
		sem <- struct{}{}

		go func(acc string) {
			defer wg.Done()
			defer func() { <-sem }()

			parser.Parse(acc)
			atomic.AddInt32(&global.CurrentProgress, 1)
		}(account)
	}

	wg.Wait()
}
