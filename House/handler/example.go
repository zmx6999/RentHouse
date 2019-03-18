package handler

import (
	"context"

		example "190316/house/proto/example"
	"encoding/json"
	"190316/models"
	"190316/utils"
		"path"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Add(ctx context.Context, req *example.AddRequest, rsp *example.AddResponse) error {
	var data map[string]interface{}
	json.Unmarshal(req.Data, &data)

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "addHouse", [][]byte{
		[]byte(utils.GetStringValue(data, "mobile", "")),
		[]byte(utils.GetStringValue(data, "title", "")),
		[]byte(utils.GetStringValue(data, "price", "0")),
		[]byte(utils.GetStringValue(data, "area_id", "")),
		[]byte(utils.GetStringValue(data, "address", "")),
		[]byte(utils.GetStringValue(data, "room_count", "1")),
		[]byte(utils.GetStringValue(data, "acreage", "0")),
		[]byte(utils.GetStringValue(data, "unit", "")),
		[]byte(utils.GetStringValue(data, "capacity", "1")),
		[]byte(utils.GetStringValue(data, "beds", "")),
		[]byte(utils.GetStringValue(data, "deposit", "0")),
		[]byte(utils.GetStringValue(data, "min_days", "1")),
		[]byte(utils.GetStringValue(data, "max_days", "0")),
		[]byte(utils.GetStringValue(data, "facility", "")),
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) GetLandlordList(ctx context.Context, req *example.GetLandlordListRequest, rsp *example.GetLandlordListResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getHouseList", [][]byte{[]byte(req.Mobile)})
	if err != nil {
		return err
	}

	rsp.Data = data

	return nil
}

func (e *Example) UploadImage(ctx context.Context, req *example.UploadImageRequest, rsp *example.UploadImageResponse) error {
	if len(req.Data) != int(req.FileSize) {
		return errors.New("file transfer error")
	}

	ext := path.Ext(req.FileName)
	fileId, err := utils.UploadFile(req.Data, ext[1:])
	if err != nil {
		return err
	}

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "uploadImage", [][]byte{[]byte(req.Mobile), []byte(req.HouseId), []byte(fileId)})
	if err != nil {
		return err
	}

	cnn, err := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if err != nil {
		return err
	}
	defer cnn.Close()

	hKey := "house_" + req.HouseId
	_, err = cnn.Do("del", hKey)
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) GetDesc(ctx context.Context, req *example.GetDescRequest, rsp *example.GetDescResponse) error {
	hKey := "house_" + req.HouseId
	cnn, _ := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if cnn != nil {
		defer cnn.Close()
		data, _ := redis.Bytes(cnn.Do("get", hKey))
		if data != nil {
			rsp.Data = data
			beego.Info(data)
			return nil
		}
	}

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getHouseDesc", [][]byte{[]byte(req.HouseId)})
	if err != nil {
		return err
	}

	rsp.Data = data

	if cnn != nil {
		cnn.Do("set", hKey, data, "EX", 3600)
	}

	return nil
}

func (e *Example) Search(ctx context.Context, req *example.SearchRequest, rsp *example.SearchResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "searchHouse", [][]byte{[]byte(req.AreaId), []byte(req.Start), []byte(req.End), []byte(req.Page)})
	if err != nil {
		return err
	}

	rsp.Data = data

	return nil
}
