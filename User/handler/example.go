package handler

import (
	"context"

		example "190105/User/proto/example"
	"190105/utils"
	"github.com/afocus/captcha"
	"image/color"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"github.com/SubmailDem/submail"
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"190105/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"path"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImageCpt(ctx context.Context, req *example.GetImageCptRequest, rsp *example.GetImageCptResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	cpt:=captcha.New()
	cpt.SetFont("comic.ttf")
	cpt.SetSize(91,41)
	cpt.SetDisturbance(captcha.MEDIUM)
	cpt.SetFrontColor(color.RGBA{255,255,255,255})
	cpt.SetBkgColor(color.RGBA{255, 0, 0, 255},color.RGBA{0, 0, 255, 255},color.RGBA{0, 153, 0, 255})
	img,str:=cpt.Create(4,captcha.NUM)

	c:=img.RGBA
	rsp.Pix=c.Pix
	rsp.Stride=int64(c.Stride)
	rsp.Min=&example.GetImageCptResponse_Point{X:int64(c.Rect.Min.X),Y:int64(c.Rect.Min.Y)}
	rsp.Max=&example.GetImageCptResponse_Point{X:int64(c.Rect.Max.X),Y:int64(c.Rect.Max.Y)}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	beego.Info(req.Uuid,str)
	_,err=conn.Do("set",req.Uuid,str,"EX",3600)
	if err!=nil {
		beego.Error(err)
	}

	return nil
}

func (e *Example) GetSmsCpt(ctx context.Context, req *example.GetSmsCptRequest, rsp *example.GetSmsCptResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	beego.Info(utils.RedisHost+":"+utils.RedisPort)
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	cpt,err:=redis.String(conn.Do("get",req.Uuid))
	beego.Info(cpt)
	if cpt=="" || req.Text!=cpt {
		rsp.ErrCode=utils.RECODE_PARAMERR
		rsp.Msg="Invalid captcha"
		return nil
	}

	messageconfig := make(map[string]string)
	messageconfig["appid"] = "29672"
	messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	messageconfig["signtype"] = "md5"

	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, "NQ1J94")
	code:=rand.New(rand.NewSource(time.Now().UnixNano())).Intn(8999)+1001
	submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(code))
	fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))

	conn.Do("set",req.Mobile,code,"EX",3600)

	return nil
}

func (e *Example) Register(ctx context.Context, req *example.RegisterRequest, rsp *example.RegisterResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()
	cpt,err:=redis.String(conn.Do("get",req.Mobile))
	beego.Info(cpt)
	if cpt=="" || req.Text!=cpt {
		rsp.ErrCode=utils.RECODE_PARAMERR
		rsp.Msg="Invalid sms captcha"
		return nil
	}

	user:=models.User{Mobile:req.Mobile,Name:req.Mobile,Password:utils.GetSha512Str(req.Password)}
	o:=orm.NewOrm()
	if _,err:=o.Insert(&user);err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	sessionId:=utils.GetSha512Str(req.Mobile+req.Password)
	conn.Do("set",sessionId+"_user_id",user.Id,"EX",3600)
	conn.Do("set",sessionId+"_mobile",user.Mobile,"EX",3600)
	conn.Do("set",sessionId+"_name",user.Name,"EX",3600)

	rsp.SessionId=sessionId

	return nil
}

func (e *Example) Login(ctx context.Context, req *example.LoginRequest, rsp *example.LoginResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	user:=models.User{Mobile:req.Mobile}
	o:=orm.NewOrm()
	if err:=o.Read(&user,"Mobile");err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	if utils.GetSha512Str(req.Password)!=user.Password {
		rsp.ErrCode=utils.RECODE_PWDERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	sessionId:=utils.GetSha512Str(req.Mobile+req.Password)
	conn.Do("set",sessionId+"_user_id",user.Id,"EX",3600)
	conn.Do("set",sessionId+"_mobile",user.Mobile,"EX",3600)
	conn.Do("set",sessionId+"_name",user.Name,"EX",3600)

	rsp.SessionId=sessionId

	return nil
}

func (e *Example) Logout(ctx context.Context, req *example.LogoutRequest, rsp *example.LogoutResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	conn.Do("del",req.SessionId+"_user_id")
	conn.Do("del",req.SessionId+"_mobile")
	conn.Do("del",req.SessionId+"_name")

	return nil
}

func (e *Example) GetUserInfo(ctx context.Context, req *example.GetUserInfoRequest, rsp *example.GetUserInfoResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	_user:=models.User{Id:userId}
	o:=orm.NewOrm()
	if err=o.Read(&_user);err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	user:=map[string]interface{}{
		"user_id":_user.Id,
		"name":_user.Name,
		"mobile":_user.Mobile,
		"real_name":_user.Real_name,
		"id_card":_user.Id_card,
		"avatar_url":_user.Avatar_url,
	}
	data,_:=json.Marshal(user)
	rsp.Data=data

	return nil
}

func (e *Example) Rename(ctx context.Context, req *example.RenameRequest, rsp *example.RenameResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	user:=models.User{Id:userId,Name:req.NewName}
	o:=orm.NewOrm()
	if _,err:=o.Update(&user,"Name");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	defer conn.Close()

	conn.Do("set",req.SessionId+"_user_id",user.Id,"EX",3600)
	conn.Do("set",req.SessionId+"_mobile",user.Mobile,"EX",3600)
	conn.Do("set",req.SessionId+"_name",user.Name,"EX",3600)

	return nil
}

func (e *Example) Auth(ctx context.Context, req *example.AuthRequest, rsp *example.AuthResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	user:=models.User{Id:userId,Real_name:req.RealName,Id_card:req.IdCard}
	o:=orm.NewOrm()
	if _,err:=o.Update(&user,"Real_name","Id_card");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}

func (e *Example) UploadAvatar(ctx context.Context, req *example.UploadAvatarRequest, rsp *example.UploadAvatarResponse) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.ErrCode)

	sessionInfo,err:=utils.GetSession(req.SessionId)
	beego.Info(sessionInfo)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}
	userId:=sessionInfo["user_id"].(int)

	if int(req.FileSize)!=len(req.Data) {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		rsp.ErrCode=utils.RECODE_IOERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	user:=models.User{Id:userId,Avatar_url:fileId}
	o:=orm.NewOrm()
	if _,err:=o.Update(&user,"Avatar_url");err!=nil {
		rsp.ErrCode=utils.RECODE_DBERR
		rsp.Msg=utils.RecodeText(rsp.ErrCode)
		return nil
	}

	return nil
}
