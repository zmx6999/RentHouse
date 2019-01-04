package handler

import (
	"context"

		example "sss/Logout/proto/example"
	"sss/181231/utils"
	"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Logout(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	conn.Do("del",req.SessionId+"_name")
	conn.Do("del",req.SessionId+"_id")
	conn.Do("del",req.SessionId+"_mobile")
	return nil
}
