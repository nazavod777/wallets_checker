package customTypes

import (
	"math/big"
)

type RequestParamsStruct struct {
	AccountHeader string
	Nonce         string
	Signature     string
	Timestamp     string
}

type TokenBalancesResultData struct {
	Amount          *big.Float `json:"amount"`
	Name            string     `json:"name"`
	ContractAddress string     `json:"contract_address"`
	BalanceUSD      *big.Float `json:"balance_usd"`
}

type PoolBalancesResultData struct {
	Amount     *big.Float `json:"amount"`
	Name       string     `json:"name"`
	BalanceUSD *big.Float `json:"balance_usd"`
}

type NftBalancesResultData struct {
	Amount     *big.Float `json:"amount"`
	Name       string     `json:"name"`
	BalanceUSD *big.Float `json:"price_usd"`
}

type RabbyReturnData struct {
	ChainName    string  `json:"chain_name"`
	ChainBalance float64 `json:"chain_balance"`
}

type ConfigStruct struct {
	DebankConfig struct {
		ParseTokens bool `json:"parse_tokens"`
		ParseNfts   bool `json:"parse_nfts"`
		ParsePools  bool `json:"parse_pools"`
	} `json:"debank_config"`
	TwoCaptchaApiKey string `json:"2captcha_apikey"`
}
