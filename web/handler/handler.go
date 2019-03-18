package handler

import (
		"encoding/json"
	"net/http"
				"190316/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	AREA "190316/area/proto/example"
	USER "190316/user/proto/example"
	HOUSE "190316/house/proto/example"
	ORDER "190316/order/proto/example"
	"context"
	"image/png"
	"image"
	"github.com/afocus/captcha"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190316/models"
	"errors"
)

/*
30296171523
e49057915f0b9e2282d64343e9e9c79dbdd65c26d12b89c5603cb7ece31e9ab00de24de6da39d02e2fb4e6df890c8a0c5620d686443a679b4ebb24dfa447c14f
2eff810301010a507269766174654b657901ff8200010201095075626c69634b657901ff840001014401ff860000002fff83030101095075626c69634b657901ff840001030105437572766501100001015801ff860001015901ff860000000aff85050102ff8800000046ff8201011963727970746f2f656c6c69707469632e703235364375727665ff890301010970323536437572766501ff8a000101010b4375727665506172616d7301ff8c00000053ff8b0301010b4375727665506172616d7301ff8c00010701015001ff860001014e01ff860001014201ff86000102477801ff86000102477901ff8600010742697453697a6501040001044e616d65010c000000fe012cff8affbd01012102ffffffff00000001000000000000000000000000ffffffffffffffffffffffff012102ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325510121025ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b0121026b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2960121024fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f501fe02000105502d3235360000012102e49057915f0b9e2282d64343e9e9c79dbdd65c26d12b89c5603cb7ece31e9ab00121020de24de6da39d02e2fb4e6df890c8a0c5620d686443a679b4ebb24dfa447c14f000121024fc555dbcd88ef6fcea6e1314208db7879eec823cf05b131b762fb1e424d54cc00
 */

/*
39526876799
1da90039edcc014b2e3026765e2752cc0dbab958d395269dc3ac23aee4ec675be2a225be03436d15ebf65ea85d017a6451818b589545bc00c358730888766f39
2eff810301010a507269766174654b657901ff8200010201095075626c69634b657901ff840001014401ff860000002fff83030101095075626c69634b657901ff840001030105437572766501100001015801ff860001015901ff860000000aff85050102ff8800000046ff8201011963727970746f2f656c6c69707469632e703235364375727665ff890301010970323536437572766501ff8a000101010b4375727665506172616d7301ff8c00000053ff8b0301010b4375727665506172616d7301ff8c00010701015001ff860001014e01ff860001014201ff86000102477801ff86000102477901ff8600010742697453697a6501040001044e616d65010c000000fe012cff8affbd01012102ffffffff00000001000000000000000000000000ffffffffffffffffffffffff012102ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325510121025ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b0121026b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2960121024fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f501fe02000105502d32353600000121021da90039edcc014b2e3026765e2752cc0dbab958d395269dc3ac23aee4ec675b012102e2a225be03436d15ebf65ea85d017a6451818b589545bc00c358730888766f390001210280fe3e20c7b30ddc3bcdaa8c066d0c5016e87e74ddc11752f6300b8d686787ca00
 */

/*
16876799888916
30296171523095
18891629820782
4782674154489504
3767998889162982
3029617152309526
 */

func handleResponse(w http.ResponseWriter, code int, msg interface{}, data interface{})  {
	w.Header().Set("content-type", "application/json")

	// we want to augment the response
	response := map[string]interface{}{
		"code": code,
		"msg": msg,
		"data": data,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func success(w http.ResponseWriter, data interface{})  {
	handleResponse(w, 200, "OK", data)
}

func handleError(w http.ResponseWriter, code int, err error)  {
	handleResponse(w, code, err.Error(), nil)
}

func handleChaincodeError(w http.ResponseWriter, err error)  {
	errInfo := make(map[string]interface{})
	json.Unmarshal([]byte(err.Error()), &errInfo)
	code := utils.GetFloatValue(errInfo, "code", 4101)
	msg := utils.GetStringValue(errInfo, "detail", "Unknown Error")
	handleResponse(w, int(code), msg, nil)
}

func getPublicKey(mobile string) (string, error) {
	ccs, err := models.Initialize(utils.ChannelId, utils.UserName, utils.OrgName, utils.ChaincodeId, utils.ConfigFile)
	if err != nil {
		return "", err
	}

	data, err := ccs.ChaincodeQuery(utils.ChaincodeId, "getUserPublicKey", [][]byte{[]byte(mobile)})
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func validateUser(mobile string, privateKeyHex string) (bool, error) {
	str := utils.GetNonce(41)

	rByte, sByte, err := utils.Sign(str, privateKeyHex)
	if err != nil {
		return false, err
	}

	publicKeyHex, err := getPublicKey(mobile)
	if err != nil {
		return false, err
	}

	return utils.Verify(str, publicKeyHex, rByte, sByte)
}

func GetAreaList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := AREA.NewExampleService("go.micro.srv.area", service.Client())
	rsp, err := exampleClient.GetList(context.TODO(), &AREA.GetListRequest{

	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	var data []map[string]interface{}
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func GenerateUserId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.GenerateId(context.TODO(), &USER.GenerateIdRequest{

	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["user_id"] = rsp.UserId

	success(w, data)
}

func Captcha(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Captcha(context.TODO(), &USER.CaptchaRequest{
		UserId: p.ByName("user_id"),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	rgba := image.RGBA{Pix: rsp.Pix, Stride: int(rsp.Stride), Rect: image.Rect(int(rsp.Min.X), int(rsp.Min.Y), int(rsp.Max.X), int(rsp.Max.Y))}
	img := captcha.Image{&rgba}

	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func SmsCaptcha(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := make(map[string]interface{})
	request["user_id"] = utils.GetParam("user_id", r)
	request["mobile"] = utils.GetParam("mobile", r)
	request["captcha"] = utils.GetParam("captcha", r)

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "user_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "user_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "captcha cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	_, err = exampleClient.SmsCaptcha(context.TODO(), &USER.SmsCaptchaRequest{
		UserId: request["user_id"].(string),
		Mobile: request["mobile"].(string),
		Captcha: request["captcha"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "user_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "user_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "sms_captcha",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "sms_captcha cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.Register(context.TODO(), &USER.RegisterRequest{
		UserId: request["user_id"].(string),
		Mobile: request["mobile"].(string),
		SmsCaptcha: request["sms_captcha"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["user_id"] = rsp.UserId
	data["mobile"] = rsp.Mobile
	data["public_key"] = rsp.PublicKey
	data["private_key"] = rsp.PrivateKey

	success(w, data)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	rsp, err := exampleClient.GetInfo(context.TODO(), &USER.GetInfoRequest{
		Mobile: request["mobile"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	data := make(map[string]interface{})
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func Avatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := make(map[string]interface{})
	request["mobile"] = r.FormValue("mobile")
	request["private_key"] = r.FormValue("private_key")

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	data, head, err := utils.PrepageUploadFile(r, "avatar", []string{"jpg", "png", "jpeg"}, 1024 * 1024 * 2)
	if err != nil {
		handleError(w, 500, err)
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	_, err = exampleClient.Avatar(context.TODO(), &USER.AvatarRequest{
		Mobile: request["mobile"].(string),
		Data: data,
		FileName: head.Filename,
		FileSize: head.Size,
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func Rename(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "new_name",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "new_name cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	_, err = exampleClient.Rename(context.TODO(), &USER.RenameRequest{
		Mobile: request["mobile"].(string),
		NewName: request["new_name"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "real_name",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "real_name cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "id_card",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "id_card cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName:"id_card",
			ValidMethodName:"ChineseIdCard",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid id_card",
			Trim:true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := USER.NewExampleService("go.micro.srv.user", service.Client())
	_, err = exampleClient.Auth(context.TODO(), &USER.AuthRequest{
		Mobile: request["mobile"].(string),
		RealName: request["real_name"].(string),
		IdCard: request["id_card"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func AddHouse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "title",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "title cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "area_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "area_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	data, _ := json.Marshal(request)

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	_, err = exampleClient.Add(context.TODO(), &HOUSE.AddRequest{
		Data: data,
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func GetLandlordHouseList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.GetLandlordList(context.TODO(), &HOUSE.GetLandlordListRequest{
		Mobile: request["mobile"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	var data []map[string]interface{}
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func UploadHouseImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := make(map[string]interface{})
	request["mobile"] = r.FormValue("mobile")
	request["private_key"] = r.FormValue("private_key")
	request["house_id"] = r.FormValue("house_id")

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "house_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "house_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	data, head, err := utils.PrepageUploadFile(r, "image", []string{"jpg", "png", "jpeg"}, 1024 * 1024 * 2)
	if err != nil {
		handleError(w, 500, err)
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	_, err = exampleClient.UploadImage(context.TODO(), &HOUSE.UploadImageRequest{
		Mobile: request["mobile"].(string),
		HouseId: request["house_id"].(string),
		Data: data,
		FileName: head.Filename,
		FileSize: head.Size,
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func GetHouseDesc(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.GetDesc(context.TODO(), &HOUSE.GetDescRequest{
		HouseId: p.ByName("house_id"),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	var data map[string]interface{}
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func SearchHouse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := HOUSE.NewExampleService("go.micro.srv.house", service.Client())
	rsp, err := exampleClient.Search(context.TODO(), &HOUSE.SearchRequest{
		AreaId: utils.GetParam("area_id", r),
		Start: utils.GetParam("start", r),
		End: utils.GetParam("end", r),
		Page: utils.GetParam("page", r),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	var data map[string]interface{}
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func AddOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "house_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "house_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "start_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "start_date cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName:"start_date",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"start_date invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName: "end_date",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "end_date cannot be empty",
			Trim: true,
			ValidEmpty: true,
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
			ErrMsg:"end date should be later than start date",
			Trim:true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	data, _ := json.Marshal(request)

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	_, err = exampleClient.Add(context.TODO(), &ORDER.AddRequest{
		Data: data,
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func GetOrderList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	rsp, err := exampleClient.GetList(context.TODO(), &ORDER.GetListRequest{
		Mobile: request["mobile"].(string),
		Role: utils.GetParam("role", r),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	var data []map[string]interface{}
	json.Unmarshal(rsp.Data, &data)

	success(w, data)
}

func HandleOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "order_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "order_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	_, err = exampleClient.Handle(context.TODO(), &ORDER.HandleRequest{
		Mobile: request["mobile"].(string),
		OrderId: request["order_id"].(string),
		Action: utils.GetParam("action", r),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}

func Comment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, 500, err)
		return
	}

	fvs := []*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "mobile cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "mobile",
			ValidMethodName: "ChineseMobile",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "invalid mobile",
			Trim: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "private_key cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "order_id",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "order_id cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
		&FormValidation.FieldValidation{
			FieldName: "comment",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg: "comment cannot be empty",
			Trim: true,
			ValidEmpty: true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		handleError(w, 4102, err)
		return
	}

	valid, err := validateUser(request["mobile"].(string), request["private_key"].(string))
	if err != nil {
		handleError(w, 4102, err)
		return
	}
	if !valid {
		handleError(w, 4102, errors.New("invalid private key"))
		return
	}

	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := ORDER.NewExampleService("go.micro.srv.order", service.Client())
	_, err = exampleClient.Comment(context.TODO(), &ORDER.CommentRequest{
		Mobile: request["mobile"].(string),
		OrderId: request["order_id"].(string),
		Comment: request["comment"].(string),
	})
	if err != nil {
		handleChaincodeError(w, err)
		return
	}

	success(w, nil)
}
