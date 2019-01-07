package handler

import (
	"context"

		example "190105/House/proto/example"
	"190105/utils"
	"github.com/astaxie/beego"
	"190105/models"
	"encoding/json"
	"strconv"
	"github.com/astaxie/beego/orm"
	"github.com/zmx6999/FormValidation/FormValidation"
	"time"
	"path"
	"github.com/garyburd/redigo/redis"
	"math"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) AddHouse(ctx context.Context, req *example.AddHouseRequest, rsp *example.AddHouseResponse) error {
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

	_house:=map[string]interface{}{}
	json.Unmarshal(req.Data,&_house)

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"title",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"title cannot be empty",
		},
	}
	gv:=FormValidation.GroupValidation{_house,fvs}
	_,err=gv.Validate()
	if err!=nil {
		rsp.ErrCode=utils.RECODE_PARAMERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)+":"+err.Error()
		return nil
	}

	var house models.House
	house.Title=_house["title"].(string)
	price,_:=strconv.Atoi(_house["price"].(string))
	house.Price=price*100

	areaId,_:=strconv.Atoi(_house["area_id"].(string))
	area:=models.Area{Id:areaId}
	house.Area=&area

	house.Address=_house["address"].(string)
	house.Room_count,_=strconv.Atoi(_house["room_count"].(string))
	house.Acreage,_=strconv.Atoi(_house["acreage"].(string))
	house.Unit=_house["unit"].(string)
	house.Capacity,_=strconv.Atoi(_house["capacity"].(string))
	house.Beds=_house["beds"].(string)
	deposit,_:=strconv.Atoi(_house["deposit"].(string))
	house.Deposit=deposit*100
	house.Min_days,_=strconv.Atoi(_house["min_days"].(string))
	house.Max_days,_=strconv.Atoi(_house["max_days"].(string))

	user:=models.User{Id:userId}
	house.User=&user

	loc,_:=time.LoadLocation(utils.TimeZone)
	house.Create_time=time.Now().In(loc)

	facilities:=[]*models.Facility{}
	for _,v:=range _house["facility"].([]interface{}){
		facilityId,_:=strconv.Atoi(v.(string))
		facility:=models.Facility{Id:facilityId}
		facilities=append(facilities,&facility)
	}

	o:=orm.NewOrm()
	o.Begin()

	if _,err:=o.Insert(&house);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if _,err=o.QueryM2M(&house,"Facilities").Add(facilities);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Commit()
	rsp.HouseId=int64(house.Id)

	return nil
}

func (e *Example) GetUserHouses(ctx context.Context, req *example.GetUserHousesRequest, rsp *example.GetUserHousesResponse) error {
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

	_houses:=[]*models.House{}
	o:=orm.NewOrm()
	if _,err=o.QueryTable("House").RelatedSel("Area","User").Filter("User__Id",userId).All(&_houses);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	houses:=[]map[string]interface{}{}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}
	data,_:=json.Marshal(houses)
	rsp.Data=data

	return nil
}

func (e *Example) UploadHouseImage(ctx context.Context, req *example.UploadHouseImageRequest, rsp *example.UploadHouseImageResponse) error {
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

	if int(req.FileSize)!=len(req.Data) {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o:=orm.NewOrm()
	var house models.House
	if err=o.QueryTable("House").RelatedSel("User").Filter("Id",req.HouseId).Filter("User__Id",userId).One(&house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Begin()

	if house.Index_image_url=="" {
		house.Index_image_url=fileId
		if _,err:=o.Update(&house);err!=nil {
			o.Rollback()
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.Msg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
	}

	houseImage:=models.HouseImage{Url:fileId,House:&house}
	if _,err=o.Insert(&houseImage);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	o.Commit()

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	conn.Do("del","house_"+strconv.Itoa(int(req.HouseId)))
	return nil
}

func (e *Example) GetHouseDetail(ctx context.Context, req *example.GetHouseDetailRequest, rsp *example.GetHouseDetailResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	data,err:=redis.Bytes(conn.Do("get","house_"+strconv.Itoa(int(req.HouseId))))
	if data!=nil {
		beego.Info(data)
		rsp.Data=data
		return nil
	}

	_house:=models.House{Id:int(req.HouseId)}
	o:=orm.NewOrm()
	if err=o.Read(&_house);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	o.LoadRelated(&_house,"User")
	o.LoadRelated(&_house,"Facilities")
	o.LoadRelated(&_house,"HouseImages")

	o.LoadRelated(&_house,"Orders")
	for _,order:=range _house.Orders{
		o.LoadRelated(order,"User")
	}

	house:=_house.Desc()
	data,_=json.Marshal(house)
	rsp.Data=data
	conn.Do("set","house_"+strconv.Itoa(int(req.HouseId)),data,"EX",3600)

	return nil
}

func (e *Example) GetIndexBanner(ctx context.Context, req *example.GetIndexBannerRequest, rsp *example.GetIndexBannerResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	data,err:=redis.Bytes(conn.Do("get","index_banner"))
	if data!=nil {
		beego.Info(data)
		rsp.Data=data
		return nil
	}

	_houses:=[]*models.House{}
	houses:=[]map[string]interface{}{}
	o:=orm.NewOrm()
	if _,err:=o.QueryTable("House").RelatedSel("Area","User").Limit(5).All(&_houses);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}
	data,_=json.Marshal(houses)
	conn.Do("set","index_banner",data,"EX",3600)
	rsp.Data=data

	return nil
}

func (e *Example) SearchHouse(ctx context.Context, req *example.SearchHouseRequest, rsp *example.SearchHouseResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	_houses:=[]*models.House{}
	houses:=[]map[string]interface{}{}
	o:=orm.NewOrm()
	qs:=o.QueryTable("House").RelatedSel("Area","User")
	if req.AreaId!="" {
		areaId,_:=strconv.Atoi(req.AreaId)
		qs=qs.Filter("Area__Id",areaId)
	}
	if req.StartDate!="" {
		start,_:=time.Parse("2006-01-02 15:04:05",req.StartDate+" 00:00:00")
		qs=qs.Filter("Create_time__gte",start)
	}
	if req.EndDate!="" {
		end,_:=time.Parse("2006-01-02 15:04:05",req.EndDate+" 23:59:59")
		qs=qs.Filter("Create_time__lte",end)
	}
	totalRows,_:=qs.Count()
	pageSize:=5
	totalPages:=int64(math.Ceil(float64(totalRows)/float64(pageSize)))
	page,err:=strconv.Atoi(req.Page)
	if err!=nil {
		page=1
	}
	offset:=pageSize*(page-1)
	_,err=qs.Limit(pageSize,offset).All(&_houses)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	for _,v:=range _houses{
		house:=v.Info()
		houses=append(houses,house)
	}
	data,_:=json.Marshal(houses)
	rsp.Data=data
	rsp.TotalPages=totalPages
	rsp.CurrentPage=int64(page)

	return nil
}
