package handler

import (
	"context"

		example "sss/Register/proto/example"
	"sss/181231/utils"
		"github.com/garyburd/redigo/redis"
	"sss/181231/models"
	"crypto/sha512"
	"encoding/hex"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

func (e *Example) Register(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	cpt,err:=redis.String(conn.Do("get",req.Mobile))
	if cpt=="" || req.SmsCode!=cpt  {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)+":Invalid SMS code"
		return nil
	}

	var user models.User
	user.Name=req.Mobile
	user.Password=GetSha512Str(req.Password)
	user.Mobile=req.Mobile
	o:=orm.NewOrm()
	if _,err=o.Insert(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

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
