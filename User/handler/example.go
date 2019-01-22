package handler

import (
	"context"

	example "190120/User/proto/example"
	"github.com/afocus/captcha"
	"image/color"
	"190120/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/SubmailDem/submail"
	"fmt"
	"math/rand"
	"strconv"
	"github.com/astaxie/beego/orm"
	"190120/models"
	"encoding/json"
	"path"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetCaptcha(ctx context.Context, req *example.GetCaptchaRequest, rsp *example.GetCaptchaResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	ca:=captcha.New()
	if err:=ca.SetFont("comic.ttf");err!=nil {
		return err
	}
	ca.SetSize(91,41)
	ca.SetDisturbance(captcha.MEDIUM)
	ca.SetFrontColor(color.RGBA{255,255,255,255})
	ca.SetBkgColor(color.RGBA{255, 0, 0, 255},color.RGBA{0, 0, 255, 255},color.RGBA{0, 153, 0, 255})
	img,str:=ca.Create(4,captcha.NUM)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	_,err=conn.Do("set","captcha_"+req.Uuid,str,"EX",3600)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	b:=*img
	c:=*(b.RGBA)
	rsp.Pix=c.Pix
	rsp.Stride=int64(c.Stride)
	rsp.Min=&example.GetCaptchaResponse_Point{X:int64(c.Rect.Min.X),Y:int64(c.Rect.Min.Y)}
	rsp.Max=&example.GetCaptchaResponse_Point{X:int64(c.Rect.Max.X),Y:int64(c.Rect.Max.Y)}
	return nil
}

func (e *Example) GetSmsCaptcha(ctx context.Context, req *example.GetSmsCaptchaRequest, rsp *example.GetSmsCaptchaResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	str,err:=redis.String(conn.Do("get","captcha_"+req.Uuid))
	if str=="" || req.Captcha!=str {
		rsp.ErrCode=utils.RECODE_SMSERR
		rsp.ErrMsg="Invalid captcha"
		return nil
	}

	messageconfig := make(map[string]string)
	messageconfig["appid"] = "29672"
	messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	messageconfig["signtype"] = "md5"

	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, "NQ1J94")
	code:=rand.Intn(8999)+1001
	submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(code))
	fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))

	_,err=conn.Do("set","sms_"+req.Mobile,strconv.Itoa(code),"EX",3600)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}

func (e *Example) Register(ctx context.Context, req *example.RegisterRequest, rsp *example.RegisterResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	str,err:=redis.String(conn.Do("get","sms_"+req.Mobile))
	if str=="" || req.SmsCaptcha!=str {
		rsp.ErrCode=utils.RECODE_SMSERR
		rsp.ErrMsg="Invalid captcha"
		return nil
	}

	o:=orm.NewOrm()
	user:=models.User{Mobile:req.Mobile}
	if err=o.Read(&user,"Mobile");err==nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg="Mobile has already registered"
		return nil
	}

	user.Password=utils.Sha512Str(req.Password)
	user.Name=user.Mobile
	if _,err=o.Insert(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	sessionId:=utils.Sha512Str(req.Mobile+req.Password)
	conn.Do("set",sessionId+"_user_id",user.Id,"EX",3600)
	conn.Do("set",sessionId+"_user_name",user.Name,"EX",3600)
	conn.Do("set",sessionId+"_mobile",user.Mobile,"EX",3600)
	rsp.SessionId=sessionId

	return nil
}

func (e *Example) Login(ctx context.Context, req *example.LoginRequest, rsp *example.LoginResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	user:=models.User{Mobile:req.Mobile}
	o:=orm.NewOrm()
	if err:=o.Read(&user,"Mobile");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg="User doesn't exit"
		return nil
	}

	if utils.Sha512Str(req.Password)!=user.Password {
		rsp.ErrCode=utils.RECODE_LOGINERR
		rsp.ErrMsg="Password error"
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	sessionId:=utils.Sha512Str(req.Mobile+req.Password)
	conn.Do("set",sessionId+"_user_id",user.Id,"EX",3600)
	conn.Do("set",sessionId+"_user_name",user.Name,"EX",3600)
	conn.Do("set",sessionId+"_mobile",user.Mobile,"EX",3600)
	rsp.SessionId=sessionId

	return nil
}

func (e *Example) Logout(ctx context.Context, req *example.LogoutRequest, rsp *example.LogoutResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	conn.Do("del",req.SessionId+"_user_id")
	conn.Do("del",req.SessionId+"_user_name")
	conn.Do("del",req.SessionId+"_mobile")

	return nil
}

func (e *Example) Info(ctx context.Context, req *example.InfoRequest, rsp *example.InfoResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	user:=models.User{Id:userId}
	o:=orm.NewOrm()
	if err=o.Read(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	data,_:=json.Marshal(user)
	rsp.Data=data

	return nil
}

func (e *Example) Avatar(ctx context.Context, req *example.AvatarRequest, rsp *example.AvatarResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	if len(req.Data)!=int(req.FileSize) {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DATAERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	user:=models.User{Id:userId,Avatar_url:fileId}
	o:=orm.NewOrm()
	if _,err=o.Update(&user,"Avatar_url");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	var houses []*models.House
	o.QueryTable("House").RelatedSel("User").Filter("User__Id",userId).All(&houses)
	if len(houses)>0 {
		conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
		if err!=nil {
			rsp.ErrCode=utils.RECODE_DBERR
			rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
			return nil
		}
		defer conn.Close()

		for _,v:=range houses{
			conn.Do("del","house_"+strconv.Itoa(v.Id))
		}
	}

	return nil
}

func (e *Example) UpdateUserName(ctx context.Context, req *example.UpdateUserNameRequest, rsp *example.UpdateUserNameResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	user:=models.User{Id:userId,Name:req.UserName}
	o:=orm.NewOrm()
	if _,err=o.Update(&user,"Name");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}

func (e *Example) Auth(ctx context.Context, req *example.AuthRequest, rsp *example.AuthResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)

	userId,err:=utils.GetUserId(req.SessionId)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_SESSIONERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	user:=models.User{Id:userId,Real_name:req.RealName,Id_card:req.IdCard}
	o:=orm.NewOrm()
	if _,err=o.Update(&user,"Real_name","Id_card");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.ErrMsg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}
