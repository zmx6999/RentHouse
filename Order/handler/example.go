package handler

import (
	"context"

		example "190105/Order/proto/example"
	"190105/utils"
	"github.com/astaxie/beego"
			"strconv"
		"github.com/astaxie/beego/orm"
	"190105/models"
	"time"
	"math"
	"encoding/json"
	"github.com/zmx6999/FormValidation/FormValidation"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) AddOrder(ctx context.Context, req *example.AddOrderRequest, rsp *example.AddOrderResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	data:=map[string]interface{}{}
	json.Unmarshal(req.Data,&data)

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"house_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"house id cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"start_date",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"start date cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"end date cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"EndDate",
			ValidMethodArgs:[]interface{}{data["start_date"]},
			Trim:true,
			ErrMsg:"end date should be later than start date",
		},
	}
	gv:=FormValidation.GroupValidation{data,fvs}
	_,err=gv.Validate()
	if err!=nil {
		rsp.ErrCode=utils.RECODE_PARAMERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)+":"+err.Error()
		return nil
	}

	houseId,_:=strconv.Atoi(data["house_id"].(string))
	var house models.House
	o:=orm.NewOrm()
	if err=o.QueryTable("House").RelatedSel("User").Filter("Id",houseId).One(&house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	if userId==house.User.Id {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	start,_:=time.Parse("2006-01-02 15:04:05",data["start_date"].(string)+" 00:00:00")
	end,_:=time.Parse("2006-01-02 15:04:05",data["end_date"].(string)+" 00:00:00")
	days:=int(math.Ceil(end.Sub(start).Hours()/24))

	loc,_:=time.LoadLocation(utils.TimeZone)
	createTime:=time.Now().In(loc)

	user:=models.User{Id:userId}

	order:=models.Order{
		Begin_date:start,
		End_date:end,
		Days:days,
		House_price:house.Price,
		Amount:house.Price*days,
		Status:utils.ORDER_STATUS_WAIT_ACCEPT,
		Create_time:createTime,
		Credit:false,
		User:&user,
		House:&house,
	}
	if _,err:=o.Insert(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	rsp.OrderId=int64(order.Id)

	return nil
}

func (e *Example) GetOrders(ctx context.Context, req *example.GetOrdersRequest, rsp *example.GetOrdersResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	var _orders []*models.Order
	o:=orm.NewOrm()
	if req.Role=="landlord" {
		var houses []*models.House
		if _,err=o.QueryTable("House").RelatedSel("User").Filter("User__Id",userId).All(&houses);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
		houseIds:=[]int{-1}
		for _,v:=range houses{
			houseIds=append(houseIds,v.Id)
		}
		if _,err=o.QueryTable("Order").RelatedSel("House").Filter("House__Id__in",houseIds).OrderBy("-Create_time").All(&_orders);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	} else {
		if _,err=o.QueryTable("Order").RelatedSel("User","House").Filter("User__Id",userId).OrderBy("-Create_time").All(&_orders);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	}
	orders:=[]map[string]interface{}{}
	for _,v:=range _orders{
		orders=append(orders,v.Info())
	}
	data,_:=json.Marshal(orders)
	rsp.Data=data

	return nil
}

func (e *Example) HandleOrder(ctx context.Context, req *example.HandleOrderRequest, rsp *example.HandleOrderResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	var order models.Order
	o:=orm.NewOrm()
	if err=o.QueryTable("Order").RelatedSel("House").Filter("Id",req.OrderId).Filter("Status",utils.ORDER_STATUS_WAIT_ACCEPT).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	house:=order.House
	if _,err=o.LoadRelated(house,"User");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if house.User.Id!=userId {
		rsp.ErrCode=utils.RECODE_ROLEERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if req.Action=="reject" {
		order.Status=utils.ORDER_STATUS_REJECTED
		if _,err=o.Update(&order);err!=nil {
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	} else {
		o.Begin()

		order.Status=utils.ORDER_STATUS_WAIT_COMMENT
		if _,err=o.Update(&order);err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}

		house.Order_count++
		if _,err=o.Update(house);err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}

		o.Commit()
	}

	return nil
}

func (e *Example) CommentOrder(ctx context.Context, req *example.CommentOrderRequest, rsp *example.CommentOrderResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	var order models.Order
	o:=orm.NewOrm()
	if err=o.QueryTable("Order").RelatedSel("User").RelatedSel("House").Filter("Id",req.OrderId).Filter("Status",utils.ORDER_STATUS_WAIT_COMMENT).Filter("User__Id",userId).One(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	order.Status=utils.ORDER_STATUS_COMPLETE
	order.Comment=req.Comment
	if _,err=o.Update(&order);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	conn.Do("del","house_"+strconv.Itoa(int(order.House.Id)))

	return nil
}
