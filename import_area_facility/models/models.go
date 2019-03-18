package models

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/astaxie/beego"
)

var (
	ChannelId string
	UserName string
	OrgName string
	ChaincodeId string
	ConfigFile string
)

type ChaincodeSpec struct {
	client *channel.Client
	chaincodeId string
}

func Initialize(channelId string, userName string, orgName string, chaincodeId string, configFile string) (*ChaincodeSpec, error) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
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

func init()  {
	ChannelId = beego.AppConfig.String("channel_id")
	UserName = beego.AppConfig.String("user_name")
	OrgName = beego.AppConfig.String("org_name")
	ChaincodeId = beego.AppConfig.String("chaincode_id")
	ConfigFile = beego.AppConfig.String("config_file")
}
