package handler

import (
	"context"

		example "sss/GetArea/proto/example"
	"sss/181231/utils"
	"sss/181231/models"
	"github.com/garyburd/redigo/redis"
	"encoding/gob"
	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	var areaList []models.Area
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	areaInfo,_:=redis.Bytes(conn.Do("get","area_info"))
	if areaInfo!=nil {
		decoder:=gob.NewDecoder(bytes.NewReader(areaInfo))
		decoder.Decode(&areaList)
		data:=[]*example.Response_Area{}
		for k,v:=range areaList{
			beego.Info(k,v)
			area:=&example.Response_Area{Id:int32(v.Id),Name:v.Name}
			data=append(data,area)
		}
		rsp.Data=data
		return nil
	}

	o:=orm.NewOrm()
	n,err:=o.QueryTable("Area").All(&areaList)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	if n==0 {
		if err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
			return nil
		}
	}

	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	encoder.Encode(&areaList)
	conn.Do("set","area_info",buffer.Bytes())

	data:=[]*example.Response_Area{}
	for _,v:=range areaList{
		area:=&example.Response_Area{Id:int32(v.Id),Name:v.Name}
		data=append(data,area)
	}
	rsp.Data=data
	return nil
}
