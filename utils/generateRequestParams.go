package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"debank_checker_v3/customTypes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

func generateRandomID() string {
	const chars = "abcdef0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[mathRand.Intn(len(chars))]
	}
	return string(result)
}

func sortQueryString(queryString string) string {

	params, _ := url.ParseQuery(queryString)

	var paramKeys []string
	for paramKey := range params {
		paramKeys = append(paramKeys, paramKey)
	}
	sort.Strings(paramKeys)

	var sortedQueryParams []string
	for _, paramKey := range paramKeys {
		sortedQueryParams = append(sortedQueryParams, fmt.Sprintf("%s=%s", paramKey, params[paramKey][0]))
	}
	return strings.Join(sortedQueryParams, "&")
}

func customSha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	hashBytes := h.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func generateNonce(length int) (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	nonce := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		nonce[i] = letters[num.Int64()]
	}

	return "n_" + string(nonce), nil
}

func hmacSha256(key []byte, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	hmacBytes := h.Sum(nil)

	return hex.EncodeToString(hmacBytes)
}

func mapToQueryString(payload map[string]interface{}) string {
	values := url.Values{}
	for key, value := range payload {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}

func GenerateSignature(payload map[string]interface{}, method string, path string) (error, customTypes.RequestParamsStruct) {
	nonce, err := generateNonce(40)

	if err != nil {
		return err, customTypes.RequestParamsStruct{}
	}

	queryString := mapToQueryString(payload)

	timestamp := time.Now().Unix()
	randStr := fmt.Sprintf(
		"debank-api\n%s\n%d",
		nonce,
		timestamp,
	)
	randStrHash := customSha256(randStr)

	requestParams := fmt.Sprintf(
		"%s\n%s\n%s",
		strings.ToUpper(method),
		strings.ToLower(path),
		sortQueryString(strings.ToLower(queryString)),
	)
	requestParamsHash := customSha256(requestParams)

	info := map[string]interface{}{
		"random_at": timestamp,
		"random_id": generateRandomID(),
		"user_addr": "",
	}
	accountHeader, err := json.Marshal(info)

	if err != nil {
		return err, customTypes.RequestParamsStruct{}
	}

	signature := hmacSha256(
		[]byte(randStrHash),
		[]byte(requestParamsHash),
	)

	result := customTypes.RequestParamsStruct{
		AccountHeader: string(accountHeader),
		Nonce:         fmt.Sprintf("%v", nonce),
		Signature:     fmt.Sprintf("%v", signature),
		Timestamp:     fmt.Sprintf("%d", timestamp),
	}

	return nil, result
}
