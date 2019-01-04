package handler

import (
	"context"
	"encoding/json"
	"net/http"
	GETAREA "sss/GetArea/proto/example"
	GETIMAGECD "sss/GetImageCd/proto/example"
	GETSMSCD "sss/GetSmsCd/proto/example"
	REGISTER "sss/Register/proto/example"
	GETSESSION "sss/GetSession/proto/example"
	LOGIN "sss/Login/proto/example"
	LOGOUT "sss/Logout/proto/example"
	GETUSERINFO "sss/GetUserInfo/proto/example"
	POSTAVATAR "sss/PostAvatar/proto/example"
	PUTUSERINFO "sss/PutUserInfo/proto/example"
	USERAUTH "sss/UserAuth/proto/example"
	POSTHOUSES "sss/PostHouses/proto/example"
	GETUSERHOUSES "sss/GetUserHouses/proto/example"
	POSTHOUSEIMAGE "sss/PostHouseImage/proto/example"
	GETHOUSEINFO "sss/GetHouseInfo/proto/example"
	GETINDEXBANNER "sss/GetIndexBanner/proto/example"
	GETHOUSES "sss/GetHouses/proto/example"
	POSTORDER "sss/PostOrder/proto/example"
	GETUSERORDER "sss/GetUserOrder/proto/example"
	PUTORDER "sss/PutOrder/proto/example"
	PUTCOMMENT "sss/PutComment/proto/example"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/png"
	"github.com/afocus/captcha"
	"github.com/micro/go-grpc"
	"regexp"
	"sss/181231/utils"
	"io/ioutil"
	"sss/181231/models"
	"strconv"
)

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", service.Client())
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{

	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	areaList:=[]map[string]interface{}{}
	for _,v:=range rsp.Data{
		area:=map[string]interface{}{
			"area_id":v.Id,
			"area_name":v.Name,
		}
		areaList=append(areaList,area)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":areaList,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", service.Client())
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid:ps.ByName("uuid"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var rgba image.RGBA
	rgba.Pix=[]uint8(rsp.Pix)
	rgba.Stride=int(rsp.Stride)
	rgba.Rect.Max.X=int(rsp.Max.X)
	rgba.Rect.Max.Y=int(rsp.Max.Y)
	rgba.Rect.Min.X=int(rsp.Min.X)
	rgba.Rect.Min.Y=int(rsp.Min.Y)

	var img captcha.Image
	img.RGBA=&rgba
	png.Encode(w,img)
}

func GetSmsCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	mobile:=ps.ByName("mobile")
	re:=regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	if !re.MatchString(mobile) {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Invalid mobile number",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	q:=r.URL.Query()
	uuid:=q["uuid"][0]
	text:=q["text"][0]

	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmsCd", service.Client())
	rsp, err := exampleClient.GetSmsCd(context.TODO(), &GETSMSCD.Request{
		Mobile:mobile,
		Uuid:uuid,
		Text:text,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
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

	if request["mobile"]=="" || request["password"]=="" || request["sms_code"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please compete filling",
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
	exampleClient := REGISTER.NewExampleService("go.micro.srv.Register", service.Client())
	rsp, err := exampleClient.Register(context.TODO(), &REGISTER.Request{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
		SmsCode:request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		http.SetCookie(w,&http.Cookie{Name:"login",Path:"/",Value:rsp.SessionId,MaxAge:600})
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", service.Client())
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=map[string]interface{}{}
	data["name"]=rsp.Data
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
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

	if request["mobile"]=="" || request["password"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please compete filling",
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
	exampleClient := LOGIN.NewExampleService("go.micro.srv.Login", service.Client())
	rsp, err := exampleClient.Login(context.TODO(), &LOGIN.Request{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	ck,err:=r.Cookie("login")
	if err!=nil || ck.Value=="" {
		http.SetCookie(w,&http.Cookie{Name:"login",Value:rsp.SessionId,Path:"/",MaxAge:600})
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}


func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := LOGOUT.NewExampleService("go.micro.srv.Logout", service.Client())

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := exampleClient.Logout(context.TODO(), &LOGOUT.Request{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	ck,err=r.Cookie("login")
	if err==nil && ck.Value!="" {
		http.SetCookie(w,&http.Cookie{Name:"login",Path:"/",MaxAge:-1})
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", service.Client())

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=map[string]interface{}{}
	data["user_id"]=rsp.UserId
	data["name"]=rsp.Name
	data["mobile"]=rsp.Mobile
	data["real_name"]=rsp.RealName
	data["id_card"]=rsp.IdCard
	data["avatar_url"]=utils.AddDomain2Url(rsp.AvatarUrl)
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	avatar,head,err:=utils.PrepareUpload(r,"avatar",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_IOERR,
			"message":utils.RecodeTest(utils.RECODE_IOERR)+":"+err.Error(),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// call the backend service
	exampleClient := POSTAVATAR.NewExampleService("go.micro.srv.PostAvatar", service.Client())

	rsp, err := exampleClient.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		SessionId:ck.Value,
		Avatar:avatar,
		Filename:head.Filename,
		FileSize:head.Size,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=map[string]interface{}{}
	data["avatar_url"]=utils.AddDomain2Url(rsp.AvatarUrl)
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PutUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["new_name"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please input a new name",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := PUTUSERINFO.NewExampleService("go.micro.srv.PutUserInfo", service.Client())
	rsp, err := exampleClient.PutUserInfo(context.TODO(), &PUTUSERINFO.Request{
		SessionId:ck.Value,
		Name:request["new_name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=map[string]interface{}{}
	data["new_name"]=rsp.Name
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["real_name"]=="" || request["id_card"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please compete filling",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	idCard:=request["id_card"].(string)
	re:=regexp.MustCompile(`^[1-9][0-9]{17}$`)
	if !re.MatchString(idCard) {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Invalid ID card",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := USERAUTH.NewExampleService("go.micro.srv.UserAuth", service.Client())
	rsp, err := exampleClient.PostUserAuth(context.TODO(), &USERAUTH.Request{
		SessionId:ck.Value,
		RealName:request["real_name"].(string),
		IdCard:idCard,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostHouse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	data,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := POSTHOUSES.NewExampleService("go.micro.srv.PostHouses", service.Client())
	rsp, err := exampleClient.PostHouse(context.TODO(), &POSTHOUSES.Request{
		SessionId:ck.Value,
		Body:data,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":map[string]interface{}{
			"house_id":rsp.HouseId,
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
	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := GETUSERHOUSES.NewExampleService("go.micro.srv.GetUserHouses", service.Client())
	rsp, err := exampleClient.GetUserHouses(context.TODO(), &GETUSERHOUSES.Request{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	var houses []*models.House
	json.Unmarshal(rsp.Data,&houses)
	data:=[]map[string]interface{}{}
	for _,house:=range houses{
		data=append(data,house.Info())
	}
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostHouseImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	avatar,head,err:=utils.PrepareUpload(r,"image",[]string{"jpg","png","jpeg"},1024*1024*2)
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_IOERR,
			"message":utils.RecodeTest(utils.RECODE_IOERR)+":"+err.Error(),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// call the backend service
	exampleClient := POSTHOUSEIMAGE.NewExampleService("go.micro.srv.PostHouseImage", service.Client())
	houseId,_:=strconv.Atoi(ps.ByName("id"))
	rsp, err := exampleClient.PostHouseImage(context.TODO(), &POSTHOUSEIMAGE.Request{
		SessionId:ck.Value,
		Image:avatar,
		FileName:head.Filename,
		FileSize:head.Size,
		HouseId:int64(houseId),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=map[string]interface{}{}
	data["url"]=utils.AddDomain2Url(rsp.Url)
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouseInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := GETHOUSEINFO.NewExampleService("go.micro.srv.GetHouseInfo", service.Client())
	houseId,_:=strconv.Atoi(ps.ByName("id"))
	rsp, err := exampleClient.GetHouseInfo(context.TODO(), &GETHOUSEINFO.Request{
		SessionId:ck.Value,
		HouseId:int64(houseId),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	var house models.House
	json.Unmarshal(rsp.Data,&house)
	data:=map[string]interface{}{}
	data["user_id"]=rsp.UserId
	data["house"]=house.Desc()
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetIndexBanner(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := GETINDEXBANNER.NewExampleService("go.micro.srv.GetIndexBanner", service.Client())
	rsp, err := exampleClient.GetIndexBanner(context.TODO(), &GETINDEXBANNER.Request{

	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&data)
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q:=r.URL.Query()
	areaIds:=q["area_id"]
	areaId:=""
	if len(areaIds)>0 {
		areaId=areaIds[0]
	}
	pages:=q["page"]
	page:=""
	if len(pages)>0 {
		page=pages[0]
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := GETHOUSES.NewExampleService("go.micro.srv.GetHouses", service.Client())
	rsp, err := exampleClient.GetHouses(context.TODO(), &GETHOUSES.Request{
		AreaId:areaId,
		Page:page,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	houses:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&houses)
	data:=map[string]interface{}{}
	data["houses"]=houses
	data["total_page"]=rsp.TotalPages
	data["page"]=rsp.Page
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	data,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

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
	exampleClient := POSTORDER.NewExampleService("go.micro.srv.PostOrder", service.Client())
	rsp, err := exampleClient.PostOrder(context.TODO(), &POSTORDER.Request{
		SessionId:ck.Value,
		Data:data,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":map[string]interface{}{
			"order_id":rsp.OrderId,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	q:=r.URL.Query()
	roles:=q["role"]
	role:=""
	if len(roles)>0 {
		role=roles[0]
	}

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := GETUSERORDER.NewExampleService("go.micro.srv.GetUserOrder", service.Client())
	rsp, err := exampleClient.GetUserOrder(context.TODO(), &GETUSERORDER.Request{
		SessionId:ck.Value,
		Role:role,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=[]map[string]interface{}{}
	json.Unmarshal(rsp.Data,&data)
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PutOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["action"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please input action",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	orderId:=ps.ByName("id")

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := PUTORDER.NewExampleService("go.micro.srv.PutOrder", service.Client())
	rsp, err := exampleClient.PutOrder(context.TODO(), &PUTORDER.Request{
		SessionId:ck.Value,
		OrderId:orderId,
		Action:request["action"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PutComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["comment"]=="" {
		response := map[string]interface{}{
			"code":utils.RECODE_NODATA,
			"message":"Please comment",
		}

		// encode and write the response as json
		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	ck,err:=r.Cookie("login")
	if err!=nil {
		response := map[string]interface{}{
			"code":utils.RECODE_SESSIONERR,
			"message":utils.RecodeTest(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	orderId:=ps.ByName("id")

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := PUTCOMMENT.NewExampleService("go.micro.srv.PutComment", service.Client())
	rsp, err := exampleClient.PutComment(context.TODO(), &PUTCOMMENT.Request{
		SessionId:ck.Value,
		OrderId:orderId,
		Comment:request["comment"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"code":rsp.ErrCode,
		"message":rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
