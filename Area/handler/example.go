package handler

import (
	"context"
	example "190222/area/proto/example"
	"190222/models"
	"190222/utils"
		"github.com/garyburd/redigo/redis"
)

type Example struct{}

func (e *Example) GetArea(ctx context.Context, req *example.GetAreaRequest, rsp *example.GetAreaResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	data,err:=redis.Bytes(conn.Do("get","area_list"))
	if data!=nil {
		rsp.Data=data
		rsp.Code=utils.RECODE_OK
		rsp.Msg=utils.RecodeText(rsp.Code)

		return nil
	}

	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	data,err=ccs.ChaincodeQuery("getAreaList",[][]byte{})
	if err!=nil {
		return err
	}

	_,err=conn.Do("set","area_list",data,"EX",3600)
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}
