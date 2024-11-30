package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"regexp"
	"strings"
)

func isMnemonic(input string) (bool, string) {
	if !bip39.IsMnemonicValid(input) {
		return false, ""
	}

	seed, err := bip39.NewSeedWithErrorChecking(input, "")
	if err != nil {
		return false, ""
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return false, ""
	}

	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return false, ""
	}
	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return false, ""
	}
	account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return false, ""
	}
	change, err := account.NewChildKey(0)
	if err != nil {
		return false, ""
	}
	addressKey, err := change.NewChildKey(0)
	if err != nil {
		return false, ""
	}

	privateKey, err := crypto.ToECDSA(addressKey.Key)
	if err != nil {
		return false, ""
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return true, address
}

func isPrivateKey(input string) (bool, string) {
	if strings.HasPrefix(input, "0x") {
		input = input[2:]
	}

	if len(input) != 64 {
		return false, ""
	}

	privateKeyBytes, err := hex.DecodeString(input)
	if err != nil {
		return false, ""
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return false, ""
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return true, address
}

func isEthAddress(input string) (bool, string) {
	loweredInput := strings.ToLower(input)

	if len(input) != 42 || !strings.HasPrefix(input, "0x") {
		return false, ""
	}

	match, _ := regexp.MatchString("^[0-9a-fA-F]{40}$", input[2:])
	if !match {
		return false, ""
	}

	address := common.HexToAddress(loweredInput)
	return loweredInput == strings.ToLower(address.Hex()), address.Hex()
}

func GetAccountAddress(target string) (string, error) {
	if valid, address := isMnemonic(target); valid {
		return address, nil
	}

	if valid, address := isPrivateKey(target); valid {
		return address, nil
	}

	if valid, address := isEthAddress(target); valid {
		return address, nil
	}

	return "", fmt.Errorf("wrong account credentials")
}
