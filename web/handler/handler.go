package handler

import (
	"net/http"
	"encoding/json"

	AREA "190222/area/proto/example"
	USER "190222/user/proto/example"
	HOUSE "190222/house/proto/example"
	ORDER "190222/order/proto/example"
		"github.com/micro/go-grpc"
	"context"
	"github.com/julienschmidt/httprouter"
	"image/png"
	"github.com/afocus/captcha"
	"image"
	"190222/utils"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190222/models"
	)

/*
b51c6f0467513aa7a3640744adb5e43c9efef9eb
2eff810301010a507269766174654b657901ff8200010201095075626c69634b657901ff840001014401ff860000002fff83030101095075626c69634b657901ff840001030105437572766501100001015801ff860001015901ff860000000aff85050102ff8800000046ff8201011963727970746f2f656c6c69707469632e703235364375727665ff890301010970323536437572766501ff8a000101010b4375727665506172616d7301ff8c00000053ff8b0301010b4375727665506172616d7301ff8c00010701015001ff860001014e01ff860001014201ff86000102477801ff86000102477901ff8600010742697453697a6501040001044e616d65010c000000fe012cff8affbd01012102ffffffff00000001000000000000000000000000ffffffffffffffffffffffff012102ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325510121025ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b0121026b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2960121024fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f501fe02000105502d3235360000012102157082d835e4a65389e7389290a707b145c4c54a9a0b14135a8ad509fdea34b5012102507ab8dec06507a68b0cb64eda6e181ed559eca4a56caf18596b3f7c5ef71fe90001210282e778b4e6ab1ebefa63266cb68290e6ae3db87514c82df5b7524674fda2d7e200
157082d835e4a65389e7389290a707b145c4c54a9a0b14135a8ad509fdea34b5507ab8dec06507a68b0cb64eda6e181ed559eca4a56caf18596b3f7c5ef71fe9
 */

/*
28c0a02a9ef3931c5b685f4f30f383312ee35b9d
2eff810301010a507269766174654b657901ff8200010201095075626c69634b657901ff840001014401ff860000002fff83030101095075626c69634b657901ff840001030105437572766501100001015801ff860001015901ff860000000aff85050102ff8800000046ff8201011963727970746f2f656c6c69707469632e703235364375727665ff890301010970323536437572766501ff8a000101010b4375727665506172616d7301ff8c00000053ff8b0301010b4375727665506172616d7301ff8c00010701015001ff860001014e01ff860001014201ff86000102477801ff86000102477901ff8600010742697453697a6501040001044e616d65010c000000fe012cff8affbd01012102ffffffff00000001000000000000000000000000ffffffffffffffffffffffff012102ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325510121025ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b0121026b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2960121024fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f501fe02000105502d32353600000121026e3176c14ed4c4fb666f6b1430d305c73bf0d885f9570c3d8f73290238eb206801210281001387f8ea8908fcfe3cca8541355139f5986088f7b1f28d1b21d7fdb4e9a6000121025606453c870aba426ae3ae576e7fec3df60d4a2d03f838e2a9ec54c6c5106b0800
6e3176c14ed4c4fb666f6b1430d305c73bf0d885f9570c3d8f73290238eb206881001387f8ea8908fcfe3cca8541355139f5986088f7b1f28d1b21d7fdb4e9a6
 */

/*
func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
 */

func handleResponse(w http.ResponseWriter,code string,msg string,data interface{})  {
	w.Header().Set("Content-Type","application/json")
	// we want to augment the response
	response := map[string]interface{}{
		"code": code,
		"msg": msg,
		"data": data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func verifyUser(privateKeyHex string,userMobile string) (bool,error) {
	ccs,err:=models.Initialize(utils.ChannelId,utils.User,utils.ChaincodeId,utils.FabricSDKConfig)
	if err!=nil {
		return false,err
	}
	defer ccs.Close()

	return ccs.VerifyUser(privateKeyHex,userMobile)
}

func GetArea(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := AREA.NewExampleService("go.micro.srv.area", service.Client())
	rsp, err := exampleClient.GetArea(context.TODO(), &AREA.GetAreaRequest{

	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data []map[string]interface{}
	err = json.Unmarshal(rsp.Data, &data)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func GenerateId(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.GenerateId(context.TODO(), &USER.GenerateIdRequest{

	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	data:=map[string]interface{}{
		"user_id":string(rsp.UserId),
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func Captcha(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	userId:=p.ByName("user_id")

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Captcha(context.TODO(), &USER.CaptchaRequest{
		UserId:userId,
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	img:=captcha.Image{&image.RGBA{Rect:image.Rect(int(rsp.Min.X),int(rsp.Min.Y),int(rsp.Max.X),int(rsp.Max.Y)),Pix:rsp.Pix,Stride:int(rsp.Stride)}}

	// encode and write the response as json
	if err := png.Encode(w,img); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
}

func SmsCaptcha(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	request:=make(map[string]interface{})
	request["user_id"]=utils.GetParam("user_id",r)
	request["mobile"]=utils.GetParam("mobile",r)
	request["captcha"]=utils.GetParam("captcha",r)

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "user_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "user_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "captcha cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.SmsCaptcha(context.TODO(), &USER.SmsCaptchaRequest{
		UserId:request["user_id"].(string),
		Mobile:request["mobile"].(string),
		Captcha:request["captcha"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func Register(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "user_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "user_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "sms_captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "sms_captcha cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Register(context.TODO(), &USER.RegisterRequest{
		UserId:request["user_id"].(string),
		Mobile:request["mobile"].(string),
		SmsCaptcha:request["sms_captcha"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	data:= map[string]interface{}{
		"user_id": rsp.UserId,
		"mobile": rsp.Mobile,
		"public_key": rsp.PublicKey,
		"private_key": rsp.PrivateKey,
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.GetInfo(context.TODO(), &USER.GetInfoRequest{
		Mobile:request["mobile"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data map[string]interface{}
	err=json.Unmarshal(rsp.Data,&data)
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func Avatar(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	request:=make(map[string]interface{})
	request["mobile"]=r.FormValue("mobile")
	request["private_key"]=r.FormValue("private_key")

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	data,head,err:=utils.PrepareUploadFile(r,"avatar",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Avatar(context.TODO(), &USER.AvatarRequest{
		Mobile:request["mobile"].(string),
		Data:data,
		FileName:head.Filename,
		FileSize:head.Size,
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func Rename(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "new_name",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "new_name cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Rename(context.TODO(), &USER.RenameRequest{
		Mobile:request["mobile"].(string),
		NewName:request["new_name"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func Auth(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "real_name",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "real_name cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "id_card",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "id_card cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"id_card",
			ValidMethodName:"ChineseIdCard",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"id_card invalid",
			Trim:true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Auth(context.TODO(), &USER.AuthRequest{
		Mobile:request["mobile"].(string),
		RealName:request["real_name"].(string),
		IdCard:request["id_card"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func AddHouse(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "title",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "title cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "area_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "area_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	request["private_key"]=nil
	data,err:=json.Marshal(request)
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.Add(context.TODO(), &HOUSE.AddRequest{
		Mobile:request["mobile"].(string),
		Data:data,
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func GetHouseList(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.GetHouseList(context.TODO(), &HOUSE.GetHouseListRequest{
		Mobile:request["mobile"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data []map[string]interface{}
	err=json.Unmarshal(rsp.Data,&data)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func GetHouseDesc(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.GetHouseDesc(context.TODO(), &HOUSE.GetHouseDescRequest{
		HouseId:p.ByName("house_id"),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data map[string]interface{}
	err=json.Unmarshal(rsp.Data,&data)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func UploadHouseImage(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	// decode the incoming request as json
	request:=make(map[string]interface{})
	request["mobile"]=r.FormValue("mobile")
	request["private_key"]=r.FormValue("private_key")

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	data,head,err:=utils.PrepareUploadFile(r,"image",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.UploadImage(context.TODO(), &HOUSE.UploadImageRequest{
		HouseId:p.ByName("house_id"),
		Mobile:request["mobile"].(string),
		Data:data,
		FileName:head.Filename,
		FileSize:head.Size,
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func GetIndexHouseList(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.GetIndexList(context.TODO(), &HOUSE.GetIndexListRequest{

	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data []map[string]interface{}
	err=json.Unmarshal(rsp.Data,&data)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func AddOrder(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "house_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "house_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "start_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "start_date cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"start_date",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"start_date invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "end_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "end_date cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"end_date invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"EndDate",
			ValidMethodArgs:[]interface{}{request["start_date"]},
			ErrMsg:"end_date should be later than start_date",
			Trim:true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	request["private_key"]=nil
	data,err:=json.Marshal(request)
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	rsp, err := exampleClient.Add(context.TODO(), &ORDER.AddRequest{
		Mobile:request["mobile"].(string),
		Data:data,
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func GetOrderList(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	rsp, err := exampleClient.GetList(context.TODO(), &ORDER.GetListRequest{
		Mobile:request["mobile"].(string),
		Role:utils.GetParam("role",r),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	var data []map[string]interface{}
	err=json.Unmarshal(rsp.Data,&data)
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,data)
}

func HandleOrder(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "order_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "order_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	rsp, err := exampleClient.Handle(context.TODO(), &ORDER.HandleRequest{
		Mobile:request["mobile"].(string),
		Action:utils.GetParam("action",r),
		OrderId:request["order_id"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}

func Comment(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "mobile cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "order_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "order_id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "comment",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "comment cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	valid,err:=verifyUser(request["private_key"].(string),request["mobile"].(string))
	if err!=nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}
	if !valid {
		handleResponse(w,"403","invalid private key",nil)
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	rsp, err := exampleClient.Comment(context.TODO(), &ORDER.CommentRequest{
		Mobile:request["mobile"].(string),
		OrderId:request["order_id"].(string),
		Comment:request["comment"].(string),
	})
	if err != nil {
		handleResponse(w,"500",err.Error(),nil)
		return
	}

	handleResponse(w,rsp.Code,rsp.Msg,nil)
}
