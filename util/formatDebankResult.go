package util

import (
	"debank_checker_v3/customTypes"
	"debank_checker_v3/global"
	"fmt"
	"strings"
)

// getKeys returns a sorted list of keys from a map.
func getKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// FormatResult generates and writes a formatted string representing the wallet balances by tokens, NFTs, and pools.
func FormatResult(
	accountData string,
	accountAddress string,
	totalUsdBalance float64,
	tokensChainsCount int,
	nftChainsCount int,
	tokenBalances map[string][]customTypes.TokenBalance,
	nftBalances map[string][]customTypes.NFTBalance,
	poolsData map[string]map[string][]customTypes.PoolBalance,
) {
	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("========== Address: %s | Balance: $%.2f | Token Chains: %d | NFT Chains: %d ==========\n", accountAddress, totalUsdBalance, tokensChainsCount, nftChainsCount))
	output.WriteString(fmt.Sprintf("Account Data: %s\n\n", accountData))

	// Token Balances
	if global.ConfigFile.Debank.ParseTokens && len(tokenBalances) > 0 {
		totalTokens := 0
		for _, tokens := range tokenBalances {
			totalTokens += len(tokens)
		}

		output.WriteString(fmt.Sprintf("----- TOKEN BALANCES (%d total tokens) -----\n", totalTokens))

		for _, chain := range getKeys(tokenBalances) {
			tokens := tokenBalances[chain]
			if len(tokens) == 0 {
				continue
			}
			output.WriteString(fmt.Sprintf("  [%s] %d tokens:\n", strings.ToUpper(chain), len(tokens)))

			for _, token := range tokens {
				output.WriteString(fmt.Sprintf("    - %s | $%s | Amount: %s | Contract: %s\n",
					token.Name,
					token.BalanceUSD.Text('f', 2),
					token.Amount.Text('f', 4),
					token.ContractAddress,
				))
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	// NFT Balances
	if global.ConfigFile.Debank.ParseNfts && len(nftBalances) > 0 {
		totalNFTs := 0
		for _, nfts := range nftBalances {
			totalNFTs += len(nfts)
		}

		output.WriteString(fmt.Sprintf("----- NFT BALANCES (%d total NFTs) -----\n", totalNFTs))

		for _, chain := range getKeys(nftBalances) {
			nfts := nftBalances[chain]
			if len(nfts) == 0 {
				continue
			}
			output.WriteString(fmt.Sprintf("  [%s] %d NFTs:\n", strings.ToUpper(chain), len(nfts)))

			for _, nft := range nfts {
				output.WriteString(fmt.Sprintf("    - %s | $%s | Amount: %s\n",
					nft.Name,
					nft.BalanceUSD.Text('f', 2),
					nft.Amount.Text('f', 0),
				))
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	// Pool Balances
	if global.ConfigFile.Debank.ParsePools && len(poolsData) > 0 {
		totalPools := 0
		for _, chainPools := range poolsData {
			for _, pools := range chainPools {
				totalPools += len(pools)
			}
		}

		output.WriteString(fmt.Sprintf("----- POOL BALANCES (%d total entries) -----\n", totalPools))

		for _, chain := range getKeys(poolsData) {
			projects := poolsData[chain]
			if len(projects) == 0 {
				continue
			}
			output.WriteString(fmt.Sprintf("  [%s] %d pools:\n", strings.ToUpper(chain), len(projects)))

			for _, poolName := range getKeys(projects) {
				pools := projects[poolName]
				if len(pools) == 0 {
					continue
				}
				output.WriteString(fmt.Sprintf("    [%s]\n", strings.ToUpper(poolName)))

				for _, pool := range pools {
					output.WriteString(fmt.Sprintf("      - %s | $%s | Amount: %s\n",
						pool.Name,
						pool.BalanceUSD.Text('f', 2),
						pool.Amount.Text('f', 4),
					))
				}
				output.WriteString("\n")
			}
			output.WriteString("\n")
		}
	}

	output.WriteString(strings.Repeat("-", 80))
	output.WriteString("\n\n")

	// Output file routing by balance range
	var fileName string
	switch {
	case totalUsdBalance < 1:
		fileName = "0_1_debank.txt"
	case totalUsdBalance < 10:
		fileName = "1_10_debank.txt"
	case totalUsdBalance < 100:
		fileName = "10_100_debank.txt"
	case totalUsdBalance < 500:
		fileName = "100_500_debank.txt"
	case totalUsdBalance < 1000:
		fileName = "500_1000_debank.txt"
	default:
		fileName = "1000_plus_debank.txt"
	}

	AppendToFile("./results/"+fileName, output.String())
}
