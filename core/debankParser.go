package core

import (
	"debank_checker_v3/customTypes"
	"debank_checker_v3/global"
	"debank_checker_v3/utils"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/big"
	"net/url"
	"sort"
	"strings"
	"time"
)

type CustomBigFloat struct {
	*big.Float
}

func (c *CustomBigFloat) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		c.Float, _, _ = big.ParseFloat(str, 10, 0, big.ToNearestEven)
		return nil
	}
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	c.Float = new(big.Float).SetFloat64(f)
	return nil
}

func doRequest(accountAddress string,
	baseURL string,
	method string,
	path string,
	params url.Values,
	payload map[string]interface{}) ([]byte, error) {
	client := GetClient(utils.GetProxy())
	var err error
	var requestParams customTypes.RequestParamsStruct
	req := fasthttp.AcquireRequest()

	if strings.ToUpper(method) == fasthttp.MethodPost {
		if payload != nil {
			body, err := json.Marshal(payload)

			if err != nil {
				return nil, fmt.Errorf("%s | Failed to marshal payload: %v", accountAddress, err)
			}

			req.SetBody(body)
			req.Header.SetContentType("application/json")
		}

		err, requestParams = utils.GenerateSignature(nil, strings.ToUpper(method), path)

	} else {
		err, requestParams = utils.GenerateSignature(payload, strings.ToUpper(method), path)
	}

	if err != nil {
		return nil, fmt.Errorf("%s | Failed to generate request params: %v", accountAddress, err)
	}

	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
	req.Header.Set("origin", "https://debank.com")
	req.Header.Set("referer", "https://debank.com/")
	req.Header.Set("source", "web")
	req.Header.Set("x-api-ver", "v2")
	req.Header.Set("account", requestParams.AccountHeader)
	req.Header.Set("x-api-nonce", requestParams.Nonce)
	req.Header.Set("x-api-sign", requestParams.Signature)
	req.Header.Set("x-api-ts", requestParams.Timestamp)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("%s | Request error: %v", accountAddress, err)
	}

	if resp.StatusCode() == 429 {
		return nil, fmt.Errorf("%s | Rate limit", accountAddress)
	}

	respBody := make([]byte, len(resp.Body()))
	copy(respBody, resp.Body())
	return respBody, nil
}

func getTotalUsdBalance(accountAddress string) float64 {
	baseURL := "https://api.debank.com/asset/net_curve_24h"
	path := "/asset/net_curve_24h"
	params := url.Values{}
	params.Set("user_addr", strings.ToLower(accountAddress))

	payload := map[string]interface{}{
		"user_addr": strings.ToLower(accountAddress),
	}

	for {
		respBody, err := doRequest(accountAddress, baseURL, "GET", path, params, payload)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		var responseData struct {
			Data struct {
				UsdValueList [][]float64 `json:"usd_value_list"`
			} `json:"data"`
		}

		if err = json.Unmarshal(respBody, &responseData); err != nil {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		usdValueList := responseData.Data.UsdValueList

		if len(usdValueList) < 1 {
			log.Printf("%s | UsdValueList is Empty", accountAddress)
			continue
		}

		lastEntry := usdValueList[len(usdValueList)-1]

		if len(lastEntry) < 2 {
			log.Printf("%s | Last Entry Does Not Contain Enough Elements", accountAddress)
			continue
		}

		return lastEntry[1]
	}
}

func getUsedChains(accountAddress string, path string) []string {
	baseURL := "https://api.debank.com" + path
	var payload map[string]interface{}
	var responseData interface{}
	params := url.Values{}

	if path == "/nft/used_chains" {
		params.Set("user_addr", strings.ToLower(accountAddress))

		payload = map[string]interface{}{
			"user_addr": strings.ToLower(accountAddress),
		}

		responseData = &struct {
			Data []string `json:"data"`
		}{}
	} else if path == "/user/used_chains" {
		params.Set("id", strings.ToLower(accountAddress))

		payload = map[string]interface{}{
			"id": strings.ToLower(accountAddress),
		}

		responseData = &struct {
			Data struct {
				Chains []string `json:"chains"`
			} `json:"data"`
		}{}
	} else {
		log.Printf("%s | Wrong Path: %s", accountAddress, path)
		return nil
	}

	for {
		respBody, err := doRequest(accountAddress, baseURL, "GET", path, params, payload)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		if err = json.Unmarshal(respBody, &responseData); err != nil {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		switch v := responseData.(type) {
		case *struct {
			Data []string `json:"data"`
		}:
			return v.Data
		case *struct {
			Data struct {
				Chains []string `json:"chains"`
			} `json:"data"`
		}:
			return v.Data.Chains
		default:
			log.Printf("%s | Unexpected Response Format", accountAddress)
			continue
		}
	}
}

func getTokenBalances(accountAddress string, chains []string) map[string][]customTypes.TokenBalancesResultData {
	type tokenData struct {
		Amount          CustomBigFloat  `json:"amount"`
		Balance         big.Int         `json:"balance"`
		Name            string          `json:"name"`
		Price           *CustomBigFloat `json:"price"`
		ContractAddress string          `json:"id"`
	}

	type responseStruct struct {
		Data      []tokenData `json:"data"`
		ErrorCode *int        `json:"error_code"`
	}

	baseURL := "https://api.debank.com/token/balance_list"
	path := "/token/balance_list"

	params := url.Values{}
	params.Set("user_addr", strings.ToLower(accountAddress))

	result := make(map[string][]customTypes.TokenBalancesResultData)

	for _, currentChain := range chains {
		for {
			responseData := &responseStruct{}
			var tokensResultData []customTypes.TokenBalancesResultData
			params.Set("chain", currentChain)
			payload := map[string]interface{}{
				"user_addr": strings.ToLower(accountAddress),
				"chain":     currentChain,
			}

			respBody, err := doRequest(accountAddress, baseURL, "GET", path, params, payload)

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			if err = json.Unmarshal(respBody, responseData); err != nil {
				log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
				continue
			}

			for _, currentToken := range responseData.Data {
				var tokenInUsd *big.Float
				if currentToken.Price != nil {
					tokenInUsd = new(big.Float).Mul(currentToken.Price.Float, currentToken.Amount.Float)
				} else {
					tokenInUsd = new(big.Float)
				}
				tokensResultData = append(tokensResultData, customTypes.TokenBalancesResultData{
					Name:            currentToken.Name,
					BalanceUSD:      tokenInUsd,
					ContractAddress: currentToken.ContractAddress,
					Amount:          currentToken.Amount.Float,
				})
			}

			result[currentChain] = tokensResultData
			break
		}
	}

	return result
}

func getPoolBalances(accountAddress string) map[string]map[string][]customTypes.PoolBalancesResultData {
	type assetToken struct {
		Amount CustomBigFloat  `json:"amount"`
		Name   string          `json:"name"`
		Price  *CustomBigFloat `json:"price"`
	}

	type poolData struct {
		Chain             string `json:"chain"`
		Name              string `json:"name"`
		PortfolioItemList []struct {
			AssetTokenList []assetToken `json:"asset_token_list"`
		} `json:"portfolio_item_list"`
	}

	type responseStruct struct {
		Data      []poolData `json:"data"`
		ErrorCode *int       `json:"error_code"`
	}

	baseURL := "https://api.debank.com/portfolio/project_list"
	path := "/portfolio/project_list"

	params := url.Values{}
	params.Set("user_addr", strings.ToLower(accountAddress))

	payload := map[string]interface{}{
		"user_addr": strings.ToLower(accountAddress),
	}

	result := make(map[string]map[string][]customTypes.PoolBalancesResultData)

	for {
		respBody, err := doRequest(accountAddress, baseURL, "GET", path, params, payload)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		responseData := &responseStruct{}

		if err = json.Unmarshal(respBody, responseData); err != nil {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		for _, currentPool := range responseData.Data {
			for _, item := range currentPool.PortfolioItemList {
				for _, token := range item.AssetTokenList {
					if _, exists := result[currentPool.Chain]; !exists {
						result[currentPool.Chain] = make(map[string][]customTypes.PoolBalancesResultData)
					}
					if _, exists := result[currentPool.Chain][currentPool.Name]; !exists {
						result[currentPool.Chain][currentPool.Name] = []customTypes.PoolBalancesResultData{}
					}

					var tokenInUsd *big.Float
					if token.Price != nil {
						tokenInUsd = new(big.Float).Mul(token.Price.Float, token.Amount.Float)
					} else {
						tokenInUsd = new(big.Float)
					}
					result[currentPool.Chain][currentPool.Name] = append(result[currentPool.Chain][currentPool.Name],
						customTypes.PoolBalancesResultData{
							Name:       token.Name,
							Amount:     token.Amount.Float,
							BalanceUSD: tokenInUsd,
						})
				}
			}
		}

		break
	}

	return result
}

func getNftBalances(accountAddress string, chains []string) map[string][]customTypes.NftBalancesResultData {
	type collectionData struct {
		Amount          CustomBigFloat `json:"amount"`
		AvgPriceLast24h CustomBigFloat `json:"avg_price_last_24h"`
		Name            string         `json:"name"`
		NFTList         []struct {
			Amount CustomBigFloat `json:"amount"`
		} `json:"nft_list"`
		RankAt     *int `json:"rank_at"`
		SpentToken struct {
			Price *CustomBigFloat `json:"price"`
		} `json:"spent_token"`
	}
	type responseStruct struct {
		Data struct {
			Job *struct {
				Status string `json:"status"`
			} `json:"job"`
			Result struct {
				Data []collectionData `json:"data"`
			} `json:"result"`
		} `json:"data"`
		ErrorCode int `json:"error_code"`
	}

	baseURL := "https://api.debank.com/nft/collection_list"
	path := "/nft/collection_list"
	result := make(map[string][]customTypes.NftBalancesResultData)

	params := url.Values{}
	params.Set("user_addr", strings.ToLower(accountAddress))

	for _, currentChain := range chains {
		for {
			params.Set("chain", currentChain)
			payload := map[string]interface{}{
				"user_addr": strings.ToLower(accountAddress),
				"chain":     currentChain,
			}

			respBody, err := doRequest(accountAddress, baseURL, "GET", path, params, payload)

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			responseData := &responseStruct{}

			if err = json.Unmarshal(respBody, responseData); err != nil {
				log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
				continue
			}

			if responseData.Data.Job != nil && responseData.Data.Job.Status == "pending" {
				log.Printf("%s | NFT Balance Pending, sleeping 3 secs...", accountAddress)
				time.Sleep(3 * time.Second)
				continue
			}

			for _, currentNftData := range responseData.Data.Result.Data {
				var nftInUsd *big.Float

				if currentNftData.SpentToken.Price != nil {
					nftInUsd = new(big.Float).Mul(new(big.Float).Mul(currentNftData.AvgPriceLast24h.Float, currentNftData.SpentToken.Price.Float), currentNftData.Amount.Float)
				} else {
					nftInUsd = new(big.Float)
				}

				amountBigInt := new(big.Int)
				currentNftData.Amount.Float.Int(amountBigInt)
				result[currentChain] = append(result[currentChain],
					customTypes.NftBalancesResultData{
						Name:       currentNftData.Name,
						Amount:     currentNftData.Amount.Float,
						BalanceUSD: nftInUsd,
					})
			}
			break
		}
	}

	return result
}

func SortByBalanceUSD[T any](slice []T, balanceUSDGetter func(T) *big.Float) {
	sort.Slice(slice, func(i, j int) bool {
		return balanceUSDGetter(slice[i]).Cmp(balanceUSDGetter(slice[j])) > 0
	})
}

func SortMapByBalanceUSD[T any](dataMap map[string][]T, balanceUSDGetter func(T) *big.Float) {
	for key := range dataMap {
		SortByBalanceUSD(dataMap[key], balanceUSDGetter)
	}
}

func SortNestedMapByBalanceUSD[T any](data map[string]map[string][]T, getBalance func(T) *big.Float) {
	for _, innerMap := range data {
		for _, balances := range innerMap {
			SortByBalanceUSD(balances, getBalance)
		}
	}
}

func ParseDebankAccount(accountData string) {
	var (
		tokenBalances map[string][]customTypes.TokenBalancesResultData
		nftBalances   map[string][]customTypes.NftBalancesResultData
		poolsData     map[string]map[string][]customTypes.PoolBalancesResultData
	)

	accountAddress, _, _, err := utils.GetAccountData(accountData)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	totalUsdBalance := getTotalUsdBalance(accountAddress)
	log.Printf("%s | Total USD Balance: %f $", accountAddress, totalUsdBalance)

	if global.ConfigFile.DebankConfig.ParseTokens == true {
		tokenChainsUsed := getUsedChains(accountAddress, "/user/used_chains")
		log.Printf("%s | Token Chains Used: %d", accountAddress, len(tokenChainsUsed))

		if len(tokenChainsUsed) > 0 {
			tokenBalances = getTokenBalances(accountAddress, tokenChainsUsed)
			log.Printf("%s | Successfully Parsed Tokens", accountAddress)

			SortMapByBalanceUSD(tokenBalances, func(t customTypes.TokenBalancesResultData) *big.Float {
				return t.BalanceUSD
			})
		}
	}

	if global.ConfigFile.DebankConfig.ParseNfts == true {
		nftChainsUsed := getUsedChains(accountAddress, "/nft/used_chains")
		log.Printf("%s | NFT Chains Used: %d", accountAddress, len(nftChainsUsed))

		if len(nftChainsUsed) > 0 {
			nftBalances = getNftBalances(accountAddress, nftChainsUsed)
			log.Printf("%s | Successfully Parsed NFTs", accountAddress)

			SortMapByBalanceUSD(nftBalances, func(n customTypes.NftBalancesResultData) *big.Float {
				return n.BalanceUSD
			})
		}
	}

	if global.ConfigFile.DebankConfig.ParsePools == true {
		poolsData = getPoolBalances(accountAddress)
		log.Printf("%s | Successfully Parsed Pools", accountAddress)

		SortNestedMapByBalanceUSD(poolsData, func(p customTypes.PoolBalancesResultData) *big.Float {
			return p.BalanceUSD
		})

	}

	utils.FormatResult(accountData, accountAddress, totalUsdBalance, tokenBalances, nftBalances, poolsData)
}
