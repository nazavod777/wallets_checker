package debank

import (
	"github.com/tidwall/gjson"
	"math/big"
	"sort"
)

func parseBigFloat(val gjson.Result) *big.Float {
	if !val.Exists() {
		return new(big.Float)
	}
	f, _, err := big.ParseFloat(val.String(), 10, 256, big.ToNearestEven)
	if err != nil {
		return new(big.Float)
	}
	return f
}

func sortByBalance[T any](data []T, balanceFunc func(T) *big.Float) {
	sort.Slice(data, func(i, j int) bool {
		return balanceFunc(data[i]).Cmp(balanceFunc(data[j])) > 0
	})
}

func sortMapByBalance[T any](data map[string][]T, balanceFunc func(T) *big.Float) {
	for _, v := range data {
		sortByBalance(v, balanceFunc)
	}
}

func sortNestedMapByBalance[T any](nested map[string]map[string][]T, balanceFunc func(T) *big.Float) {
	for _, m := range nested {
		for _, v := range m {
			sortByBalance(v, balanceFunc)
		}
	}
}
