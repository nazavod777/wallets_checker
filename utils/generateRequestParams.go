package utils

import (
	"bufio"
	"debank_checker_v3/customTypes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type RequestParams struct {
	RandomAt string `json:"random_at"`
	RandomID string `json:"random_id"`
	UserAddr string `json:"user_addr"`
}

func generateRandomID() string {
	const chars = "abcdef0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func GenerateRequestParams(payload map[string]interface{}, method, path string) (customTypes.RequestParamsStruct, error) {
	cmd := exec.Command("node", "js/main.js")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	if err := cmd.Start(); err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	method = strings.ToUpper(method)
	fmt.Fprintf(stdin, "%s|%s|%s\n", payloadData, method, path)
	stdin.Close()
	scanner := bufio.NewScanner(stdout)
	var outputData string
	if scanner.Scan() {
		outputData = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	var returnData map[string]interface{}
	if err := json.Unmarshal([]byte(outputData), &returnData); err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	rTime := fmt.Sprintf("%d", time.Now().Unix())
	info := map[string]interface{}{
		"random_at": rTime,
		"random_id": generateRandomID(),
		"user_addr": "", // Set to empty as per Python code's default value
	}
	accountHeader, err := json.Marshal(info)
	if err != nil {
		return customTypes.RequestParamsStruct{}, err
	}
	returnData["account_header"] = string(accountHeader)
	requestParams := customTypes.RequestParamsStruct{
		AccountHeader: string(accountHeader),
		Nonce:         fmt.Sprintf("%v", returnData["nonce"]),     // Assuming "nonce" exists in returnData
		Signature:     fmt.Sprintf("%v", returnData["signature"]), // Assuming "signature" exists in returnData
		Timestamp:     rTime,
	}
	return requestParams, nil
}
