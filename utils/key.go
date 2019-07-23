package utils

import (
	"crypto/ecdsa"
	ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
)

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	privateKeyData := ethereum_crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyData)
	return privateKeyHex
}

func DecodePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	if strings.HasPrefix(strings.ToLower(privateKeyHex), "0x") {
		if len(privateKeyHex) > 2 {
			privateKeyHex = string([]byte(privateKeyHex)[2:])
		}
	}

	return ethereum_crypto.HexToECDSA(privateKeyHex)
}

func AddressFromPrivateKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := ethereum_crypto.PubkeyToAddress(*publicKey).Hex()
	return address
}
