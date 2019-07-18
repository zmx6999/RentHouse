package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
)

func AddressFromPrivateKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*publicKey).Hex()
	return addr
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	privateKeyData := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyData)
	return privateKeyHex
}

func DecodePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	if strings.HasPrefix(strings.ToLower(privateKeyHex), "0x") {
		privateKeyData := []byte(privateKeyHex)
		if len(privateKeyData) > 2 {
			privateKeyData = privateKeyData[2:]
			privateKeyHex = string(privateKeyData)
		}
	}

	return crypto.HexToECDSA(privateKeyHex)
}
