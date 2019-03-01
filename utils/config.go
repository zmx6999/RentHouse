package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

var (
	APPPort string

	FDFSHost string
	FDFSPort string
	FDFSConfig string

	FabricSDKConfig string
	ChannelId string
	User string
	ChaincodeId string

	RedisHost string
	RedisPort string

	MessageAppId string
	MessageAppKey string
	MessageProject string
)

func init()  {
	conf,err:=config.NewConfig("ini","/root/go/src/190222/conf/app.conf")
	if err!=nil {
		beego.Error(err)
		return
	}

	APPPort=conf.String("APP_PORT")

	FDFSHost=conf.String("FDFS_HOST")
	FDFSPort=conf.String("FDFS_PORT")
	FDFSConfig=conf.String("FDFS_CONFIG")

	FabricSDKConfig=conf.String("FABRIC_SDK_CONFIG")
	ChannelId=conf.String("CHANNEL_ID")
	User=conf.String("USER")
	ChaincodeId=conf.String("CHAINCODE_ID")

	RedisHost=conf.String("REDIS_HOST")
	RedisPort=conf.String("REDIS_PORT")

	MessageAppId=conf.String("MESSAGE_APP_ID")
	MessageAppKey=conf.String("MESSAGE_APP_KEY")
	MessageProject=conf.String("MESSAGE_PROJECT")
}
