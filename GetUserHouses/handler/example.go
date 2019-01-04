package handler

import (
	"context"

		example "sss/GetUserHouses/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

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

	var houses []*models.House
	o:=orm.NewOrm()
	if _,err:=o.QueryTable("House").RelatedSel("User","Area").Filter("User__Id",userId).All(&houses);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	data,err:=json.Marshal(houses)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	rsp.Data=data
	return nil
}
