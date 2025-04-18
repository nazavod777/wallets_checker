package util

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// isMnemonic checks if the given input is a valid BIP39 mnemonic phrase,
// and returns (true, privateKeyHex, address) if so.
func isMnemonic(input string) (bool, string, string) {
	if !bip39.IsMnemonicValid(input) {
		return false, "", ""
	}

	seed, err := bip39.NewSeedWithErrorChecking(input, "")
	if err != nil {
		return false, "", ""
	}

	// Standard BIP44 derivation path: m/44'/60'/0'/0/0
	path := []uint32{
		bip32.FirstHardenedChild + 44,
		bip32.FirstHardenedChild + 60,
		bip32.FirstHardenedChild + 0,
		0,
		0,
	}

	key := deriveBIP32(seed, path)
	if key == nil {
		return false, "", ""
	}

	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		return false, "", ""
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return true, hex.EncodeToString(crypto.FromECDSA(privateKey)), address
}

// isPrivateKey checks if the input is a valid hex-encoded private key.
func isPrivateKey(input string) (bool, string) {
	input = RemoveHexPrefix(input)

	if len(input) != 64 {
		return false, ""
	}

	bytes, err := hex.DecodeString(input)
	if err != nil {
		return false, ""
	}

	privateKey, err := crypto.ToECDSA(bytes)
	if err != nil {
		return false, ""
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return true, address
}

// isEthAddress validates an Ethereum address format (0x-prefixed, 40 hex chars).
func isEthAddress(input string) (bool, string) {
	if len(input) != 42 || !strings.HasPrefix(input, "0x") {
		return false, ""
	}

	match, _ := regexp.MatchString("^[0-9a-fA-F]{40}$", input[2:])
	if !match {
		return false, ""
	}

	lowered := strings.ToLower(input)
	return lowered == strings.ToLower(common.HexToAddress(lowered).Hex()), common.HexToAddress(lowered).Hex()
}

// GetAccountData determines the type of input (mnemonic, private key, or address)
// and returns (address, type, privateKey, error).
// type: 1 = mnemonic, 2 = privateKey, 3 = address
func GetAccountData(target string) (string, int, string, error) {
	if ok, privKey, addr := isMnemonic(target); ok {
		return addr, 1, privKey, nil
	}
	if ok, addr := isPrivateKey(target); ok {
		return addr, 2, target, nil
	}
	if ok, addr := isEthAddress(target); ok {
		return addr, 3, "", nil
	}
	return "", 0, "", fmt.Errorf("invalid account input: %s", target)
}

// deriveBIP32 applies a derivation path to the seed and returns the final key.
func deriveBIP32(seed []byte, path []uint32) *bip32.Key {
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil
	}

	for _, step := range path {
		key, err = key.NewChildKey(step)
		if err != nil {
			return nil
		}
	}
	return key
}
