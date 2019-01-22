package handler

import (
	"context"

		example "190120/Area/proto/example"
	"github.com/gomodule/redigo/redis"
	"190120/utils"
	"190120/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.GetAreaRequest, rsp *example.GetAreaResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	data,err:=redis.Bytes(conn.Do("get","areas"))
	beego.Info(data)
	if data==nil {
		var _areaList []*models.Area
		o:=orm.NewOrm()
		if _,err=o.QueryTable("Area").All(&_areaList);err!=nil {
			rsp.ErrCode=utils.RECODE_NODATA
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}

		var areaList []map[string]interface{}
		for _,v:=range _areaList{
			areaList=append(areaList,map[string]interface{}{
				"area_id": v.Id,
				"area_name": v.Name,
			})
		}

		data,_:=json.Marshal(areaList)
		conn.Do("set","areas",data,"EX",3600)
	}
	rsp.Data=data

	return nil
}
