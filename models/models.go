package models

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/astaxie/beego"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
)

var (
	ChannelId string
	UserName string
	OrgName string
	ChaincodeId string
	ConfigFile string
)

func init()  {
	ChannelId = beego.AppConfig.String("channel_id")
	UserName = beego.AppConfig.String("username")
	OrgName = beego.AppConfig.String("org_name")
	ChaincodeId = beego.AppConfig.String("chaincode_id")
	ConfigFile = beego.AppConfig.String("config_file")
}

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

	return &ChaincodeSpec{client,chaincodeId}, nil
}

func (this *ChaincodeSpec) ChaincodeQuery(fcn string, args [][]byte) ([]byte, error) {
	request := channel.Request{ChaincodeID:this.chaincodeId, Fcn:fcn, Args:args}
	r, err := this.client.Query(request)
	return r.Payload, err
}

func (this *ChaincodeSpec) ChaincodeUpdate(fcn string, args [][]byte) ([]byte, error) {
	request := channel.Request{ChaincodeID:this.chaincodeId, Fcn:fcn, Args:args}
	r, err := this.client.Execute(request)
	return []byte(r.TransactionID), err
}
