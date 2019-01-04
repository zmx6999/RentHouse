package handler

import (
	"context"

		example "sss/PutOrder/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/orm"
	"strconv"
	"sss/181231/models"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PutOrder(ctx context.Context, req *example.Request, rsp *example.Response) error {
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

	orderId,_:=strconv.Atoi(req.OrderId)

	var order models.Order
	o:=orm.NewOrm()
	if err=o.QueryTable("Order").RelatedSel("House").Filter("Id",orderId).Filter("Status",models.ORDER_STATUS_WAIT_ACCEPT).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	house:=order.House
	if _,err=o.LoadRelated(house,"User");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	if house.User.Id!=userId {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	if req.Action=="accept" {
		order.Status=models.ORDER_STATUS_WAIT_COMMENT
	} else if req.Action=="reject" {
		order.Status=models.ORDER_STATUS_REJECTED
	}
	if _,err:=o.Update(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	return nil
}
