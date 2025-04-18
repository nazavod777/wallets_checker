package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"debank_checker_v3/customTypes"
)

// === Utility Functions ===

// sha256Hex returns the SHA256 hash of the given string as a hex-encoded string.
func sha256Hex(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

// sha256Bytes returns the SHA256 hash of the given byte slice.
func sha256Bytes(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// xorHash applies XOR with the given byte to every character of the input hash string.
func xorHash(hash string, xorByte byte) string {
	result := make([]byte, len(hash))
	for i := 0; i < len(hash); i++ {
		result[i] = hash[i] ^ xorByte
	}
	return string(result)
}

// reverseQueryString builds and reverses the query string from a map.
func reverseQueryString(payload map[string]interface{}) string {
	params := make([]string, 0, len(payload))
	for key, value := range payload {
		params = append(params, fmt.Sprintf("%s=%v", url.QueryEscape(key), url.QueryEscape(fmt.Sprintf("%v", value))))
	}

	// Reverse the order
	for i, j := 0, len(params)-1; i < j; i, j = i+1, j-1 {
		params[i], params[j] = params[j], params[i]
	}

	return strings.Join(params, "&")
}

// generateNonce creates a random alphanumeric string of the given length.
func generateNonce(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	nonce := make([]byte, length)
	for i := range nonce {
		nonce[i] = charset[rand.Intn(len(charset))]
	}
	return string(nonce), nil
}

// generateRandomID returns a random 32-character hex string.
func generateRandomID() string {
	const hexChars = "abcdef0123456789"
	id := make([]byte, 32)
	for i := range id {
		id[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(id)
}

// === Signature Generator ===

// GenerateSignatureDebank builds the required headers and signature for accessing Debank API endpoints.
func GenerateSignatureDebank(payload map[string]interface{}, method string, path string) (error, customTypes.RequestParams) {
	nonce, err := generateNonce(40)
	if err != nil {
		return err, customTypes.RequestParams{}
	}

	timestamp := time.Now().Unix()
	queryString := reverseQueryString(payload)

	// Build the string to be hashed
	method = strings.ToUpper(method)
	data1 := method + "\n" + path + "\n" + queryString
	data2 := fmt.Sprintf("debank-api\nn_%s\n%d", nonce, timestamp)

	hash1 := sha256Hex(data1)
	hash2 := sha256Hex(data2)

	xor1 := xorHash(hash2, 54)
	xor2 := xorHash(hash2, 92)

	h1 := sha256Bytes([]byte(xor1 + hash1))
	h2 := sha256Bytes(append([]byte(xor2), h1...))

	signature := hex.EncodeToString(h2)

	accountInfo := map[string]interface{}{
		"random_at": timestamp,
		"random_id": generateRandomID(),
		"user_addr": nil,
	}

	accountHeader, err := json.Marshal(accountInfo)
	if err != nil {
		return err, customTypes.RequestParams{}
	}

	result := customTypes.RequestParams{
		AccountHeader: string(accountHeader),
		Nonce:         "n_" + nonce,
		Signature:     signature,
		Timestamp:     fmt.Sprintf("%d", timestamp),
	}

	return nil, result
}
