package customTypes

import (
	"fmt"
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

type Proxy struct {
	User     string
	Password string
	IP       string
	Port     string
	Scheme   string
}

func (p *Proxy) GetAsString() string {
	var proxyString string

	if p.User != "" && p.Password != "" {
		proxyString = fmt.Sprintf("%s:%s@", p.User, p.Password)
	}

	proxyString += fmt.Sprintf("%s:%s", p.IP, p.Port)

	if p.Scheme != "" {
		if p.Scheme == "https" {
			p.Scheme = "http"
		}
		proxyString = fmt.Sprintf("%s://%s", p.Scheme, proxyString)
	} else {
		proxyString = fmt.Sprintf("http://%s", proxyString)
	}

	return proxyString
}

type RabbyReturnData struct {
	ChainName    string  `json:"chain_name"`
	ChainBalance float64 `json:"chain_balance"`
}
