package core

import (
	"debank_checker_v3/customTypes"
	"debank_checker_v3/utils"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strings"
)

func getSignL2(accountAddress string) string {
	baseURL := "https://api.debank.com/user/sign_v2"
	path := "/user/sign_v2"
	params := url.Values{}

	payload := map[string]interface{}{
		"id":       strings.ToLower(accountAddress),
		"chain_id": "1",
	}

	for {
		respBody, err := doRequest(accountAddress, baseURL, "POST", path, params, payload)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		var responseData struct {
			Data struct {
				Id   string `json:"id"`
				Text string `json:"text"`
			} `json:"data"`
			ErrorCode int `json:"error_code"`
		}

		if err = json.Unmarshal(respBody, &responseData); err != nil || responseData.ErrorCode != 0 {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		return responseData.Data.Text
	}
}

func signMessage(
	messageText string,
	privateKeyHex string,
) string {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)

	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	message := accounts.TextHash([]byte(messageText))

	signature, err := crypto.Sign(message, privateKey)

	if err != nil {
		log.Fatalf("Failed to sign message: %v", err)
	}

	signature[64] += 27

	return hexutil.Encode(signature)
}

func doAuth(
	accountAddress string,
	captchaResponse string,
	signedMessage string,
) string {
	baseURL := "https://api.debank.com/user/login_v2"
	path := "/user/login_v2"
	params := url.Values{}

	payload := map[string]interface{}{
		"token":     captchaResponse,
		"id":        strings.ToLower(accountAddress),
		"chain_id":  "1",
		"signature": signedMessage,
	}

	for {
		respBody, err := doRequest(accountAddress, baseURL, "POST", path, params, payload)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		var responseData struct {
			Data struct {
				SessionId string `json:"session_id"`
			} `json:"data"`
			ErrorCode int `json:"error_code"`
		}

		if err = json.Unmarshal(respBody, &responseData); err != nil || responseData.ErrorCode != 0 {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		return responseData.Data.SessionId
	}
}

func getL2Balance(
	accountAddress string,
	sessionId string,
) float64 {
	var err error
	var requestParams customTypes.RequestParamsStruct

	baseURL := "https://api.debank.com/user/l2_account"
	path := "/user/l2_account"
	params := url.Values{}
	params.Set("id", strings.ToLower(accountAddress))

	payload := map[string]interface{}{
		"id": strings.ToLower(accountAddress),
	}

	for {
		randomProxy := utils.GetProxy()
		client := GetClient(randomProxy)

		err, requestParams = utils.GenerateSignature(payload, "GET", path)

		if err != nil {
			log.Printf("%s | Failed to generate request params: %v", accountAddress, err)
			continue
		}

		var accountHeaderUnmarhsalled map[string]interface{}
		if err = json.Unmarshal([]byte(requestParams.AccountHeader), &accountHeaderUnmarhsalled); err != nil {
			log.Printf("%s | Error Decoding Account Headers: %s:", accountAddress, err)
			continue
		}

		accountHeaderUnmarhsalled["user_addr"] = strings.ToLower(accountAddress)
		accountHeaderUnmarhsalled["session_id"] = sessionId
		accountHeaderUnmarhsalled["wallet_type"] = "wallet"
		accountHeaderUnmarhsalled["is_verified"] = true

		accountHeaderFormatted, err := json.Marshal(accountHeaderUnmarhsalled)

		if err != nil {
			log.Printf("%s | Error Encoding Acccount Headers: %s", accountAddress, err)
			continue
		}

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
		req.Header.SetMethod(strings.ToUpper("GET"))
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
		req.Header.Set("origin", "https://debank.com")
		req.Header.Set("referer", "https://debank.com/")
		req.Header.Set("source", "web")
		req.Header.Set("x-api-ver", "v2")
		req.Header.Set("account", fmt.Sprintf("%s", accountHeaderFormatted))
		req.Header.Set("x-api-nonce", requestParams.Nonce)
		req.Header.Set("x-api-sign", requestParams.Signature)
		req.Header.Set("x-api-ts", requestParams.Timestamp)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		if err = client.Do(req, resp); err != nil {
			log.Printf("%s | Request error: %v", accountAddress, err)
			continue
		}

		if resp.StatusCode() == 429 {
			log.Printf("%s | Rate limit", accountAddress)
			continue
		}

		var responseData struct {
			Data struct {
				Balance float64 `json:"balance"`
			} `json:"data"`
			ErrorCode int `json:"error_code"`
		}

		if err = json.Unmarshal(resp.Body(), &responseData); err != nil || responseData.ErrorCode != 0 {
			log.Printf("%s | Failed To Parse JSON Response: %s", accountAddress, err)
			continue
		}

		return responseData.Data.Balance
	}
}

func DebankL2BalanceParser(accountData string) {
	accountData = utils.RemoveHexPrefix(accountData)
	privateKey := accountData
	accountAddress, accountType, privateKey, err := utils.GetAccountData(accountData)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	if accountType != 2 && accountType != 1 {
		log.Printf("%s | Account Data is Not Private Key/Mnemonic", accountAddress)
		return
	}

	messageText := getSignL2(accountAddress)
	messageSigned := signMessage(messageText, privateKey)
	captchaResponse := SolveCaptcha(accountAddress)

	sessionId := doAuth(accountAddress, captchaResponse, messageSigned)

	accountBalance := getL2Balance(accountAddress, sessionId)

	log.Printf("%s | Balance: %f $", accountAddress, accountBalance)

	if accountBalance > 0 {
		utils.AppendFile("./results/debank_l2_balances.txt",
			fmt.Sprintf("%s | %s | %f $\n", accountData, accountAddress, accountBalance))
	}
}
