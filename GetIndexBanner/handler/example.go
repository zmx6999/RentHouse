package handler

import (
	"context"

		example "sss/GetIndexBanner/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetIndexBanner(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	data,_:=redis.Bytes(conn.Do("get","index_banner"))
	if data!=nil {
		beego.Info(data)
		rsp.Data=data
		return nil
	}

	var _houses []*models.House
	houses:=[]map[string]interface{}{}
	o:=orm.NewOrm()
	if _,err=o.QueryTable("House").Limit(5).All(&_houses);err==nil {
		for _,house:=range _houses{
			o.LoadRelated(house,"Area")
			o.LoadRelated(house,"User")
			houses=append(houses,house.Info())
		}
	}
	data,_=json.Marshal(houses)
	conn.Do("set","index_banner",data)
	return nil
}
