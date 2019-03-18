package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego"
)

var (
	HttpPort string

	ChannelId string
	UserName string
	OrgName string
	ChaincodeId string
	ConfigFile string

	FastDFSConfigFile string

	RedisHost string
	RedisPort string
	RedisPassword string

	MessageAppId string
	MessageAppKey string
	MessageProject string
)

func init()  {
	conf, err := config.NewConfig("ini", "/root/go/src/190316/conf/app.conf")
	if err != nil {
		beego.Error(err)
		return
	}

	HttpPort = conf.String("httpport")

	ChannelId = conf.String("channel_id")
	UserName = conf.String("user_name")
	OrgName = conf.String("org_name")
	ChaincodeId = conf.String("chaincode_id")
	ConfigFile = conf.String("config_file")

	FastDFSConfigFile = conf.String("fast_dfs_config_file")

	RedisHost = conf.String("redis_host")
	RedisPort = conf.String("redis_port")
	RedisPassword = conf.String("redis_password")

	MessageAppId = conf.String("message_app_id")
	MessageAppKey = conf.String("message_app_key")
	MessageProject = conf.String("message_project")
}
