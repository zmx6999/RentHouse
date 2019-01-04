package handler

import (
	"context"

		example "sss/PostHouseImage/proto/example"
	myUtils "sss/PostHouseImage/utils"
	"sss/181231/utils"
	"path"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/orm"
	"sss/181231/models"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouseImage(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	if len(req.Image)!=int(req.FileSize) {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=myUtils.UploadFile(req.Image,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)+":"+err.Error()
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
	_,err=conn.Do("del","house_"+strconv.Itoa(int(req.HouseId)))
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	o:=orm.NewOrm()
	o.Begin()

	var house models.House
	if err=o.QueryTable("House").RelatedSel("User").Filter("User__Id",userId).Filter("Id",req.HouseId).One(&house);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if house.Index_image_url=="" {
		house.Index_image_url=fileId
	}

	if _,err:=o.Update(&house);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	houseImage:=models.HouseImage{Url:fileId,House:&house}
	if _,err=o.Insert(&houseImage);err!=nil {
		o.Rollback()
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	o.Commit()
	rsp.Url=fileId
	return nil
}
