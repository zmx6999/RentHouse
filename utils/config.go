package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego"
)

var (
	APPPort string

	MySQLHost string
	MySQLPort string
	MySQLDBName string
	MySQLUserName string
	MySQLPassword string

	RedisHost string
	RedisPort string

	FDFSHost string
	FDFSPort string

	TimeZone string
)

func init()  {
	cfg,err:=config.NewConfig("ini","/root/go/src/190120/conf/app.conf")
	if err!=nil {
		beego.Error(err)
		return
	}

	APPPort=cfg.String("APPPort")

	MySQLHost=cfg.String("MySQLHost")
	MySQLPort=cfg.String("MySQLPort")
	MySQLDBName=cfg.String("MySQLDBName")
	MySQLUserName=cfg.String("MySQLUserName")
	MySQLPassword=cfg.String("MySQLPassword")

	RedisHost=cfg.String("RedisHost")
	RedisPort=cfg.String("RedisPort")

	FDFSHost=cfg.String("FDFSHost")
	FDFSPort=cfg.String("FDFSPort")

	TimeZone=cfg.String("TimeZone")
}
