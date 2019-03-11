package handler

import (
	"context"

		example "190303/area/proto/example"
	"github.com/garyburd/redigo/redis"
	"190303/utils"
	"github.com/astaxie/beego"
	"190303/models"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetAreaList(ctx context.Context, req *example.GetAreaListRequest, rsp *example.GetAreaListResponse) error {
	aKey:="area_list"
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err==nil {
		data,_:=redis.Bytes(conn.Do("get",aKey))
		if data!=nil {
			beego.Info(data)
			rsp.Data=data
			return nil
		}
	}
	defer conn.Close()

	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getAreaList",[][]byte{})
	if err!=nil {
		return err
	}
	rsp.Data=data

	conn.Do("set",aKey,data,"EX",3600)

	return nil
}
