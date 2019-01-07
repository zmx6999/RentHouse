package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego"
)

var (
	APPHost string
	APPPort string

	MySQLHost string
	MySQLPort string
	MySQLDB string
	MySQLUserName string
	MySQLPassword string

	RedisHost string
	RedisPort string

	FastDFSHost string
	FastDFSPort string

	TimeZone string
)

func initConfig()  {
	appConf,err:=config.NewConfig("ini","/root/go/src/190105/conf/app.conf")
	if err!=nil {
		beego.Debug(err)
		return
	}
	APPHost=appConf.String("APPHost")
	APPPort=appConf.String("APPPort")

	MySQLHost=appConf.String("MySQLHost")
	MySQLPort=appConf.String("MySQLPort")
	MySQLDB=appConf.String("MySQLDB")
	MySQLUserName=appConf.String("MySQLUserName")
	MySQLPassword=appConf.String("MySQLPassword")

	RedisHost=appConf.String("RedisHost")
	RedisPort=appConf.String("RedisPort")

	FastDFSHost=appConf.String("FastDFSHost")
	FastDFSPort=appConf.String("FastDFSPort")

	TimeZone=appConf.String("TimeZone")
}

func init()  {
	initConfig()
}
