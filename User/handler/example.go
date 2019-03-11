package handler

import (
	"context"

		example "190303/user/proto/example"
	"190303/utils"
	"190303/models"
	"github.com/afocus/captcha"
	"image/color"
	"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/SubmailDem/submail"
	"strconv"
	"math/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	crypto_rand "crypto/rand"
			"encoding/hex"
	"path"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GenerateUserId(ctx context.Context, req *example.GenerateUserIdRequest, rsp *example.GenerateUserIdResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"generateUserId",[][]byte{})
	if err!=nil {
		return err
	}
	rsp.UserId=string(data)

	return nil
}

func (e *Example) Captcha(ctx context.Context, req *example.CaptchaRequest, rsp *example.CaptchaResponse) error {
	capt:=captcha.New()

	err:=capt.SetFont("comic.ttf")
	if err!=nil {
		return err
	}
	capt.SetDisturbance(captcha.MEDIUM)
	capt.SetFrontColor(color.RGBA{255,255,255,255})
	capt.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	capt.SetSize(90,40)

	img,str:=capt.Create(4,captcha.NUM)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	uKey:="captcha_"+req.UserId
	_,err=conn.Do("set",uKey,str,"EX",3600)
	if err!=nil {
		return err
	}

	rgba:=img.RGBA
	rsp.Pix=rgba.Pix
	rsp.Stride=int64(rgba.Stride)
	rsp.Max=&example.CaptchaResponse_Point{X:int64(rgba.Rect.Max.X),Y:int64(rgba.Rect.Max.Y)}
	rsp.Min=&example.CaptchaResponse_Point{X:int64(rgba.Rect.Min.X),Y:int64(rgba.Rect.Min.Y)}

	return nil
}

func (e *Example) SmsCaptcha(ctx context.Context, req *example.SmsCaptchaRequest, rsp *example.SmsCaptchaResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	uKey:="captcha_"+req.UserId
	str,err:=redis.String(conn.Do("get",uKey))
	if err!=nil {
		return err
	}

	if str=="" || req.Captcha!=str {
		return errors.New("invalid captcha")
	}

	code:=strconv.Itoa(rand.Intn(9000)+1000)

	messageconfig := make(map[string]string)
	messageconfig["appid"] = utils.MessageAppId
	messageconfig["appkey"] = utils.MessageAppKey
	messageconfig["signtype"] = "md5"

	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, utils.MessageProject)
	submail.MessageXSendAddVar(messagexsend, "code", code)
	submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)

	mKey:="sms_"+req.Mobile
	_,err=conn.Do("set",mKey,code,"EX",3600)
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) Register(ctx context.Context, req *example.RegisterRequest, rsp *example.RegisterResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	mKey:="sms_"+req.Mobile
	str,err:=redis.String(conn.Do("get",mKey))
	if err!=nil {
		return err
	}

	if str=="" || req.SmsCaptcha!=str {
		return errors.New("invalid sms captcha")
	}

	privateKey,err:=ecdsa.GenerateKey(elliptic.P256(),crypto_rand.Reader)
	if err!=nil {
		return err
	}

	rawPublicKey:=privateKey.PublicKey
	publicKey:=append(rawPublicKey.X.Bytes(),rawPublicKey.Y.Bytes()...)
	publicKeyHex:=hex.EncodeToString(publicKey)

	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"register",[][]byte{[]byte(req.UserId),[]byte(req.Mobile),[]byte(publicKeyHex)})
	if err!=nil {
		return err
	}

	privateKeyHex,err:=utils.EncodePrivateKey(privateKey)
	if err!=nil {
		return err
	}

	rsp.UserId=req.UserId
	rsp.Mobile=req.Mobile
	rsp.PublicKey=publicKeyHex
	rsp.PrivateKey=privateKeyHex

	return nil
}

func (e *Example) GetInfo(ctx context.Context, req *example.GetInfoRequest, rsp *example.GetInfoResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	data,err:=ccs.ChaincodeQuery(utils.ChaincodeId,"getUserInfo",[][]byte{[]byte(req.Mobile)})
	if err!=nil {
		return err
	}
	rsp.Data=data

	return nil
}

func (e *Example) Avatar(ctx context.Context, req *example.AvatarRequest, rsp *example.AvatarResponse) error {
	if len(req.Data)!=int(req.FileSize) {
		return errors.New("file transfer error")
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"avatar",[][]byte{[]byte(req.Mobile),[]byte(fileId)})
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) Rename(ctx context.Context, req *example.RenameRequest, rsp *example.RenameResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"rename",[][]byte{[]byte(req.Mobile),[]byte(req.NewName)})
	if err!=nil {
		return err
	}

	return nil
}

func (e *Example) Auth(ctx context.Context, req *example.AuthRequest, rsp *example.AuthResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err!=nil {
		return err
	}

	_,err=ccs.ChaincodeUpdate(utils.ChaincodeId,"auth",[][]byte{[]byte(req.Mobile),[]byte(req.RealName),[]byte(req.IdCard)})
	if err!=nil {
		return err
	}

	return nil
}
