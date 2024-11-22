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

	formattedResult += fmt.Sprintf("==================== Address: %s (%f $)\n", accountAddress, totalUsdBalance)
	formattedResult += fmt.Sprintf("==================== Account Data: %s\n", accountData)

	if ConfigFile.DebankConfig.ParseTokens == true && len(tokenBalances) > 0 {
		totalTokens := 0

		for _, tokens := range tokenBalances {
			totalTokens += len(tokens)
		}

		formattedResult += fmt.Sprintf("\n========== Token Balances (%d tokens)\n", totalTokens)

		for _, chainName := range getKeys(tokenBalances) {
			if len(tokenBalances[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("----- %s (%d tokens)\n", strings.ToUpper(chainName), len(tokenBalances[chainName]))

			for _, tokenData := range tokenBalances[chainName] {
				formattedResult += fmt.Sprintf("    Name: %s | Balance (in usd): %s $ | Amount: %s | CA: %s\n", tokenData.Name, tokenData.BalanceUSD.Text('f', -1), tokenData.Amount.Text('f', -1), tokenData.ContractAddress)
			}
			formattedResult += "\n"
		}
	} else {
		formattedResult += "\n"
	}

	if ConfigFile.DebankConfig.ParseNfts == true && len(nftBalances) > 0 {
		totalNFTs := 0

		for _, balances := range nftBalances {
			totalNFTs += len(balances)
		}

		formattedResult += fmt.Sprintf("\n========== NFT Balances (%d nfts)\n", totalNFTs)

		for _, chainName := range getKeys(nftBalances) {
			if len(nftBalances[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("----- %s (%d nfts)\n", strings.ToUpper(chainName), len(nftBalances[chainName]))

			for _, nftData := range nftBalances[chainName] {
				formattedResult += fmt.Sprintf("    Name: %s | Price (in usd): %s $ | Amount: %s\n", nftData.Name, nftData.BalanceUSD.Text('f', -1), nftData.Amount.Text('f', -1))
			}
			formattedResult += "\n"
		}
	} else {
		formattedResult += "\n"
	}

	if ConfigFile.DebankConfig.ParsePools == true && len(poolsData) > 0 {
		totalPools := 0

		for _, chainPools := range poolsData {
			for _, pools := range chainPools {
				totalPools += len(pools)
			}
		}

		formattedResult += fmt.Sprintf("\n========== Pool Balances (%d pools)\n", totalPools)

		for _, chainName := range getKeys(poolsData) {
			if len(poolsData[chainName]) < 1 {
				continue
			}

			formattedResult += fmt.Sprintf("----- %s (%d pools)\n", strings.ToUpper(chainName), len(poolsData[chainName]))

			for _, poolName := range getKeys(poolsData[chainName]) {
				if len(poolsData[chainName][poolName]) < 1 {
					continue
				}

				formattedResult += fmt.Sprintf("===== %s\n", strings.ToUpper(poolName))

				for _, poolData := range poolsData[chainName][poolName] {
					formattedResult += fmt.Sprintf("    Name: %s | Balance (in usd): %s $ | Amount: %s\n", poolData.Name, poolData.BalanceUSD.Text('f', -1), poolData.Amount.Text('f', -1))
				}
				formattedResult += "\n"
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
