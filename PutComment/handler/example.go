package handler

import (
	"context"

		example "sss/PutComment/proto/example"
	"sss/181231/utils"
	"strconv"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PutComment(ctx context.Context, req *example.Request, rsp *example.Response) error {
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
	o.Begin()

	if err=o.QueryTable("Order").RelatedSel("User").Filter("Id",orderId).Filter("Status",models.ORDER_STATUS_WAIT_COMMENT).Filter("User__Id",userId).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	order.Status=models.ORDER_STATUS_COMPLETE
	order.Comment=req.Comment
	if _,err:=o.Update(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if _,err=o.LoadRelated(&order,"House");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	house:=order.House
	house.Order_count++
	if _,err=o.Update(house);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	o.Commit()

	conn.Do("del","house_"+strconv.Itoa(house.Id))

	return nil
}
