package handler

import (
	"context"

		example "sss/Login/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"crypto/sha512"
	"encoding/hex"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Login(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	var user models.User
	o:=orm.NewOrm()
	if err:=o.QueryTable("User").Filter("Mobile",req.Mobile).One(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if GetSha512Str(req.Password)!=user.Password {
		rsp.ErrCode=utils.RECODE_PWDERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	rsp.SessionId=GetSha512Str(req.Mobile+req.Password)
	conn.Do("set",rsp.SessionId+"_name",req.Mobile,"EX",3600)
	conn.Do("set",rsp.SessionId+"_id",user.Id,"EX",3600)
	conn.Do("set",rsp.SessionId+"_mobile",req.Mobile,"EX",3600)
	return nil
}

func GetSha512Str(x string) string {
	h:=sha512.New()
	h.Write([]byte(x))
	return hex.EncodeToString(h.Sum(nil))
}
