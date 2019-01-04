package handler

import (
	"context"

		example "sss/GetSession/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

func (e *Example) GetSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_SESSIONERR
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	name,err:=redis.String(conn.Do("get",req.SessionId+"_name"))
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	if name!="" {
		rsp.ErrCode=utils.RECODE_OK
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		rsp.Data=name
	}
	return nil
}
