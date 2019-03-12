package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

var (
	AppPort string

	ChannelId string
	UserName string
	OrgName string
	ChaincodeId string
	ConfigFile string

	FastDFSHost string
	FastDFSPort string
	FastDFSConfigFile string

	RedisHost string
	RedisPort string
	RedisPassword string

	MessageAppId string
	MessageAppKey string
	MessageProject string
)

func init()  {
	conf,err:=config.NewConfig("ini","conf/app.conf")
	if err!=nil {
		beego.Error(err)
		return
	}

	AppPort=conf.String("app_port")

	ChannelId=conf.String("channel_id")
	UserName=conf.String("username")
	OrgName=conf.String("org_name")
	ChaincodeId=conf.String("chaincode_id")
	ConfigFile=conf.String("config_file")

	FastDFSHost=conf.String("fast_dfs_host")
	FastDFSPort=conf.String("fast_dfs_port")
	FastDFSConfigFile=conf.String("fast_dfs_config_file")

	RedisHost=conf.String("redis_host")
	RedisPort=conf.String("redis_port")
	RedisPassword=conf.String("redis_password")

	MessageAppId=conf.String("message_app_id")
	MessageAppKey=conf.String("message_app_key")
	MessageProject=conf.String("message_project")
}
