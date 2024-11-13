package utils

import (
	"debank_checker_v3/customTypes"
	"fmt"
	"strings"
)

func getKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func FormatResult(accountData string,
	accountAddress string,
	totalUsdBalance float64,
	tokenBalances map[string][]customTypes.TokenBalancesResultData,
	nftBalances map[string][]customTypes.NftBalancesResultData,
	poolsData map[string]map[string][]customTypes.PoolBalancesResultData) {
	var formattedResult string

	formattedResult += fmt.Sprintf("==================== Account Data: %s\n", accountData)
	formattedResult += fmt.Sprintf("==================== Address: %s\n", accountAddress)
	formattedResult += fmt.Sprintf("==================== USD Balance: %f $\n", totalUsdBalance)

	if ConfigFile.DebankConfig.ParseTokens == true && len(tokenBalances) > 0 {
		formattedResult += fmt.Sprintf("\n========== Token Balances\n")

		for _, chainName := range getKeys(tokenBalances) {
			if len(tokenBalances[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("===== %s\n", strings.ToUpper(chainName))

			for _, tokenData := range tokenBalances[chainName] {
				formattedResult += fmt.Sprintf("Name: %s | Balance (in usd): %s $ | Amount: %s\n", tokenData.Name, tokenData.BalanceUSD.Text('f', -1), tokenData.Amount.Text('f', -1))
			}
		}
	} else {
		formattedResult += "\n"
	}

	if ConfigFile.DebankConfig.ParseNfts == true && len(nftBalances) > 0 {
		formattedResult += fmt.Sprintf("\n========== NFT Balances\n")

		for _, chainName := range getKeys(nftBalances) {
			if len(nftBalances[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("===== %s\n", strings.ToUpper(chainName))

			for _, nftData := range nftBalances[chainName] {
				formattedResult += fmt.Sprintf("Name: %s | Price (in usd): %s $ | Amount: %s\n", nftData.Name, nftData.BalanceUSD.Text('f', -1), nftData.Amount.Text('f', -1))
			}
		}
	} else {
		formattedResult += "\n"
	}

	if ConfigFile.DebankConfig.ParsePools == true && len(poolsData) > 0 {
		formattedResult += fmt.Sprintf("\n========== Pool Balances\n")

		for _, chainName := range getKeys(poolsData) {
			if len(poolsData[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("===== %s\n", strings.ToUpper(chainName))

			for _, poolName := range getKeys(poolsData[chainName]) {
				if len(poolsData[chainName][poolName]) < 1 {
					continue
				}

				formattedResult += fmt.Sprintf("===== %s\n", strings.ToUpper(poolName))

				for _, poolData := range poolsData[chainName][poolName] {
					formattedResult += fmt.Sprintf("Name: %s | Balance (in usd): %s $ | Amount: %s\n", poolData.Name, poolData.BalanceUSD.Text('f', -1), poolData.Amount.Text('f', -1))
				}
			}
		}
	}

	formattedResult += "\n\n\n"

	var filePath string

	switch {
	case totalUsdBalance >= 0 && totalUsdBalance < 1:
		filePath = "0_1_debank.txt"

	case totalUsdBalance >= 1 && totalUsdBalance < 10:
		filePath = "1_10_debank.txt"

	case totalUsdBalance >= 10 && totalUsdBalance < 100:
		filePath = "10_100_debank.txt"

	case totalUsdBalance >= 100 && totalUsdBalance < 500:
		filePath = "100_500_debank.txt"

	case totalUsdBalance >= 500 && totalUsdBalance < 1000:
		filePath = "500_1000_debank.txt"

	default:
		filePath = "1000_plus_debank.txt"
	}

	AppendFile("./results/"+filePath,
		formattedResult)
}
