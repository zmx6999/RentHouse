package handler

import (
	"context"

		example "190303/order/proto/example"
	"encoding/json"
	"190303/utils"
	"190303/models"
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

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"addOrder",[][]byte{
		[]byte(data["mobile"].(string)),
		[]byte(data["house_id"].(string)),
		[]byte(data["start_date"].(string)),
		[]byte(data["end_date"].(string)),
	})
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) GetList(ctx context.Context, req *example.GetListRequest, rsp *example.GetListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getOrderList",[][]byte{[]byte(req.Mobile),[]byte(req.Role)})
	if err!=nil {
		return err
	}
	rsp.Data=data

	return nil
}

func (e *Example) Handle(ctx context.Context, req *example.HandleRequest, rsp *example.HandleResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"handleOrder",[][]byte{[]byte(req.Mobile),[]byte(req.OrderId),[]byte(req.Action)})
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) Comment(ctx context.Context, req *example.CommentRequest, rsp *example.CommentResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"comment",[][]byte{[]byte(req.Mobile),[]byte(req.OrderId),[]byte(req.Comment)})
	if err!=nil {
		return err
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort,redis.DialPassword(utils.RedisPassword))
	if err!=nil {
		return err
	}

	houseId,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getOrderHouseId",[][]byte{[]byte(req.OrderId)})
	if err!=nil {
		return err
	}

	hKey:="house_"+string(houseId)
	_,err=conn.Do("del",hKey)
	if err!=nil {
		return err
	}

	return nil
}
