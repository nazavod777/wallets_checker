package debank

import (
	"debank_checker_v3/customTypes"
	"debank_checker_v3/global"
	"debank_checker_v3/util"
	"log"
	"math/big"
)

type ParseDebank struct{}

func (p ParseDebank) Parse(accountData string) {
	addr, _, _, err := util.GetAccountData(accountData)
	if err != nil {
		log.Printf("[%d/%d] | %s | failed to get account info: %v", global.CurrentProgress, global.TargetProgress, accountData, err)
		return
	}

	totalBalance := getTotalUsdBalance(addr)
	log.Printf("[%d/%d] | %s | Total Balance: $%.2f", global.CurrentProgress, global.TargetProgress, addr, totalBalance)

	var (
		tokens                 map[string][]customTypes.TokenBalance
		nfts                   map[string][]customTypes.NFTBalance
		pools                  map[string]map[string][]customTypes.PoolBalance
		tokenChains, nftChains []string
	)

	if global.ConfigFile.Debank.ParseTokens && totalBalance > 0 {
		tokenChains = getUsedChains(addr, "/user/used_chains")
		if len(tokenChains) > 0 {
			tokens = getTokenBalances(addr, tokenChains)
			sortMapByBalance(tokens, func(t customTypes.TokenBalance) *big.Float { return t.BalanceUSD })
		}
	}

	if global.ConfigFile.Debank.ParseNfts {
		nftChains = getUsedChains(addr, "/nft/used_chains")
		if len(nftChains) > 0 {
			nfts = getNftBalances(addr, nftChains)
			sortMapByBalance(nfts, func(n customTypes.NFTBalance) *big.Float { return n.BalanceUSD })
		}
	}

	if global.ConfigFile.Debank.ParsePools && totalBalance > 0 {
		pools = getPoolBalances(addr)
		sortNestedMapByBalance(pools, func(p customTypes.PoolBalance) *big.Float { return p.BalanceUSD })
	}

	util.FormatResult(accountData, addr, totalBalance, len(tokenChains), len(nftChains), tokens, nfts, pools)
}
