package handler

import (
	"context"

		example "190316/user/proto/example"
	"190316/models"
	"190316/utils"
	"github.com/afocus/captcha"
	"image/color"
	"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/SubmailDem/submail"
	"math/rand"
	"strconv"
	"crypto/ecdsa"
	"crypto/elliptic"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"path"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GenerateId(ctx context.Context, req *example.GenerateIdRequest, rsp *example.GenerateIdResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "generateUserId", [][]byte{})
	if err != nil {
		return err
	}

	rsp.UserId = string(data)

	return nil
}

func (e *Example) Captcha(ctx context.Context, req *example.CaptchaRequest, rsp *example.CaptchaResponse) error {
	cpt := captcha.New()
	err := cpt.SetFont("comic.ttf")
	if err != nil {
		return err
	}
	cpt.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	cpt.SetFrontColor(color.RGBA{255,255,255,255})
	cpt.SetDisturbance(captcha.MEDIUM)
	cpt.SetSize(90, 40)

	img, str := cpt.Create(4, captcha.NUM)

	cnn, err := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if err != nil {
		return err
	}
	defer cnn.Close()

	key := "captcha_" + req.UserId
	_, err = cnn.Do("set", key, str, "EX", 3600)
	if err != nil {
		return err
	}

	rsp.Pix = img.Pix
	rsp.Stride = int64(img.Stride)
	rsp.Min = &example.CaptchaResponse_Point{X: int64(img.Rect.Min.X), Y: int64(img.Rect.Min.Y)}
	rsp.Max = &example.CaptchaResponse_Point{X: int64(img.Rect.Max.X), Y: int64(img.Rect.Max.Y)}

	return nil
}

func (e *Example) SmsCaptcha(ctx context.Context, req *example.SmsCaptchaRequest, rsp *example.SmsCaptchaResponse) error {
	cnn, err := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if err != nil {
		return err
	}
	defer cnn.Close()

	aKey := "captcha_" + req.UserId
	str, err := redis.String(cnn.Do("get", aKey))
	if err != nil {
		return err
	}

	if str == "" || str != req.Captcha {
		return errors.New("invalid captcha")
	}

	code := rand.Intn(8999) + 1001

	sKey := "sms_" + req.Mobile
	_, err = cnn.Do("set", sKey, strconv.Itoa(code), "EX", 3600)
	if err != nil {
		return err
	}

	messageconfig := make(map[string]string)
	messageconfig["appid"] = utils.MessageAppId
	messageconfig["appkey"] = utils.MessageAppKey
	messageconfig["signtype"] = "md5"

	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	submail.MessageXSendSetProject(messagexsend, utils.MessageProject)
	submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(code))
	submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)

	return nil
}

func (e *Example) Register(ctx context.Context, req *example.RegisterRequest, rsp *example.RegisterResponse) error {
	cnn, err := redis.Dial("tcp", utils.RedisHost + ":" + utils.RedisPort, redis.DialPassword(utils.RedisPassword))
	if err != nil {
		return err
	}
	defer cnn.Close()

	sKey := "sms_" + req.Mobile
	str, err := redis.String(cnn.Do("get", sKey))
	if err != nil {
		return err
	}

	if str == "" || str != req.SmsCaptcha {
		return errors.New("invalid sms captcha")
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), crypto_rand.Reader)
	if err != nil {
		return err
	}

	rawPublicKey := privateKey.PublicKey
	publicKey := append(rawPublicKey.X.Bytes(), rawPublicKey.Y.Bytes()...)
	publicKeyHex := hex.EncodeToString(publicKey)

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "register", [][]byte{[]byte(req.UserId), []byte(req.Mobile), []byte(publicKeyHex)})
	if err != nil {
		return err
	}

	privateKeyHex, err := utils.EncodePrivateKey(privateKey)
	if err != nil {
		return err
	}

	rsp.UserId = req.UserId
	rsp.Mobile = req.Mobile
	rsp.PublicKey = publicKeyHex
	rsp.PrivateKey = privateKeyHex

	return nil
}

func (e *Example) GetInfo(ctx context.Context, req *example.GetInfoRequest, rsp *example.GetInfoResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getUserInfo", [][]byte{[]byte(req.Mobile)})
	if err != nil {
		return err
	}

	rsp.Data = data

	return nil
}

func (e *Example) Avatar(ctx context.Context, req *example.AvatarRequest, rsp *example.AvatarResponse) error {
	if len(req.Data) != int(req.FileSize) {
		return errors.New("file transfer error")
	}

	ext := path.Ext(req.FileName)
	fileId, err := utils.UploadFile(req.Data, ext[1:])
	if err != nil {
		return err
	}

	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "avatar", [][]byte{[]byte(req.Mobile), []byte(fileId)})
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) Rename(ctx context.Context, req *example.RenameRequest, rsp *example.RenameResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "rename", [][]byte{[]byte(req.Mobile), []byte(req.NewName)})
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) Auth(ctx context.Context, req *example.AuthRequest, rsp *example.AuthResponse) error {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return err
	}

	_, err = ccs.ChaincodeUpdate(utils.ChaincodeId, "auth", [][]byte{[]byte(req.Mobile), []byte(req.RealName), []byte(req.IdCard)})
	if err != nil {
		return err
	}

	return nil
}
