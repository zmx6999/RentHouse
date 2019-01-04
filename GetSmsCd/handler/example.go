package handler

import (
	"context"

		example "sss/GetSmsCd/proto/example"
	"sss/181231/utils"
	"github.com/astaxie/beego/orm"
	"sss/181231/models"
		"github.com/garyburd/redigo/redis"
	"github.com/SubmailDem/submail"
	"fmt"
	"math/rand"
	"time"
	"strconv"
)

type Example struct{}

func (e *Example) GetSmsCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	o:=orm.NewOrm()
	user:=models.User{Mobile:req.Mobile}
	err:=o.Read(&user,"Mobile")
	if err==nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)+":User already exists"
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	cpt,err:=redis.String(conn.Do("get",req.Uuid))
	if cpt=="" || req.Text!=cpt  {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)+":Invalid Captcha"
		return nil
	}

	n:=rand.New(rand.NewSource(time.Now().UnixNano())).Intn(8999)+1001
	code:=strconv.Itoa(n)

	messageconfig := make(map[string]string)
	messageconfig["appid"] = "29672"
	messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	messageconfig["signtype"] = "md5"
	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, "NQ1J94")
	submail.MessageXSendAddVar(messagexsend, "code", code)
	fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))

	_,err=conn.Do("set",req.Mobile,code,"EX",3600)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	return nil
}
