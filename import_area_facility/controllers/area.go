package controllers

import (
	"os"
	"encoding/csv"
	"190316/import_area_facility/models"
	"io"
	"strconv"
	"github.com/astaxie/beego"
	"encoding/json"
)

type AreaController struct {
	BaseController
}

func (this *AreaController) Add()  {
	file, err := os.Open("area.csv")
	if err != nil {
		this.handleResponse(500, err.Error(), nil)
		return
	}
	defer file.Close()

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.handleResponse(500, err.Error(), nil)
		return
	}

	reader := csv.NewReader(file)
	areaId := 0
	for  {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			this.handleResponse(500, err.Error(), nil)
			return
		}
		if len(line) < 1 {
			this.handleResponse(500, "invalid format", nil)
			return
		}

		areaId++
		txId, err := ccs.ChaincodeUpdate(models.ChaincodeId, "addArea", [][]byte{[]byte(strconv.Itoa(areaId)), []byte(line[0])})
		if err != nil {
			this.handleResponse(500, err.Error(), nil)
			return
		}
		beego.Info(string(txId))
	}

	this.handleResponse(200, "ok", nil)
}

func (this *AreaController) List() {
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.handleResponse(500, err.Error(), nil)
		return
	}

	_data, err := ccs.ChaincodeQuery(models.ChaincodeId, "getAreaList", [][]byte{})
	if err != nil {
		this.handleResponse(500, err.Error(), nil)
		return
	}

	var data []map[string]interface{}
	err = json.Unmarshal(_data, &data)
	if err != nil {
		this.handleResponse(500, err.Error(), nil)
		return
	}

	this.handleResponse(200, "ok", data)
}
