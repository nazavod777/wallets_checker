package customTypes

import "math/big"

type RequestParams struct {
	AccountHeader string
	Nonce         string
	Signature     string
	Timestamp     string
}

type TokenBalance struct {
	Name            string     `json:"name"`
	Amount          *big.Float `json:"amount"`
	ContractAddress string     `json:"contract_address"`
	BalanceUSD      *big.Float `json:"balance_usd"`
}

type PoolBalance struct {
	Name       string     `json:"name"`
	Amount     *big.Float `json:"amount"`
	BalanceUSD *big.Float `json:"balance_usd"`
}

type NFTBalance struct {
	Name       string     `json:"name"`
	Amount     *big.Float `json:"amount"`
	BalanceUSD *big.Float `json:"price_usd"`
}

type RabbyChainBalance struct {
	ChainName    string  `json:"chain_name"`
	ChainBalance float64 `json:"chain_balance"`
}

type Config struct {
	Debank struct {
		ParseTokens bool `json:"parse_tokens"`
		ParseNfts   bool `json:"parse_nfts"`
		ParsePools  bool `json:"parse_pools"`
	} `json:"debank_config"`
	TwoCaptchaAPIKey string `json:"2captcha_apikey"`
}
