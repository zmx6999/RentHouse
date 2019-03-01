package handler

import (
	"context"

		example "190222/user/proto/example"
	"190222/utils"
	"190222/models"
	"github.com/afocus/captcha"
	"image/color"
	"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/SubmailDem/submail"
	"strconv"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/rand"
	crypto_rand "crypto/rand"
		"encoding/hex"
	"path"
	)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GenerateId(ctx context.Context, req *example.GenerateIdRequest, rsp *example.GenerateIdResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	data,err:=ccs.ChaincodeQuery("generateUserId",[][]byte{})
	if err!=nil {
		return err
	}

	rsp.UserId=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Captcha(ctx context.Context, req *example.CaptchaRequest, rsp *example.CaptchaResponse) error {
	capt:=captcha.New()
	err:=capt.SetFont("comic.ttf")
	if err!=nil {
		return err
	}

	capt.SetBkgColor(color.RGBA{255, 0, 0, 255},color.RGBA{0, 0, 255, 255},color.RGBA{0, 153, 0, 255})
	capt.SetFrontColor(color.RGBA{255,255,255,255})
	capt.SetSize(91,41)
	capt.SetDisturbance(captcha.MEDIUM)

	img,str:=capt.Create(4,captcha.NUM)

	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	_,err=conn.Do("set","captcha_"+req.UserId,str,"EX",3600)
	if err!=nil {
		return err
	}

	rgba:=img.RGBA
	rsp.Max=&example.CaptchaResponse_Point{X:int64(rgba.Rect.Max.X),Y:int64(rgba.Rect.Max.Y)}
	rsp.Min=&example.CaptchaResponse_Point{X:int64(rgba.Rect.Min.X),Y:int64(rgba.Rect.Min.Y)}
	rsp.Pix=rgba.Pix
	rsp.Stride=int64(rgba.Stride)

	return nil
}

func (e *Example) SmsCaptcha(ctx context.Context, req *example.SmsCaptchaRequest, rsp *example.SmsCaptchaResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	str,err:=redis.String(conn.Do("get","captcha_"+req.UserId))
	if err!=nil {
		return err
	}

	if str=="" || req.Captcha!=str {
		return errors.New("invalid captcha")
	}

	code:=rand.Intn(8999)+1001
	codeStr:=strconv.Itoa(code)
	_,err=conn.Do("set","mobile_"+req.Mobile,codeStr,"EX",3600)
	if err!=nil {
		return err
	}

	messageconfig := make(map[string]string)
	messageconfig["appid"] = utils.MessageAppId
	messageconfig["appkey"] = utils.MessageAppKey
	messageconfig["signtype"] = "md5"

	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, utils.MessageProject)
	submail.MessageXSendAddVar(messagexsend, "code", codeStr)
	submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Register(ctx context.Context, req *example.RegisterRequest, rsp *example.RegisterResponse) error {
	conn,err:=redis.Dial("tcp",utils.RedisHost+":"+utils.RedisPort)
	if err!=nil {
		return err
	}
	defer conn.Close()

	str,err:=redis.String(conn.Do("get","mobile_"+req.Mobile))
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

	privateKeyData,err:=models.EncodePrivateKey(privateKey)
	if err!=nil {
		return err
	}

	privateKeyHex:=hex.EncodeToString(privateKeyData)

	rawPublicKey:=privateKey.PublicKey
	publicKey:=append(rawPublicKey.X.Bytes(),rawPublicKey.Y.Bytes()...)
	publicKeyHex:=hex.EncodeToString(publicKey)

	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	_,err=ccs.ChaincodeUpdate("register",[][]byte{[]byte(req.UserId),[]byte(req.Mobile),[]byte(publicKeyHex)})
	if err!=nil {
		return err
	}

	rsp.PrivateKey=privateKeyHex
	rsp.PublicKey=publicKeyHex
	rsp.UserId=req.UserId
	rsp.Mobile=req.Mobile
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) GetInfo(ctx context.Context, req *example.GetInfoRequest, rsp *example.GetInfoResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	data,err:=ccs.ChaincodeQuery("getUserInfo",[][]byte{[]byte(req.Mobile)})
	if err!=nil {
		return err
	}

	rsp.Data=data
	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Avatar(ctx context.Context, req *example.AvatarRequest, rsp *example.AvatarResponse) error {
	if int64(len(req.Data))!=req.FileSize {
		return errors.New("file transfer error")
	}

	ext:=path.Ext(req.FileName)
	fileId,err:=utils.UploadFile(req.Data,ext[1:])
	if err!=nil {
		return err
	}

	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	_,err=ccs.ChaincodeUpdate("setUserAvatar",[][]byte{[]byte(req.Mobile),[]byte(fileId)})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Rename(ctx context.Context, req *example.RenameRequest, rsp *example.RenameResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	_,err=ccs.ChaincodeUpdate("setUserName",[][]byte{[]byte(req.Mobile),[]byte(req.NewName)})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}

func (e *Example) Auth(ctx context.Context, req *example.AuthRequest, rsp *example.AuthResponse) error {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return err
	}
	defer ccs.Close()

	_,err=ccs.ChaincodeUpdate("auth",[][]byte{[]byte(req.Mobile),[]byte(req.RealName),[]byte(req.IdCard)})
	if err!=nil {
		return err
	}

	rsp.Code=utils.RECODE_OK
	rsp.Msg=utils.RecodeText(rsp.Code)

	return nil
}
