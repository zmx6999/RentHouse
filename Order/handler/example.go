package handler

import (
	"context"

		example "190316/order/proto/example"
	"encoding/json"
	"190316/models"
	"190316/utils"
	"github.com/garyburd/redigo/redis"
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

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "addOrder", [][]byte{
		[]byte(utils.GetStringValue(data, "mobile", "")),
		[]byte(utils.GetStringValue(data, "house_id", "")),
		[]byte(utils.GetStringValue(data, "start_date", "")),
		[]byte(utils.GetStringValue(data, "end_date", "")),
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) GetList(ctx context.Context, req *example.GetListRequest, rsp *example.GetListResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getOrderList", [][]byte{[]byte(req.Mobile), []byte(req.Role)})
	if err != nil {
		return err
	}

	rsp.Data = data

	return nil
}

func (e *Example) Handle(ctx context.Context, req *example.HandleRequest, rsp *example.HandleResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "handleOrder", [][]byte{[]byte(req.Mobile), []byte(req.OrderId), []byte(req.Action)})
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) Comment(ctx context.Context, req *example.CommentRequest, rsp *example.CommentResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "comment", [][]byte{[]byte(req.Mobile), []byte(req.OrderId), []byte(req.Comment)})
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getOrderHouseId", [][]byte{[]byte(req.OrderId)})
	if err != nil {
		return err
	}

	cnn, err := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if err != nil {
		return err
	}
	defer cnn.Close()

	hKey := "house_" + string(data)
	_, err = cnn.Do("del", hKey)
	if err != nil {
		return err
	}

	return nil
}
