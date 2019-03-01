package models

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"crypto/ecdsa"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	)

type ChaincodeSpec struct {
	client apitxn.ChannelClient
	chaincodeId string
}

func Initialize(channelId string,user string,chaincodeId string,configFile string) (*ChaincodeSpec,error) {
	sdk,err:=fabapi.NewSDK(fabapi.Options{ConfigFile:configFile})
	if err!=nil {
		return nil,err
	}

	client,err:=sdk.NewChannelClient(channelId,user)
	if err!=nil {
		return nil,err
	}

	return &ChaincodeSpec{client,chaincodeId}, nil
}

func (this *ChaincodeSpec) ChaincodeQuery(function string,args [][]byte) ([]byte,error) {
	request:=apitxn.QueryRequest{ChaincodeID:this.chaincodeId,Fcn:function,Args:args}
	return this.client.Query(request)
}

func (this *ChaincodeSpec) ChaincodeUpdate(function string,args [][]byte) ([]byte,error) {
	request:=apitxn.ExecuteTxRequest{ChaincodeID:this.chaincodeId,Fcn:function,Args:args}
	id,err:=this.client.ExecuteTx(request)
	return []byte(id.ID),err
}

func (this *ChaincodeSpec) Close()  {
	this.client.Close()
}

func Sign(privateKeyHex string,data []byte) ([]byte,[]byte,error) {
	privateKeyData,err:=hex.DecodeString(privateKeyHex)
	if err!=nil {
		return nil,nil,err
	}

	privateKey,err:=DecodePrivateKey(privateKeyData)
	if err!=nil {
		return nil,nil,err
	}

	r,s,err:=ecdsa.Sign(crypto_rand.Reader,privateKey,data)
	if err!=nil {
		return nil,nil,err
	}

	return r.Bytes(),s.Bytes(),nil
}

func Verify(publicKeyHex string,data []byte,rbyte []byte,sbyte []byte) (bool,error) {
	publicKeyData,err:=hex.DecodeString(publicKeyHex)
	if err!=nil {
		return false,err
	}

	x:=big.Int{}
	x.SetBytes(publicKeyData[:len(publicKeyData)/2])

	y:=big.Int{}
	y.SetBytes(publicKeyData[len(publicKeyData)/2:])

	publicKey:=ecdsa.PublicKey{elliptic.P256(),&x,&y}

	r:=big.Int{}
	r.SetBytes(rbyte)

	s:=big.Int{}
	s.SetBytes(sbyte)

	return ecdsa.Verify(&publicKey,data,&r,&s),nil
}

func (this *ChaincodeSpec) GetPublicKeyHex(userMobile string) (string,error) {
	data,err:=this.ChaincodeQuery("getUserPublicKey",[][]byte{[]byte(userMobile)})
	if err!=nil {
		return "",err
	}

	return string(data),nil
}

func GenerateNonce(length int) string {
	x:=[]string{"0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	n:=len(x)
	y:=""
	for i:=0; i<length; i++ {
		y+=x[rand.Intn(n)]
	}
	return y
}

func (this *ChaincodeSpec) VerifyUser(privateKeyHex string,userMobile string) (bool,error) {
	nonce:=GenerateNonce(41)
	nonceData,err:=json.Marshal(nonce)
	if err!=nil {
		return false,err
	}

	rbyte,sbyte,err:=Sign(privateKeyHex,nonceData)
	if err!=nil {
		return false,err
	}

	publicKeyHex,err:=this.GetPublicKeyHex(userMobile)
	if err!=nil {
		return false,err
	}

	return Verify(publicKeyHex,nonceData,rbyte,sbyte)
}

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) ([]byte,error) {
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	gob.Register(elliptic.P256())
	err:=encoder.Encode(privateKey)
	if err!=nil {
		return nil,err
	}

	return buffer.Bytes(),nil
}

func DecodePrivateKey(privateKeyData []byte) (*ecdsa.PrivateKey,error) {
	var buffer bytes.Buffer
	_,err:=buffer.Write(privateKeyData)
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

func (this *ChaincodeSpec) GetUserId(userMobile string) (string,error) {
	_data,err:=this.ChaincodeQuery("getUserInfo",[][]byte{[]byte(userMobile)})
	if err!=nil {
		return "",err
	}

	var data map[string]interface{}
	err=json.Unmarshal(_data,&data)
	if err!=nil {
		return "",err
	}

	return data["user_id"].(string),nil
}
