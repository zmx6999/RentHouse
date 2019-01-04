package handler

import (
	"context"

		example "sss/PostHouses/proto/example"
	"sss/181231/utils"
		"encoding/json"
	"sss/181231/models"
	"strconv"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouse(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	houseMap:=map[string]interface{}{}
	err:=json.Unmarshal(req.Body,&houseMap)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	var house models.House
	house.Title=houseMap["title"].(string)
	price,_:=strconv.ParseFloat(houseMap["price"].(string),64)
	house.Price=int(price*100)
	house.Address=houseMap["address"].(string)
	house.Room_count,_=strconv.Atoi(houseMap["room_count"].(string))
	house.Acreage,_=strconv.Atoi(houseMap["acreage"].(string))
	house.Unit=houseMap["unit"].(string)
	house.Capacity,_=strconv.Atoi(houseMap["capacity"].(string))
	house.Beds=houseMap["beds"].(string)
	deposit,_:=strconv.ParseFloat(houseMap["deposit"].(string),64)
	house.Deposit=int(deposit*100)
	house.Max_days,_=strconv.Atoi(houseMap["max_days"].(string))
	house.Min_days,_=strconv.Atoi(houseMap["min_days"].(string))

	loc,_:=time.LoadLocation("Asia/Saigon")
	house.Create_time=time.Now().In(loc)

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

	user:=models.User{Id:userId}
	house.User=&user

	areaId,_:=strconv.Atoi(houseMap["area_id"].(string))
	area:=models.Area{Id:areaId}
	house.Area=&area

	facilities:=[]*models.Facility{}
	for _,_fid:=range houseMap["facilities"].([]interface{}){
		fid,_:=strconv.Atoi(_fid.(string))
		facility:=models.Facility{Id:fid}
		facilities=append(facilities,&facility)
	}

	o:=orm.NewOrm()
	o.Begin()

	if _,err:=o.Insert(&house);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if _,err=o.QueryM2M(&house,"Facilities").Add(facilities);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	o.Commit()
	rsp.HouseId=int64(house.Id)
	return nil
}
