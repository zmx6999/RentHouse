package handler

import (
	"context"

		example "sss/GetHouses/proto/example"
	"sss/181231/utils"
	"sss/181231/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"math"
	"encoding/json"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	rsp.ErrCode=utils.RECODE_OK
	rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)

	var _houses []*models.House
	o:=orm.NewOrm()
	qs:=o.QueryTable("House")
	if req.AreaId!="" {
		qs=qs.RelatedSel("Area").Filter("Area__Id",req.AreaId)
	}
	totalRows,_:=qs.Count()
	pageSize:=5
	totalPages:=int(math.Ceil(float64(totalRows)/float64(pageSize)))
	page,err:=strconv.Atoi(req.Page)
	if err!=nil {
		page=1
	}
	offset:=pageSize*(page-1)
	_,err=qs.Limit(pageSize,offset).All(&_houses)
	if err!=nil {
		rsp.ErrCode=utils.RECODE_NODATA
		rsp.ErrMsg=utils.RecodeTest(rsp.ErrCode)
		return nil
	}

	var houses []map[string]interface{}
	for _,house:=range _houses{
		o.LoadRelated(house,"Area")
		o.LoadRelated(house,"User")
		houses=append(houses,house.Info())
	}
	data,_:=json.Marshal(houses)
	rsp.Data=data
	rsp.TotalPages=int64(totalPages)
	rsp.Page=int64(page)
	return nil
}
