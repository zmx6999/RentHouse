package utils

import (
	"net/http"
	"mime/multipart"
	"path"
	"strings"
	"errors"
	"github.com/weilaihui/fdfs_client"
	"crypto/ecdsa"
	"encoding/gob"
	"bytes"
	"crypto/elliptic"
	"encoding/hex"
	"math/rand"
	"crypto/sha512"
	crypto_rand "crypto/rand"
	"math/big"
)

func Find(x []string, needle string) int {
	for i := 0; i < len(x); i++ {
		if x[i] == needle {
			return i
		}
	}
	return -1
}

func PrepageUploadFile(r *http.Request, key string, allowedTypes []string, allowedMaxSize int64) ([]byte, *multipart.FileHeader, error) {
	file, head, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	ext := path.Ext(head.Filename)
	if Find(allowedTypes, strings.ToLower(ext[1:])) < 0 {
		return nil, nil, errors.New("FILE TYPE SHOULD BE " + strings.Join(allowedTypes, ","))
	}
	if head.Size > allowedMaxSize {
		return nil, nil, errors.New("FILE SIZE EXCEED")
	}

	data := make([]byte, head.Size)
	_, err = file.Read(data)
	if err != nil {
		return nil, nil, err
	}

	return data, head, nil
}

func UploadFile(data []byte, ext string) (string, error) {
	client, err := fdfs_client.NewFdfsClient(FastDFSConfigFile)
	if err != nil {
		return "", err
	}

	r, err := client.UploadByBuffer(data, ext)
	if err != nil {
		return "", err
	}

	return r.RemoteFileId, nil
}

func GetParam(key string, r *http.Request) string {
	query := r.URL.Query()[key]
	if query != nil && len(query) > 0 {
		return query[0]
	} else {
		return ""
	}
}

func GetStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if r, ok := data[key]; ok {
		return r.(string)
	}

	return defaultValue
}

func GetIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if r, ok := data[key]; ok {
		return r.(int)
	}

	return defaultValue
}

func GetFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if r, ok := data[key]; ok {
		return r.(float64)
	}

	return defaultValue
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) (string, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	gob.Register(elliptic.P256())
	err := encoder.Encode(privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer.Bytes()), nil
}

func DecodePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	data, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	gob.Register(elliptic.P256())
	var privateKey ecdsa.PrivateKey
	err = decoder.Decode(&privateKey)
	if err != nil {
		return nil, err
	}

	return &privateKey, nil
}

func Sign(data string, privateKeyHex string) ([]byte, []byte, error) {
	hash := GetSha512(data)

	privateKey, err := DecodePrivateKey(privateKeyHex)
	if err != nil {
		return nil, nil, err
	}

	r, s, err := ecdsa.Sign(crypto_rand.Reader, privateKey, hash)
	if err != nil {
		return nil, nil, err
	}

	return r.Bytes(), s.Bytes(), nil
}

func Verify(data string, publicKeyHex string, rByte []byte, sByte []byte) (bool, error) {
	hash := GetSha512(data)

	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, err
	}

	var r big.Int
	var s big.Int
	r.SetBytes(rByte)
	s.SetBytes(sByte)

	var x big.Int
	var y big.Int
	x.SetBytes(publicKey[:len(publicKey) / 2])
	y.SetBytes(publicKey[len(publicKey) / 2:])
	rawPublicKey := ecdsa.PublicKey{elliptic.P256(), &x, &y}

	return ecdsa.Verify(&rawPublicKey, hash, &r, &s), nil
}

func GetNonce(length int) string {
	x := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	str := ""
	for i := 0; i < length; i++ {
		str += x[rand.Intn(len(x))]
	}
	return str
}

func GetSha512(x string) []byte {
	h := sha512.New()
	h.Write([]byte(x))
	return h.Sum(nil)
}
