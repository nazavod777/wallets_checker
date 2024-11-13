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
	Amount     *big.Float `json:"amount"`
	Name       string     `json:"name"`
	BalanceUSD *big.Float `json:"balance_usd"`
}

type PoolBalancesResultData struct {
	Amount     *big.Float `json:"amount"`
	Name       string     `json:"name"`
	BalanceUSD *big.Float `json:"balance_usd"`
}

type NftBalancesResultData struct {
	Amount   *big.Float `json:"amount"`
	Name     string     `json:"name"`
	PriceUSD *big.Float `json:"price_usd"`
}

type RabbyReturnData struct {
	ChainName    string  `json:"chain_name"`
	ChainBalance float64 `json:"chain_balance"`
}
