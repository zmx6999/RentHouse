package handler

import (
	"context"

		example "190105/Area/proto/example"
	"190105/utils"
	"190105/models"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"encoding/json"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.GetAreaRequest, rsp *example.GetAreaResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	data,_:=redis.Bytes(conn.Do("get","areas"))
	if data!=nil {
		beego.Info(data)
		rsp.Data=data
		return nil
	}

	var _areas []*models.Area
	o:=orm.NewOrm()
	if _,err:=o.QueryTable("Area").All(&_areas);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	beego.Info(_areas)

	areas:=[]map[string]interface{}{}
	for _,v:=range _areas{
		area:=map[string]interface{}{
			"area_id":v.Id,
			"area_name":v.Name,
		}
		areas=append(areas,area)
	}

	data,_=json.Marshal(areas)
	conn.Do("set","areas",data)
	rsp.Data=data
	return nil
}
