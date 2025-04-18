package debankl2

import (
	"debank_checker_v3/core/debankRequest"
	"debank_checker_v3/global"
	"github.com/ethereum/go-ethereum/crypto"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func getSignL2(address string) string {
	urlPath := "/user/sign_v2"
	base := "https://api.debank.com" + urlPath
	params := url.Values{}
	payload := map[string]interface{}{
		"id":       strings.ToLower(address),
		"chain_id": "1",
	}

	for {
		body, statusCode, err := debankRequest.MakeRequest(address, base, "POST", urlPath, params, payload, "")
		if err != nil {
			log.Printf("[%d/%d] | %s | error requesting sign L2 [%d]: %v",
				global.CurrentProgress, global.TargetProgress, address, statusCode, err)
			continue
		}
		if msg := gjson.GetBytes(body, "data.text").String(); msg != "" {
			return msg
		}

		log.Printf("[%d/%d] | %s | unexpected sign L2 response [%d]: %s",
			global.CurrentProgress, global.TargetProgress, address, statusCode, body)
	}
}

func signMessage(messageText, privKeyHex string) string {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	msgHash := accounts.TextHash([]byte(messageText))
	sig, err := crypto.Sign(msgHash, privKey)
	if err != nil {
		log.Fatalf("Failed to sign message: %v", err)
	}

	sig[64] += 27
	return hexutil.Encode(sig)
}

func doAuth(address, captcha, signature string) string {
	urlPath := "/user/login_v2"
	base := "https://api.debank.com" + urlPath
	params := url.Values{}
	payload := map[string]interface{}{
		"token":     captcha,
		"id":        strings.ToLower(address),
		"chain_id":  "1",
		"signature": signature,
	}

	for {
		body, statusCode, err := debankRequest.MakeRequest(address, base, "POST", urlPath, params, payload, "")
		if err != nil {
			log.Printf("[%d/%d] | %s | error during login [%d]: %v",
				global.CurrentProgress, global.TargetProgress, address, statusCode, err)
			continue
		}
		if session := gjson.GetBytes(body, "data.session_id").String(); session != "" {
			return session
		}

		log.Printf("[%d/%d] | %s | invalid login response [%d]: %s",
			global.CurrentProgress, global.TargetProgress, address, statusCode, body)
	}
}
