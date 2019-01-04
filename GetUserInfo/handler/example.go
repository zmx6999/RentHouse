package handler

import (
	"context"

		example "sss/GetUserInfo/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
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

	var user models.User
	user.Id=userId
	o:=orm.NewOrm()
	if err:=o.Read(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	rsp.UserId=int64(user.Id)
	rsp.Name=user.Name
	rsp.Mobile=user.Mobile
	rsp.RealName=user.Real_name
	rsp.IdCard=user.Id_card
	rsp.AvatarUrl=user.Avatar_url
	return nil
}
