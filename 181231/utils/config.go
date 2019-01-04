package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego"
)

var (
	ServerName string
	ServerHost string
	ServerPort string
	RedisHost string
	RedisPort string
	RedisDBNum string
	MySQLHost string
	MySQLPort string
	MySQLDBName string
	FastDFSHost string
	FastDFSPort string
)

func initConfig()  {
	appconf,err:=config.NewConfig("ini","/root/go/src/sss/181231/conf/app.conf")
	if err!=nil {
		beego.Debug(err)
		return
	}
	ServerName=appconf.String("ServerName")
	ServerHost=appconf.String("ServerHost")
	ServerPort=appconf.String("ServerPort")
	RedisHost=appconf.String("RedisHost")
	RedisPort=appconf.String("RedisPort")
	RedisDBNum=appconf.String("RedisDBNum")
	MySQLHost=appconf.String("MySQLHost")
	MySQLPort=appconf.String("MySQLPort")
	MySQLDBName=appconf.String("MySQLDBName")
	FastDFSHost=appconf.String("FastDFSHost")
	FastDFSPort=appconf.String("FastDFSPort")
}

func init()  {
	initConfig()
}
