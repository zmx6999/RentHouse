package utils

import (
	"net/http"
	"mime/multipart"
	"path"
	"errors"
	"strings"
	"github.com/weilaihui/fdfs_client"
	"crypto/ecdsa"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"encoding/hex"
	"math/rand"
				)

func PrepareUploadFile(r *http.Request,key string,allowTypes []string,allowMaxSize int) ([]byte,*multipart.FileHeader,error) {
	file,head,err:=r.FormFile(key)
	if err!=nil {
		return nil,nil,err
	}
	defer file.Close()

	ext:=strings.ToLower(path.Ext(head.Filename))
	if Find(allowTypes,ext[1:])<0 {
		return nil,nil,errors.New("file type should be "+strings.Join(allowTypes,","))
	}
	if int(head.Size)>allowMaxSize {
		return nil,nil,errors.New("file size exceeds")
	}

	data:=make([]byte,head.Size)
	_,err=file.Read(data)
	if err!=nil {
		return nil,nil,err
	}

	return data,head,nil
}

func UploadFile(data []byte,ext string) (string,error) {
	client,err:=fdfs_client.NewFdfsClient(FastDFSConfigFile)
	if err!=nil {
		return "", err
	}

	r,err:=client.UploadByBuffer(data,ext)
	if err!=nil {
		return "", err
	}

	return r.RemoteFileId, nil
}

func GetParam(key string,r *http.Request) string {
	query:=r.URL.Query()[key]
	if query!=nil && len(query)>0 {
		return query[0]
	} else {
		return ""
	}
}

func Find(list []string,needle string) int {
	for i:=0; i<len(list); i++ {
		if list[i]==needle {
			return i
		}
	}
	return -1
}

func GenerateNonce(length int) string {
	x:=[]string{"0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	n:=len(x)
	str:=""
	for i:=0; i<n; i++ {
		str+=x[rand.Intn(n)]
	}
	return str
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) (string,error) {
	var buffer bytes.Buffer

	encoder:=gob.NewEncoder(&buffer)
	gob.Register(elliptic.P256())
	err:=encoder.Encode(privateKey)
	if err!=nil {
		return "",err
	}

	privateKeyHex:=hex.EncodeToString(buffer.Bytes())
	return privateKeyHex,nil
}

func DecodePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey,error) {
	data,err:=hex.DecodeString(privateKeyHex)
	if err!=nil {
		return nil,err
	}

	var buffer bytes.Buffer
	_,err=buffer.Write(data)
	if err!=nil {
		return nil,err
	}

	decoder:=gob.NewDecoder(&buffer)
	gob.Register(elliptic.P256())

	var privateKey ecdsa.PrivateKey
	err=decoder.Decode(&privateKey)
	if err!=nil {
		return nil,err
	}

	return &privateKey,nil
}

func GetStringValue(data map[string]interface{},key string,defaultValue string) string {
	value,ok:=data[key]
	if ok {
		return value.(string)
	} else {
		return defaultValue
	}
}

func GetFloat64Value(data map[string]interface{},key string,defaultValue float64) float64 {
	value,ok:=data[key]
	if ok {
		return value.(float64)
	} else {
		return defaultValue
	}
}
