package handler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"encoding/json"
	AREA "190120/Area/proto/example"
	"context"
)

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

	var areaList []map[string]interface{}
	json.Unmarshal(rsp.Data,&areaList)
	// we want to augment the response
	response := map[string]interface{}{
		"code": rsp.ErrCode,
		"msg": rsp.ErrMsg,
		"data": map[string]interface{}{
			"areas": areaList,
		},
	}

	// encode and write the response as json
	w.Header().Set("Content-type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
