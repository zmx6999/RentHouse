package handler

import (
	"context"

		example "190222/house/proto/example"
	"encoding/json"
	"190222/utils"
	"190222/models"
	"github.com/garyburd/redigo/redis"
	"path"
	"errors"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Add(ctx context.Context, req *example.AddRequest, rsp *example.AddResponse) error {
	var data map[string]interface{}
	err:=json.Unmarshal(req.Data,&data)
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	userId,err:=ccs.GetUserId(data["mobile"].(string))
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate("addHouse",[][]byte{
		[]byte(data["title"].(string)),
		[]byte(data["price"].(string)),
		[]byte(data["area_id"].(string)),
		[]byte(data["address"].(string)),
		[]byte(data["room_count"].(string)),
		[]byte(data["acreage"].(string)),
		[]byte(data["unit"].(string)),
		[]byte(data["capacity"].(string)),
		[]byte(data["beds"].(string)),
		[]byte(data["deposit"].(string)),
		[]byte(data["min_days"].(string)),
		[]byte(data["max_days"].(string)),
		[]byte(data["facilities"].(string)),
		[]byte(userId),
	})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) GetHouseList(ctx context.Context, req *example.GetHouseListRequest, rsp *example.GetHouseListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	data,err:=ccs.ChaincodeQuery("getHouseList",[][]byte{[]byte(req.Mobile)})
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) GetHouseDesc(ctx context.Context, req *example.GetHouseDescRequest, rsp *example.GetHouseDescResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	data,err:=redis.Bytes(conn.Do("get","house_"+req.HouseId))
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

	data,err=ccs.ChaincodeQuery("getHouseInfo",[][]byte{[]byte(req.HouseId)})
	if err!=nil {
		return err
	}

	_,err=conn.Do("set","house_"+req.HouseId,data,"EX",3600)
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) UploadImage(ctx context.Context, req *example.UploadImageRequest, rsp *example.UploadImageResponse) error {
	if int64(len(req.Data))!=req.FileSize {
		return errors.New("file transfer error")
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	userId,err:=ccs.GetUserId(req.Mobile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate("addHouseImage",[][]byte{[]byte(req.HouseId),[]byte(fileId),[]byte(userId)})
	if err!=nil {
		return err
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	_,err=conn.Do("del","house_"+req.HouseId)
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) GetIndexList(ctx context.Context, req *example.GetIndexListRequest, rsp *example.GetIndexListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	data,err:=ccs.ChaincodeQuery("getIndexHouseList",[][]byte{})
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}
