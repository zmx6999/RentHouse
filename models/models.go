package models

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"190303/utils"
	"encoding/hex"
	"math/big"
	"crypto/ecdsa"
	"crypto/elliptic"
	crypto_rand "crypto/rand"
)

type ChaincodeSpec struct {
	client *channel.Client
	chaincodeId string
}

func Initialize(channelId string, userName string, orgName string, chaincodeId string, configFile string) (*ChaincodeSpec, error) {
	sdk,err := fabsdk.New(config.FromFile(utils.ConfigFile))
	if err != nil {
		return nil, err
	}

	clientContext := sdk.ChannelContext(channelId, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientContext)
	if err != nil {
		return nil, err
	}

	return &ChaincodeSpec{client, chaincodeId}, nil
}

func (this *ChaincodeSpec) ChaincodeQuery(chaincodeId string, function string, args [][]byte) ([]byte, error) {
	request := channel.Request{ChaincodeID: chaincodeId, Fcn: function, Args: args}
	r, err := this.client.Query(request)
	return r.Payload, err
}

func (this *ChaincodeSpec) ChaincodeUpdate(chaincodeId string, function string, args [][]byte) ([]byte, error) {
	request := channel.Request{ChaincodeID: chaincodeId, Fcn: function, Args: args}
	r, err := this.client.Execute(request)
	return []byte(r.TransactionID), err
}



func GetPublicKey(mobile string) ([]byte,error) {
	ccs,err:=Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return nil,err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getUserPublicKey",[][]byte{[]byte(mobile)})
	if err!=nil {
		return nil,err
	}

	return data,nil
}

func Sign(privateKeyHex string,hash []byte) ([]byte,[]byte,error) {
	privateKey,err:=utils.DecodePrivateKey(privateKeyHex)
	if err!=nil {
		return nil,nil,err
	}

	r,s,err:=ecdsa.Sign(crypto_rand.Reader,privateKey,hash)
	if err!=nil {
		return nil,nil,err
	}

	return r.Bytes(),s.Bytes(),nil
}

func Verify(publicKeyHex string,rData []byte,sData []byte,hash []byte) (bool,error) {
	publicKey,err:=hex.DecodeString(publicKeyHex)
	if err!=nil {
		return false,nil
	}

	var x,y big.Int
	x.SetBytes(publicKey[:len(publicKey)/2])
	y.SetBytes(publicKey[len(publicKey)/2:])
	rawPublicKey:=ecdsa.PublicKey{elliptic.P256(),&x,&y}

	var r,s big.Int
	r.SetBytes(rData)
	s.SetBytes(sData)

	return ecdsa.Verify(&rawPublicKey,hash,&r,&s),nil
}

func VerifyUser(privateKeyHex string,mobile string) (bool,error) {
	str:=utils.GenerateNonce(41)

	rData,sData,err:=Sign(privateKeyHex,[]byte(str))
	if err!=nil {
		return false,nil
	}

	publicKeyHex,err:=GetPublicKey(mobile)
	if err!=nil {
		return false,nil
	}

	return Verify(string(publicKeyHex),rData,sData,[]byte(str))
}
