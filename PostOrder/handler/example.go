package handler

import (
	"context"

		example "sss/PostOrder/proto/example"
	"sss/181231/utils"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"strconv"
	"github.com/astaxie/beego/orm"
	"time"
	"math"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostOrder(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	orderMap:=map[string]interface{}{}
	err:=json.Unmarshal(req.Data,&orderMap)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if (orderMap["house_id"]==nil || orderMap["house_id"].(string)=="") || (orderMap["begin_date"]==nil || orderMap["begin_date"].(string)=="") || (orderMap["end_date"]==nil || orderMap["end_date"].(string)=="") {
		rsp.ErrCode=utils.RECODE_REQERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

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

	houseId,_:=strconv.Atoi(orderMap["house_id"].(string))
	house:=models.House{Id:houseId}
	o:=orm.NewOrm()
	if err=o.Read(&house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	o.LoadRelated(&house,"User")
	if userId==house.User.Id {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	beginDate,_:=time.Parse("2006-01-02 15:04:05",orderMap["begin_date"].(string)+" 00:00:00")
	endDate,_:=time.Parse("2006-01-02 15:04:05",orderMap["end_date"].(string)+" 00:00:00")
	if endDate.Before(beginDate) {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.ErrMsg="Begin date should not be later than end date"
		return nil
	}

	var order models.Order
	order.User=&models.User{Id:userId}
	order.House=&house
	order.Begin_date=beginDate
	order.End_date=endDate

	days:=int(math.Ceil(endDate.Sub(beginDate).Hours()/24))
	order.Days=days
	order.House_price=house.Price
	order.Amount=house.Price*days
	order.Status=models.ORDER_STATUS_WAIT_ACCEPT

	loc,_:=time.LoadLocation("Asia/Saigon")
	order.Create_time=time.Now().In(loc)

	order.Credit=false

	if _,err=o.Insert(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	rsp.OrderId=int64(order.Id)
	return nil
}
