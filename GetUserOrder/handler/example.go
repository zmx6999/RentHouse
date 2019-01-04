package handler

import (
	"context"

		example "sss/GetUserOrder/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserOrder(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	userId,err:=redis.Int(conn.Do("get",req.SessionId+"_id"))
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	var _orders []*models.Order
	o:=orm.NewOrm()
	if req.Role=="landlord" {
		var houses []models.House
		_,err:=o.QueryTable("House").RelatedSel("User").Filter("User__Id",userId).All(&houses)
		if err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
			return nil
		}
		houseIds:=[]int{-1}
		for _,house:=range houses{
			houseIds=append(houseIds,house.Id)
		}
		_,err=o.QueryTable("Order").RelatedSel("House").Filter("House__Id__in",houseIds).All(&_orders)
		if err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
			return nil
		}
	} else {
		_,err=o.QueryTable("Order").RelatedSel("User").Filter("User__Id",userId).All(&_orders)
		if err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
			return nil
		}
	}
	var orders []map[string]interface{}
	for _,order:=range _orders{
		o.LoadRelated(order,"House")
		orders=append(orders,order.Info())
	}
	data,_:=json.Marshal(orders)
	rsp.Data=data
	return nil
}
