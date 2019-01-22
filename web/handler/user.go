package handler

import (
	"net/http"
			USER "190120/User/proto/example"
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"image/png"
		"github.com/afocus/captcha"
	"image"
	"encoding/json"
			"190120/utils"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190120/models"
)

func GetCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.GetCaptcha(context.TODO(), &USER.GetCaptchaRequest{
		Uuid:ps.ByName("uuid"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rgba:=image.RGBA{Pix:rsp.Pix,Stride:int(rsp.Stride),Rect:image.Rect(int(rsp.Min.X),int(rsp.Min.Y),int(rsp.Max.X),int(rsp.Max.Y))}
	img:=captcha.Image{RGBA:&rgba}

	// encode and write the response as json
	if err := png.Encode(w,img); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSmsCaptcha(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request:=map[string]interface{}{
		"mobile":utils.GetParam("mobile",r),
		"uuid":utils.GetParam("uuid",r),
		"captcha":utils.GetParam("captcha",r),
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
			FieldName:       "uuid",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "uuid cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "captcha cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "Invalid mobile",
			Trim:            true,
		},
	}
	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.GetSmsCaptcha(context.TODO(), &USER.GetSmsCaptchaRequest{
		Uuid:request["uuid"].(string),
		Captcha:request["captcha"].(string),
		Mobile:request["mobile"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
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
			FieldName:       "password",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "password cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "sms_captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "sms captcha cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "Invalid mobile",
			Trim:            true,
		},
	}
	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Register(context.TODO(), &USER.RegisterRequest{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
		SmsCaptcha:request["sms_captcha"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	http.SetCookie(w,&http.Cookie{Name:"login",Value:rsp.SessionId,MaxAge:3600,Path:"/"})

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
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
			FieldName:       "password",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "password cannot be empty",
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "Invalid mobile",
			Trim:            true,
		},
	}
	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err:=gv.Validate()
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Login(context.TODO(), &USER.LoginRequest{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	http.SetCookie(w,&http.Cookie{Name:"login",Value:rsp.SessionId,MaxAge:3600,Path:"/"})

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Logout(context.TODO(), &USER.LogoutRequest{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	http.SetCookie(w,&http.Cookie{Name:"login",Value:"",MaxAge:-1,Path:"/"})

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Info(context.TODO(), &USER.InfoRequest{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	user:=models.User{}
	json.Unmarshal(rsp.Data,&user)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"user_id":user.Id,
				"name":user.Name,
				"mobile":user.Mobile,
				"real_name":user.Real_name,
				"id_card":user.Id_card,
				"avatar_url":utils.AddDomain2Url(user.Avatar_url),
			},
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Avatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	data,head,err:=utils.PrepareUpload(r,"avatar",[]string{"jpg","jpeg","png"},1024*1024*2)
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_DATAERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Avatar(context.TODO(), &USER.AvatarRequest{
		SessionId:ck.Value,
		Data:data,
		FileSize:head.Size,
		FileName:head.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func UpdateUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "username",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "username cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}
	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err=gv.Validate()
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.UpdateUserName(context.TODO(), &USER.UpdateUserNameRequest{
		SessionId:ck.Value,
		UserName:request["username"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "real_name",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "real name cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "id_card",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "id card cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "id_card",
			ValidMethodName: "ChineseIdCard",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "invalid id card",
			Trim:            true,
		},
	}
	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_,err=gv.Validate()
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	rsp, err := exampleClient.Auth(context.TODO(), &USER.AuthRequest{
		SessionId:ck.Value,
		RealName:request["real_name"].(string),
		IdCard:request["id_card"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
