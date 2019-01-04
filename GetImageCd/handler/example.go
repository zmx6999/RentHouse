package handler

import (
	"context"

		example "sss/GetImageCd/proto/example"
	"sss/181231/utils"
	"github.com/afocus/captcha"
	"image/color"
	"github.com/garyburd/redigo/redis"
		)

type Example struct{}

func (e *Example) GetImageCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	cpt:=captcha.New()
	cpt.SetFont("comic.ttf")
	cpt.SetSize(91,41)
	cpt.SetDisturbance(captcha.MEDIUM)
	cpt.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cpt.SetBkgColor(color.RGBA{255, 0, 0, 255},color.RGBA{0, 0, 255, 255},color.RGBA{0, 153, 0, 255})
	img,str:=cpt.Create(4,captcha.NUM)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	_,err=conn.Do("set",req.Uuid,str,"EX",3600)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	b:=*img
	c:=*(b.RGBA)
	rsp.Pix=[]byte(c.Pix)
	rsp.Stride=int64(c.Stride)
	rsp.Max=&example.Response_Point{X:int64(c.Rect.Max.X),Y:int64(c.Rect.Max.Y)}
	rsp.Min=&example.Response_Point{X:int64(c.Rect.Min.X),Y:int64(c.Rect.Min.Y)}
	return nil
}
