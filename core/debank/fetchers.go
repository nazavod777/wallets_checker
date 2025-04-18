package debank

import (
	"debank_checker_v3/core/debankRequest"
	"debank_checker_v3/customTypes"
	"debank_checker_v3/global"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math/big"
	"net/url"
	"strings"
	"time"
)

// getTotalUsdBalance retrieves the total USD balance for a given account.
func getTotalUsdBalance(addr string) float64 {
	base := "https://api.debank.com/asset/net_curve_24h"
	path := "/asset/net_curve_24h"
	params := url.Values{"user_addr": {strings.ToLower(addr)}}
	payload := map[string]interface{}{"user_addr": strings.ToLower(addr)}

	for {
		body, statusCode, err := debankRequest.MakeRequest(addr, "GET", base, path, params, payload, "")
		if err != nil {
			log.Printf("[%d/%d] | %s | error getting total balance [%d]: %v",
				global.CurrentProgress, global.TargetProgress, addr, statusCode, err)
			continue
		}

		data := gjson.GetBytes(body, "data.usd_value_list")
		if !data.Exists() || data.Type == gjson.Null || len(data.Array()) < 1 {
			return 0
		}

		last := data.Array()[len(data.Array())-1].Array()
		if len(last) < 2 {
			log.Printf("[%d/%d] | %s | malformed last entry in balance history",
				global.CurrentProgress, global.TargetProgress, addr)
			return 0
		}

		return last[1].Float()
	}
}

// getUsedChains retrieves a list of chains used by an account (for tokens or NFTs).
func getUsedChains(addr string, path string) []string {
	base := "https://api.debank.com" + path
	params := url.Values{}
	payload := make(map[string]interface{})
	var key string

	switch path {
	case "/user/used_chains":
		params.Set("id", strings.ToLower(addr))
		payload["id"] = strings.ToLower(addr)
		key = "data.chains"
	case "/nft/used_chains":
		params.Set("user_addr", strings.ToLower(addr))
		payload["user_addr"] = strings.ToLower(addr)
		key = "data"
	default:
		log.Errorf("[%d/%d] | %s | invalid chain path: %s",
			global.CurrentProgress, global.TargetProgress, addr, path)
		return nil
	}

	for {
		body, statusCode, err := debankRequest.MakeRequest(addr, "GET", base, path, params, payload, "")
		if err != nil {
			log.Printf("[%d/%d] | %s | error fetching used chains [%d]: %v",
				global.CurrentProgress, global.TargetProgress, addr, statusCode, err)
			continue
		}

		result := gjson.GetBytes(body, key)
		if result.Exists() && result.IsArray() {
			var chains []string
			result.ForEach(func(_, val gjson.Result) bool {
				chains = append(chains, val.String())
				return true
			})
			return chains
		}

		log.Printf("[%d/%d] | %s | unexpected chain response [%d]",
			global.CurrentProgress, global.TargetProgress, addr, statusCode)
	}
}

// getTokenBalances retrieves ERC-20 token balances by chain.
func getTokenBalances(addr string, chains []string) map[string][]customTypes.TokenBalance {
	base := "https://api.debank.com/token/balance_list"
	path := "/token/balance_list"
	result := make(map[string][]customTypes.TokenBalance)

	for _, chain := range chains {
		params := url.Values{
			"user_addr": {strings.ToLower(addr)},
			"chain":     {chain},
		}
		payload := map[string]interface{}{
			"user_addr": strings.ToLower(addr),
			"chain":     chain,
		}

		for {
			body, statusCode, err := debankRequest.MakeRequest(addr, "GET", base, path, params, payload, "")
			if err != nil {
				log.Printf("[%d/%d] | %s | error fetching token balances [%d]: %v",
					global.CurrentProgress, global.TargetProgress, addr, statusCode, err)
				continue
			}

			tokens := gjson.GetBytes(body, "data")

			if !tokens.Exists() || !tokens.IsArray() {
				log.Printf("[%d/%d] | %s | wrong response when getting token balances [%d]: %s",
					global.CurrentProgress, global.TargetProgress, addr, statusCode, string(body))
				continue
			}

			var tokenList []customTypes.TokenBalance

			tokens.ForEach(func(_, t gjson.Result) bool {
				name := t.Get("name").String()
				contract := t.Get("id").String()
				amount := parseBigFloat(t.Get("amount"))
				price := parseBigFloat(t.Get("price"))

				usd := new(big.Float).Mul(price, amount)

				tokenList = append(tokenList, customTypes.TokenBalance{
					Name:            name,
					ContractAddress: contract,
					Amount:          amount,
					BalanceUSD:      usd,
				})

				return true
			})

			result[chain] = tokenList
			break
		}
	}

	return result
}

// getPoolBalances retrieves liquidity pool balances grouped by project and chain.
func getPoolBalances(addr string) map[string]map[string][]customTypes.PoolBalance {
	base := "https://api.debank.com/portfolio/project_list"
	path := "/portfolio/project_list"
	params := url.Values{"user_addr": {strings.ToLower(addr)}}
	payload := map[string]interface{}{"user_addr": strings.ToLower(addr)}

	result := make(map[string]map[string][]customTypes.PoolBalance)

	for {
		body, statusCode, err := debankRequest.MakeRequest(addr, "GET", base, path, params, payload, "")
		if err != nil {
			log.Printf("[%d/%d] | %s | error fetching pool balances [%d]: %v",
				global.CurrentProgress, global.TargetProgress, addr, statusCode, err)
			continue
		}

		data := gjson.GetBytes(body, "data")

		if !data.Exists() || !data.IsArray() {
			log.Printf("[%d/%d] | %s | wrong response when getting pool balances [%d]: %s",
				global.CurrentProgress, global.TargetProgress, addr, statusCode, string(body))
			continue
		}

		data.ForEach(func(_, pool gjson.Result) bool {
			chain := pool.Get("chain").String()
			project := pool.Get("name").String()

			if result[chain] == nil {
				result[chain] = make(map[string][]customTypes.PoolBalance)
			}

			pool.Get("portfolio_item_list").ForEach(func(_, item gjson.Result) bool {
				item.Get("asset_token_list").ForEach(func(_, token gjson.Result) bool {
					amount := parseBigFloat(token.Get("amount"))
					price := parseBigFloat(token.Get("price"))
					usd := new(big.Float).Mul(price, amount)

					result[chain][project] = append(result[chain][project], customTypes.PoolBalance{
						Name:       token.Get("name").String(),
						Amount:     amount,
						BalanceUSD: usd,
					})
					return true
				})
				return true
			})

			return true
		})

		break
	}

	return result
}

// getNftBalances retrieves NFT balances and their estimated USD value by chain.
func getNftBalances(addr string, chains []string) map[string][]customTypes.NFTBalance {
	base := "https://api.debank.com/nft/collection_list"
	path := "/nft/collection_list"
	result := make(map[string][]customTypes.NFTBalance)

	for _, chain := range chains {
		params := url.Values{
			"user_addr": {strings.ToLower(addr)},
			"chain":     {chain},
		}
		payload := map[string]interface{}{
			"user_addr": strings.ToLower(addr),
			"chain":     chain,
		}

		for {
			body, statusCode, err := debankRequest.MakeRequest(addr, "GET", base, path, params, payload, "")
			if err != nil {
				log.Printf("[%d/%d] | %s | error fetching NFT balances [%d]: %v",
					global.CurrentProgress, global.TargetProgress, addr, statusCode, err)
				continue
			}

			json := gjson.ParseBytes(body)
			if json.Get("data.job.status").String() == "pending" {
				log.Infof("[%d/%d] | %s | NFT job pending, waiting 3s...", global.CurrentProgress, global.TargetProgress, addr)
				time.Sleep(3 * time.Second)
				continue
			}

			data := gjson.GetBytes(body, "data.result.data")

			if !data.Exists() || !data.IsArray() {
				log.Printf("[%d/%d] | %s | wrong response when getting nft balances [%d]: %s",
					global.CurrentProgress, global.TargetProgress, addr, statusCode, string(body))
				continue
			}

			data.ForEach(func(_, v gjson.Result) bool {
				amount := parseBigFloat(v.Get("amount"))
				price := parseBigFloat(v.Get("avg_price_last_24h"))
				tokenPrice := parseBigFloat(v.Get("spent_token.price"))

				usd := new(big.Float)
				if tokenPrice.Cmp(big.NewFloat(0)) > 0 {
					usd.Mul(new(big.Float).Mul(price, tokenPrice), amount)
				}

				result[chain] = append(result[chain], customTypes.NFTBalance{
					Name:       v.Get("name").String(),
					Amount:     amount,
					BalanceUSD: usd,
				})
				return true
			})

			break
		}
	}

	return result
}
