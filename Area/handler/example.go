package handler

import (
	"context"

		example "190316/area/proto/example"
	"github.com/gomodule/redigo/redis"
	"190316/utils"
	"190316/models"
	)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetList(ctx context.Context, req *example.GetListRequest, rsp *example.GetListResponse) error {
	aKey := "area_list"
	cnn, _ := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if cnn != nil {
		defer cnn.Close()
		data, _ := redis.Bytes(cnn.Do("get", aKey))
		if data != nil {
			rsp.Data = data
			return nil
		}
	}

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getAreaList", [][]byte{})
	if err != nil {
		return err
	}

	rsp.Data = data

	if cnn != nil {
		cnn.Do("set", aKey, data, "EX", 3600)
	}

	return nil
}
