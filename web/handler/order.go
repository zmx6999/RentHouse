package handler

import (
	"net/http"
	"encoding/json"
		ORDER "190120/Order/proto/example"
	"context"
	"github.com/micro/go-grpc"
	"github.com/zmx6999/FormValidation/FormValidation"
	"github.com/julienschmidt/httprouter"
	"190120/utils"
	"strconv"
)

func AddOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
			FieldName:       "house_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "house id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "start_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "start date cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:       "end_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "end date id cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
		&FormValidation.FieldValidation{
			FieldName:"start_date",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"start date invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"end date invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end_date",
			ValidMethodName:"EndDate",
			ValidMethodArgs:[]interface{}{request["start_date"]},
			ErrMsg:"end date should be later than start date",
			Trim:true,
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
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.Add(context.TODO(), &ORDER.AddRequest{
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

func GetOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	role:=utils.GetParam("role",r)

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

	var orders []map[string]interface{}
	json.Unmarshal(rsp.Data,&orders)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"orders":orders,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func HandleOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	orderId,err:=strconv.Atoi(ps.ByName("order_id"))
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

	action:=utils.GetParam("action",r)

	service:=grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.Handle(context.TODO(), &ORDER.HandleRequest{
		SessionId:ck.Value,
		Action:action,
		OrderId:int64(orderId),
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

func Comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	orderId,err:=strconv.Atoi(ps.ByName("order_id"))
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

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fvs:=[]*FormValidation.FieldValidation{
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
	exampleClient := ORDER.NewExampleService("go.micro.srv.Order", service.Client())
	rsp, err := exampleClient.Comment(context.TODO(), &ORDER.CommentRequest{
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
		"msg": rsp.ErrMsg,
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
