package handler

import (
	"context"

		example "190303/house/proto/example"
	"encoding/json"
	"190303/utils"
	"190303/models"
	"path"
	"errors"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Add(ctx context.Context, req *example.AddRequest, rsp *example.AddResponse) error {
	var data map[string]interface{}
	err:=json.Unmarshal(req.Data,&data)
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"addHouse",[][]byte{
		[]byte(data["mobile"].(string)),
		[]byte(data["title"].(string)),
		[]byte(utils.GetStringValue(data,"price","0")),
		[]byte(data["area_id"].(string)),
		[]byte(utils.GetStringValue(data,"address","")),
		[]byte(utils.GetStringValue(data,"room_count","1")),
		[]byte(utils.GetStringValue(data,"acreage","0")),
		[]byte(utils.GetStringValue(data,"unit","")),
		[]byte(utils.GetStringValue(data,"capacity","1")),
		[]byte(utils.GetStringValue(data,"beds","")),
		[]byte(utils.GetStringValue(data,"deposit","0")),
		[]byte(utils.GetStringValue(data,"min_days","1")),
		[]byte(utils.GetStringValue(data,"max_days","0")),
		[]byte(utils.GetStringValue(data,"facility","")),
	})
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) UploadImage(ctx context.Context, req *example.UploadImageRequest, rsp *example.UploadImageResponse) error {
	if len(req.Data)!=int(req.FileSize) {
		return errors.New("file transfer error")
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"uploadHouseImage",[][]byte{[]byte(req.Mobile),[]byte(req.HouseId),[]byte(fileId)})
	if err!=nil {
		return err
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}

	hKey:="house_"+req.HouseId
	_,err=conn.Do("del",hKey)
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) GetLandlordList(ctx context.Context, req *example.GetLandlordListRequest, rsp *example.GetLandlordListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getLandlordHouseList",[][]byte{[]byte(req.Mobile)})
	if err!=nil {
		return err
	}
	rsp.Data=data

	return nil
}

func (e *Example) GetDesc(ctx context.Context, req *example.GetDescRequest, rsp *example.GetDescResponse) error {
	hKey:="house_"+req.HouseId
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err==nil {
		data,_:=redis.Bytes(conn.Do("get",hKey))
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

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getHouseDesc",[][]byte{[]byte(req.HouseId)})
	if err!=nil {
		return err
	}
	rsp.Data=data

	conn.Do("set",hKey,data,"EX",3600)

	return nil
}

func (e *Example) GetIndexList(ctx context.Context, req *example.GetIndexListRequest, rsp *example.GetIndexListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getIndexHouseList",[][]byte{})
	if err!=nil {
		return err
	}
	rsp.Data=data

	return nil
}

func (e *Example) Search(ctx context.Context, req *example.SearchRequest, rsp *example.SearchResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"searchHouse",[][]byte{[]byte(req.AreaId),[]byte(req.Start),[]byte(req.End),[]byte(req.Page)})
	if err!=nil {
		return err
	}
	rsp.Data=data

	return nil
}
