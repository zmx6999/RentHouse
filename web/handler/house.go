package handler

import (
	"net/http"
	"encoding/json"
	HOUSE "190120/House/proto/example"
	"context"
	"github.com/micro/go-grpc"
	"190120/utils"
	"github.com/zmx6999/FormValidation/FormValidation"
	"github.com/julienschmidt/httprouter"
	"strconv"
)

func AddHouse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
			ErrMsg:          "area id cannot be empty",
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
	data,_:=json.Marshal(request)

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.Add(context.TODO(), &HOUSE.AddRequest{
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
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.GetHouses(context.TODO(), &HOUSE.GetHousesRequest{
		SessionId:ck.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var houses []map[string]interface{}
	json.Unmarshal(rsp.Data,&houses)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"houses":houses,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func UploadHouseImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	houseId,err:=strconv.Atoi(ps.ByName("house_id"))
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": utils.RecodeText(utils.RECODE_PARAMERR),
		}

		// encode and write the response as json
		w.Header().Set("Content-type","application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	data,head,err:=utils.PrepareUpload(r,"image",[]string{"jpg","jpeg","png"},1024*1024*2)
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
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.UploadImage(context.TODO(), &HOUSE.UploadImageRequest{
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
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouseDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	houseId,err:=strconv.Atoi(ps.ByName("house_id"))
	if err!=nil {
		response := map[string]interface{}{
			"code": utils.RECODE_PARAMERR,
			"msg": utils.RecodeText(utils.RECODE_PARAMERR),
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
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.GetHouseDetail(context.TODO(), &HOUSE.GetHouseDetailRequest{
		HouseId:int64(houseId),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var house map[string]interface{}
	json.Unmarshal(rsp.Data,&house)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"house":house,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
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

	var houses []map[string]interface{}
	json.Unmarshal(rsp.Data,&houses)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"houses":houses,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request:=map[string]interface{}{}
	request["area_id"]=utils.GetParam("area_id",r)
	request["start_date"]=utils.GetParam("start_date",r)
	request["end_date"]=utils.GetParam("end_date",r)
	request["page"]=utils.GetParam("page",r)
	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "start_date",
			ValidMethodName: "Date",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "Invalid start date",
			Trim:            true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "end_date",
			ValidMethodName: "Date",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "Invalid end date",
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
	exampleClient := HOUSE.NewExampleService("go.micro.srv.House", service.Client())
	rsp, err := exampleClient.Search(context.TODO(), &HOUSE.SearchRequest{
		AreaId:request["area_id"].(string),
		StartDate:request["start_date"].(string),
		EndDate:request["end_date"].(string),
		Page:request["page"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var houses []map[string]interface{}
	json.Unmarshal(rsp.Data,&houses)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"houses":houses,
			"total_page":rsp.TotalPage,
			"page":rsp.Page,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
