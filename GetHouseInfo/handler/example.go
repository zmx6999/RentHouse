package handler

import (
	"context"

		example "sss/GetHouseInfo/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouseInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
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
	data,_:=redis.Bytes(conn.Do("get","house_"+strconv.Itoa(int(req.HouseId))))
	if data!=nil {
		rsp.UserId=int64(userId)
		rsp.Data=data
		return nil
	}

	house:=models.House{Id:int(req.HouseId)}
	o:=orm.NewOrm()
	o.Read(&house)
	o.LoadRelated(&house,"User")
	o.LoadRelated(&house,"Facilities")
	o.LoadRelated(&house,"HouseImages")
	data,_=json.Marshal(house)
	conn.Do("set","house_"+strconv.Itoa(int(req.HouseId)),data,"EX",3600)
	rsp.UserId=int64(userId)
	rsp.Data=data
	return nil
}
