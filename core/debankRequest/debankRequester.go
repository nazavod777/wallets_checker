package debankRequest

import (
	"bytes"
	"debank_checker_v3/customTypes"
	"debank_checker_v3/util"
	"encoding/json"
	"fmt"
	fhttp "github.com/bogdanfinn/fhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"net/url"
	"strings"
)

func MakeRequest(
	address string,
	method string,
	baseURL string,
	path string,
	params url.Values,
	payload map[string]interface{},
	sessionID string) ([]byte, int, error) {
	client := util.GetClient()

	var (
		reqBody      io.Reader
		paramsStruct customTypes.RequestParams
		err          error
	)

	upperMethod := strings.ToUpper(method)
	if upperMethod == fhttp.MethodPost {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to encode payload: %v", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
		err, paramsStruct = util.GenerateSignatureDebank(nil, upperMethod, path)
	} else {
		reqBody = bytes.NewReader([]byte{})
		err, paramsStruct = util.GenerateSignatureDebank(payload, upperMethod, path)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate signature: %v", err)
	}

	// Подставляем session_id
	if gjson.Valid(paramsStruct.AccountHeader) {
		accountHeader := paramsStruct.AccountHeader
		accountHeader, _ = sjson.Set(accountHeader, "user_addr", strings.ToLower(address))
		accountHeader, _ = sjson.Set(accountHeader, "session_id", sessionID)
		accountHeader, _ = sjson.Set(accountHeader, "wallet_type", "wallet")
		accountHeader, _ = sjson.Set(accountHeader, "is_verified", true)
		paramsStruct.AccountHeader = accountHeader
	}

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := fhttp.NewRequest(upperMethod, reqURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build request: %v", err)
	}

	req.Header.Set("Account", paramsStruct.AccountHeader)
	req.Header.Set("X-Api-Nonce", paramsStruct.Nonce)
	req.Header.Set("X-Api-Sign", paramsStruct.Signature)
	req.Header.Set("X-Api-Ts", paramsStruct.Timestamp)
	req.Header.Set("X-Api-Ver", "v2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (DebankBot)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://debank.com")
	req.Header.Set("Referer", "https://debank.com/")
	req.Header.Set("Source", "web")
	req.Header.Set("Accept-Language", "ru,en;q=0.9")

	if upperMethod == fhttp.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("error closing body: %v", err)
		}
	}()

	status := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("failed to read response: %v", err)
	}

	return body, status, nil
}
