package handler

import (
	"context"

		example "190120/House/proto/example"
	"encoding/json"
	"strconv"
	"190120/models"
	"190120/utils"
	"github.com/astaxie/beego/orm"
	"reflect"
	"time"
	"path"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"strings"
	"math"
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

	title:=data["title"].(string)
	price,_:=strconv.Atoi(data["price"].(string))
	price*=100
	areaId,_:=strconv.Atoi(data["area_id"].(string))
	address:=data["address"].(string)
	roomCount,_:=strconv.Atoi(data["room_count"].(string))
	acreage,_:=strconv.Atoi(data["acreage"].(string))
	unit:=data["unit"].(string)
	capacity,_:=strconv.Atoi(data["capacity"].(string))
	beds:=data["beds"].(string)
	deposit,_:=strconv.Atoi(data["deposit"].(string))
	deposit*=100
	minDays,_:=strconv.Atoi(data["min_days"].(string))
	maxDays,_:=strconv.Atoi(data["max_days"].(string))

	user:=models.User{Id:userId}
	area:=models.Area{Id:areaId}
	house:=models.House{
		Title:title,
		Price:price,
		Address:address,
		Room_count:roomCount,
		Acreage:acreage,
		Unit:unit,
		Capacity:capacity,
		Beds:beds,
		Deposit:deposit,
		Min_days:minDays,
		Max_days:maxDays,
		Area:&area,
		User:&user,
	}
	loc,_:=time.LoadLocation(utils.TimeZone)
	house.Create_time=time.Now().In(loc)

	o:=orm.NewOrm()
	o.Begin()

	if _,err:=o.Insert(&house);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if data["facility"]!=nil && strings.Replace(reflect.TypeOf(data["facility"]).String()," ","",-1)=="[]interface{}" {
		var facilities []*models.Facility
		for _,v:=range data["facility"].([]interface{}){
			facilityId,_:=strconv.Atoi(v.(string))
			facility:=models.Facility{Id:facilityId}
			facilities=append(facilities,&facility)
		}

		if _,err=o.QueryM2M(&house,"Facilities").Add(facilities);err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	}

	o.Commit()
	return nil
}

func (e *Example) GetHouses(ctx context.Context, req *example.GetHousesRequest, rsp *example.GetHousesResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var _houses []*models.House
	o:=orm.NewOrm()
	if _,err=o.QueryTable("House").RelatedSel("Area").RelatedSel("User").Filter("User__Id",userId).All(&_houses);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var houses []map[string]interface{}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}

	data,_:=json.Marshal(houses)
	rsp.Data=data

	return nil
}

func (e *Example) UploadImage(ctx context.Context, req *example.UploadImageRequest, rsp *example.UploadImageResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if len(req.Data)!=int(req.FileSize) {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var house models.House
	o:=orm.NewOrm()
	if err=o.QueryTable("House").RelatedSel("User").Filter("Id",req.HouseId).Filter("User__Id",userId).One(&house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Begin()

	if house.Index_image_url=="" {
		house.Index_image_url=fileId
		if _,err:=o.Update(&house,"Index_image_url");err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	}

	houseImage:=models.HouseImage{House:&house,Url:fileId}
	if _,err=o.Insert(&houseImage);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Commit()

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	conn.Do("del","house_"+strconv.Itoa(int(req.HouseId)))

	return nil
}

func (e *Example) GetHouseDetail(ctx context.Context, req *example.GetHouseDetailRequest, rsp *example.GetHouseDetailResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	data,err:=redis.Bytes(conn.Do("get","house_"+strconv.Itoa(int(req.HouseId))))
	beego.Info(data)
	if data==nil {
		var _house models.House
		o:=orm.NewOrm()
		if err:=o.QueryTable("House").RelatedSel("User").Filter("Id",req.HouseId).One(&_house);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
		o.LoadRelated(&_house,"HouseImages")
		o.LoadRelated(&_house,"Facilities")
		var orders []*models.Order
		o.QueryTable("Order").RelatedSel("User").RelatedSel("House").Filter("House__Id",_house.Id).Filter("Status",utils.ORDER_STATUS_COMPLETE).All(&orders)
		_house.Orders=orders
		house:=_house.Desc()
		data,_=json.Marshal(house)
		conn.Do("set","house_"+strconv.Itoa(int(req.HouseId)),data,"EX",3600)
	}
	rsp.Data=data

	return nil
}

func (e *Example) GetIndexBanner(ctx context.Context, req *example.GetIndexBannerRequest, rsp *example.GetIndexBannerResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	var _houses []*models.House
	o:=orm.NewOrm()
	if _,err:=o.QueryTable("House").RelatedSel("Area").RelatedSel("User").Limit(5).All(&_houses);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var houses []map[string]interface{}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}

	data,_:=json.Marshal(houses)
	rsp.Data=data

	return nil
}

func (e *Example) Search(ctx context.Context, req *example.SearchRequest, rsp *example.SearchResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	var _houses []*models.House
	o:=orm.NewOrm()
	q:=o.QueryTable("House").RelatedSel("Area").RelatedSel("User")

	areaId,err:=strconv.Atoi(req.AreaId)
	if err==nil {
		q=q.Filter("Area__Id",areaId)
	}

	if req.StartDate!="" {
		q=q.Filter("Create_time__gte",req.StartDate+" 00:00:00")
	}

	if req.EndDate!="" {
		q=q.Filter("Create_time__lte",req.EndDate+" 23:59:59")
	}

	total,_:=q.Count()
	pageSize:=2
	totalPage:=int(math.Ceil(float64(total)/float64(pageSize)))

	page,err:=strconv.Atoi(req.Page)
	if err!=nil || page<1 || page>totalPage {
		page=1
	}
	offset:=pageSize*(page-1)

	if _,err:=q.Limit(pageSize,offset).All(&_houses);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var houses []map[string]interface{}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}

	data,_:=json.Marshal(houses)
	rsp.Data=data
	rsp.TotalPage=int64(totalPage)
	rsp.Page=int64(page)

	return nil
}
