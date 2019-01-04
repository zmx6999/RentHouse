package handler

import (
	"context"

		example "sss/PostAvatar/proto/example"
	"sss/181231/utils"
	myUtils "sss/PostAvatar/utils"
	"path"
	"sss/181231/models"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/orm"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	if len(req.Avatar)!=int(req.FileSize) {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
	}

	ext:=path.Ext(req.Filename)
	fileId,err:=myUtils.UploadFile(req.Avatar,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)+":"+err.Error()
		return nil
	}

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

	user:=models.User{Id:userId,Avatar_url:fileId}
	o:=orm.NewOrm()
	if _,err=o.Update(&user,"Avatar_url");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	rsp.AvatarUrl=fileId
	return nil
}
