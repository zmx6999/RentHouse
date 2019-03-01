package handler

import (
	"context"

		example "190222/order/proto/example"
	"encoding/json"
	"190222/utils"
	"190222/models"
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

	_,err=ccs.ChaincodeUpdate("addOrder",[][]byte{
		[]byte(userId),
		[]byte(data["house_id"].(string)),
		[]byte(data["start_date"].(string)),
		[]byte(data["end_date"].(string)),
	})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) GetList(ctx context.Context, req *example.GetListRequest, rsp *example.GetListResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	userId,err:=ccs.GetUserId(req.Mobile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery("getOrderList",[][]byte{[]byte(userId),[]byte(req.Role)})
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Handle(ctx context.Context, req *example.HandleRequest, rsp *example.HandleResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	userId,err:=ccs.GetUserId(req.Mobile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate("handleOrder",[][]byte{
		[]byte(req.OrderId),
		[]byte(userId),
		[]byte(req.Action),
	})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Comment(ctx context.Context, req *example.CommentRequest, rsp *example.CommentResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	userId,err:=ccs.GetUserId(req.Mobile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate("comment",[][]byte{
		[]byte(req.OrderId),
		[]byte(userId),
		[]byte(req.Comment),
	})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}
