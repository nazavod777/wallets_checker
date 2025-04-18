package debankl2

import (
	"debank_checker_v3/core/debankRequest"
	"debank_checker_v3/global"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
)

func getL2Balance(
	address string,
	sessionID string,
) float64 {
	path := "/user/l2_account"
	base := "https://api.debank.com" + path

	params := url.Values{}
	params.Set("id", strings.ToLower(address))

	payload := map[string]interface{}{
		"id": strings.ToLower(address),
	}

	for {
		body, statusCode, err := debankRequest.MakeRequest(address, "GET", base, path, params, payload, sessionID)
		if err != nil {
			log.Printf("[%d/%d] | %s | error requesting L2 balance [%d]: %v",
				global.CurrentProgress, global.TargetProgress, address, statusCode, err)
			continue
		}

		balance := gjson.GetBytes(body, "balance")

		if gjson.GetBytes(body, "error_code").Int() != 0 || balance.Type == gjson.Null {
			log.Printf("[%d/%d] | %s | bad L2 balance response [%d]: %s",
				global.CurrentProgress, global.TargetProgress, address, statusCode, string(body))
			continue
		}

		return balance.Float()
	}
}
