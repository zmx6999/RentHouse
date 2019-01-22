package handler

import (
	"context"

		example "190120/Order/proto/example"
	"190120/utils"
	"encoding/json"
	"strconv"
	"time"
	"math"
	"github.com/astaxie/beego/orm"
	"190120/models"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Add(ctx context.Context, req *example.AddRequest, rsp *example.AddResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var data map[string]interface{}
	json.Unmarshal(req.Data,&data)

	houseId,_:=strconv.Atoi(data["house_id"].(string))
	startDate:=data["start_date"].(string)
	start,_:=time.Parse("2006-01-02 15:04:05",startDate+" 00:00:00")
	endDate:=data["end_date"].(string)
	end,_:=time.Parse("2006-01-02 15:04:05",endDate+" 00:00:00")
	days:=int(math.Ceil(float64(end.Sub(start).Seconds())/60/60/24))+1

	o:=orm.NewOrm()
	var house models.House
	if err=o.QueryTable("House").RelatedSel("User").Filter("Id",houseId).One(&house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if house.User.Id==userId {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	price:=house.Price
	amount:=price*days
	status:=utils.ORDER_STATUS_WAIT_ACCEPT

	user:=models.User{Id:userId}

	order:=models.Order{
		Begin_date:start,
		End_date:end,
		Days:days,
		House_price:price,
		Amount:amount,
		Status:status,
		Credit:false,
		House:&house,
		User:&user,
	}
	loc,_:=time.LoadLocation(utils.TimeZone)
	order.Create_time=time.Now().In(loc)

	if _,err=o.Insert(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}

func (e *Example) GetOrders(ctx context.Context, req *example.GetOrdersRequest, rsp *example.GetOrdersResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o:=orm.NewOrm()
	var _orders []*models.Order
	if req.Role=="landlord" {
		var houses []*models.House
		if _,err=o.QueryTable("House").RelatedSel("User").Filter("User__Id",userId).All(&houses);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}

		houseIds:=[]int{-1}
		for _,v:=range houses{
			houseIds=append(houseIds,v.Id)
		}

		if _,err=o.QueryTable("Order").RelatedSel("House").Filter("House__Id__in",houseIds).All(&_orders);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	} else {
		if _,err=o.QueryTable("Order").RelatedSel("House").RelatedSel("User").Filter("User__Id",userId).All(&_orders);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	}

	var orders []map[string]interface{}
	for _,v:=range _orders{
		order:=v.Info()
		orders=append(orders,order)
	}
	data,_:=json.Marshal(orders)
	rsp.Data=data

	return nil
}

func (e *Example) Handle(ctx context.Context, req *example.HandleRequest, rsp *example.HandleResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var order models.Order
	o:=orm.NewOrm()
	if err=o.QueryTable("Order").RelatedSel("House").Filter("Id",req.OrderId).Filter("Status",utils.ORDER_STATUS_WAIT_ACCEPT).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	house:=order.House
	if _,err=o.LoadRelated(house,"User");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if house.User.Id!=userId {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Begin()

	if req.Action=="reject" {
		order.Status=utils.ORDER_STATUS_REJECTED
	} else {
		house.Order_count++
		if _,err=o.Update(house,"Order_count");err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
		order.Status=utils.ORDER_STATUS_WAIT_COMMENT
	}
	if _,err=o.Update(&order,"Status");err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Commit()

	return nil
}

func (e *Example) Comment(ctx context.Context, req *example.CommentRequest, rsp *example.CommentResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var order models.Order
	o:=orm.NewOrm()
	if err=o.QueryTable("Order").RelatedSel("House").RelatedSel("User").Filter("Id",req.OrderId).Filter("Status",utils.ORDER_STATUS_WAIT_COMMENT).Filter("User__Id",userId).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	order.Comment=req.Comment
	order.Status=utils.ORDER_STATUS_COMPLETE
	if _,err=o.Update(&order,"Comment","Status");err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	conn.Do("del","house_"+strconv.Itoa(order.House.Id))

	return nil
}
