package handler

import (
	"net/http"
	"encoding/json"
	AREA "190105/Area/proto/example"
	USER "190105/User/proto/example"
	HOUSE "190105/House/proto/example"
	ORDER "190105/Order/proto/example"
	"github.com/micro/go-grpc"
	"context"
	"github.com/julienschmidt/httprouter"
	"image/png"
	"github.com/afocus/captcha"
	"image"
	"github.com/astaxie/beego"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190105/utils"
	"io/ioutil"
	"strconv"
)

/*func ExampleCall(w http.ResponseWriter, r *http.Request) {
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
}*/

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := AREA.NewExampleService("go.micro.srv.Area", service.Client())
	rsp, err := exampleClient.GetArea(context.TODO(), &AREA.GetAreaRequest{

	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	areas:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&areas)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"areas":areas,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCpt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.User", service.Client())
	beego.Info(ps.ByName("id"))
	rsp, err := exampleClient.GetImageCpt(context.TODO(), &USER.GetImageCptRequest{
		Uuid:ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response

	rgba:=image.RGBA{Pix:rsp.Pix,Stride:int(rsp.Stride),Rect:image.Rect(int(rsp.Min.X),int(rsp.Min.Y),int(rsp.Max.X),int(rsp.Max.Y))}
	img:=captcha.Image{RGBA:&rgba}

	// encode and write the response as json
	if err := png.Encode(w,img); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSmsCpt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid:=ps.ByName("id")
	mobile:=utils.GetParam(r,"mobile")
	text:=utils.GetParam(r,"text")

	request:=map[string]interface{}{
		"mobile":mobile,
	}
	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"Invalid mobile",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err:=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
	rsp, err := exampleClient.GetSmsCpt(context.TODO(), &USER.GetSmsCptRequest{
		Uuid:uuid,
		Mobile:mobile,
		Text:text,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"mobile cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"password",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ValidEmpty:true,
			ErrMsg:"password cannot be empty",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err:=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
		Text:request["text"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp.ErrCode==utils.RECODE_OK {
		http.SetCookie(w,&http.Cookie{Name:"session_id",Value:rsp.SessionId,MaxAge:600,Path:"/"})
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"mobile cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"password",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ValidEmpty:true,
			ErrMsg:"password cannot be empty",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err:=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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

	if rsp.ErrCode==utils.RECODE_OK {
		http.SetCookie(w,&http.Cookie{Name:"session_id",Value:rsp.SessionId,MaxAge:600,Path:"/"})
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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

	if rsp.ErrCode==utils.RECODE_OK {
		http.SetCookie(w,&http.Cookie{Name:"session_id",Value:"",MaxAge:-1,Path:"/"})
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &USER.GetUserInfoRequest{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp.ErrCode==utils.RECODE_OK {
		http.SetCookie(w,&http.Cookie{Name:"session_id",Value:"",MaxAge:-1,Path:"/"})
	}

	// we want to augment the response
	user:=map[string]interface{}{}
	json.Unmarshal(rsp.Data,&user)
	user["avatar_url"]=utils.AddDomain2Url(user["avatar_url"].(string))
	data:=map[string]interface{}{
		"user":user,
	}
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Rename(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"new_name",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"new name cannot be empty",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
	rsp, err := exampleClient.Rename(context.TODO(), &USER.RenameRequest{
		SessionId:ck.Value,
		NewName:request["new_name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"real_name",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"real name cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"id_card",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"ID card cannot be empty",
		},
		&FormValidation.FieldValidation{
			FieldName:"id_card",
			ValidMethodName:"ChineseIdCard",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ErrMsg:"invalid ID card",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func UploadAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	data,head,err:=utils.PrepareUpload(r,"avatar",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_IOERR,
			"msg": utils.RecodeText(utils.RECODE_IOERR)+":"+err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
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
	rsp, err := exampleClient.UploadAvatar(context.TODO(), &USER.UploadAvatarRequest{
		SessionId:ck.Value,
		Data:data,
		FileName:head.Filename,
		FileSize:head.Size,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func AddHouse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	data,_:=ioutil.ReadAll(r.Body)

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.AddHouse(context.TODO(), &HOUSE.AddHouseRequest{
		SessionId:ck.Value,
		Data:data,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"house_id": rsp.HouseId,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.GetUserHouses(context.TODO(), &HOUSE.GetUserHousesRequest{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	houses:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&houses)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"houses":houses,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func UploadHouseImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	data,head,err:=utils.PrepareUpload(r,"image",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_IOERR,
			"msg": utils.RecodeText(utils.RECODE_IOERR)+":"+err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	houseId,_:=strconv.Atoi(ps.ByName("id"))

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.UploadHouseImage(context.TODO(), &HOUSE.UploadHouseImageRequest{
		SessionId:ck.Value,
		HouseId:int64(houseId),
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
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouseDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	houseId,_:=strconv.Atoi(ps.ByName("id"))

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.GetHouseDetail(context.TODO(), &HOUSE.GetHouseDetailRequest{
		HouseId:int64(houseId),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp.ErrCode==utils.RECODE_OK {
		http.SetCookie(w,&http.Cookie{Name:"session_id",Value:"",MaxAge:-1,Path:"/"})
	}

	// we want to augment the response
	house:=map[string]interface{}{}
	json.Unmarshal(rsp.Data,&house)
	data:=map[string]interface{}{
		"house":house,
	}
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetIndexBanner(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.GetIndexBanner(context.TODO(), &HOUSE.GetIndexBannerRequest{

	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	houses:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&houses)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"houses":houses,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func SearchHouse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	areaId:=utils.GetParam(r,"area_id")
	start:=utils.GetParam(r,"start")
	end:=utils.GetParam(r,"end")
	page:=utils.GetParam(r,"page")

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.SearchHouse(context.TODO(), &HOUSE.SearchHouseRequest{
		AreaId:areaId,
		StartDate:start,
		EndDate:end,
		Page:page,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	houses:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&houses)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"houses":houses,
			"total_pages":rsp.TotalPages,
			"current_page":rsp.CurrentPage,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func AddOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	data,_:=ioutil.ReadAll(r.Body)

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.AddOrder(context.TODO(), &ORDER.AddOrderRequest{
		SessionId:ck.Value,
		Data:data,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"order_id": rsp.OrderId,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	role:=utils.GetParam(r,"role")

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.GetOrders(context.TODO(), &ORDER.GetOrdersRequest{
		SessionId:ck.Value,
		Role:role,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	orders:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&orders)
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
		"data": map[string]interface{}{
			"orders":orders,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func HandleOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	orderId,_:=strconv.Atoi(ps.ByName("id"))
	action:=utils.GetParam(r,"action")

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.HandleOrder(context.TODO(), &ORDER.HandleOrderRequest{
		SessionId:ck.Value,
		OrderId:int64(orderId),
		Action:action,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func CommentOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ck,err:=r.Cookie("session_id")
	if err!=nil || ck.Value=="" {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_SESSIONERR,
			"msg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	orderId,_:=strconv.Atoi(ps.ByName("id"))

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"comment",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			Trim:true,
			ValidEmpty:true,
			ErrMsg:"input comment please",
		},
	}
	gv:=FormValidation.GroupValidation{request,fvs}
	_,err=gv.Validate()
	if err!=nil {
		// we want to augment the response
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": err.Error(),
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.CommentOrder(context.TODO(), &ORDER.CommentOrderRequest{
		SessionId:ck.Value,
		OrderId:int64(orderId),
		Comment:request["comment"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.Msg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
