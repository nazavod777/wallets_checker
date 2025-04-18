package debankl2

import (
	"bytes"
	"debank_checker_v3/global"
	"debank_checker_v3/util"
	"encoding/json"
	fhttp "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
	"time"
)

const (
	createTaskURL    = "https://api.2captcha.com/createTask"
	getTaskResultURL = "https://api.2captcha.com/getTaskResult"
	siteKey          = "6Lcw7ewpAAAAAPtZi4LTNCAWmWj-1h5ACTD_CQHE"
	siteURL          = "https://debank.com/"
	softID           = 4759
)

func createCaptchaTask(client tlsclient.HttpClient, key string) string {
	payload := map[string]interface{}{
		"clientKey": global.ConfigFile.TwoCaptchaAPIKey,
		"softId":    softID,
		"task": map[string]interface{}{
			"type":       "RecaptchaV2TaskProxyless",
			"websiteURL": siteURL,
			"websiteKey": siteKey,
			"minScore":   0.9,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Panicf("%s | Failed to encode captcha task payload: %v", key, err)
	}

	for {
		req, err := fhttp.NewRequest("POST", createTaskURL, bytes.NewReader(body))
		if err != nil {
			log.Errorf("%s | Failed to create request (createTask): %v", key, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("%s | Request error (createTask): %v", key, err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("%s | Read body error (createTask): %v", key, err)

			if err = resp.Body.Close(); err != nil {
				log.Warnf("error closing body: %v", err)
			}

			continue
		}

		jsonResp := gjson.ParseBytes(respBody)
		if jsonResp.Get("errorId").Int() != 0 {
			log.Errorf("%s | 2Captcha createTask error: %s", key, string(respBody))

			if err = resp.Body.Close(); err != nil {
				log.Warnf("error closing body: %v", err)
			}

			return ""
		}

		taskID := jsonResp.Get("taskId").Int()

		if err = resp.Body.Close(); err != nil {
			log.Warnf("error closing body: %v", err)
		}

		return strconv.Itoa(int(taskID))
	}
}

func pollCaptchaResult(client tlsclient.HttpClient, key, taskID string) *string {
	payload := map[string]interface{}{
		"clientKey": global.ConfigFile.TwoCaptchaAPIKey,
		"softId":    softID,
		"taskId":    taskID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Panicf("%s | Failed to encode result payload: %v", key, err)
	}

	for {
		req, err := fhttp.NewRequest("POST", getTaskResultURL, bytes.NewReader(body))
		if err != nil {
			log.Errorf("%s | Failed to create request (getTaskResult): %v", key, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("%s | Request error (getTaskResult): %v", key, err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("%s | Read body error (getTaskResult): %v", key, err)

			if err = resp.Body.Close(); err != nil {
				log.Warnf("error closing body: %v", err)
			}

			continue
		}

		jsonResp := gjson.ParseBytes(respBody)

		if jsonResp.Get("errorId").Int() != 0 {
			log.Errorf("%s | 2Captcha getTaskResult error: %s", key, string(respBody))

			if err = resp.Body.Close(); err != nil {
				log.Warnf("error closing body: %v", err)
			}

			return nil
		}

		if jsonResp.Get("status").String() == "ready" {
			token := jsonResp.Get("solution.token").String()

			if err = resp.Body.Close(); err != nil {
				log.Warnf("error closing body: %v", err)
			}

			return &token
		}

		log.Infof("%s | Captcha still processing... retrying in 5s", key)
		time.Sleep(5 * time.Second)
	}
}

func solveCaptcha(privateKey string) string {
	client := util.CreateClient("")

	for {
		taskID := createCaptchaTask(client, privateKey)
		if taskID == "" {
			continue
		}

		result := pollCaptchaResult(client, privateKey, taskID)
		if result != nil {
			return *result
		}
	}
}
