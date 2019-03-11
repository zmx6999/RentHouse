package controllers

import (
	"os"
	"190305/models"
	"encoding/csv"
	"io"
	"github.com/pkg/errors"
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
		this.handleError(500, err)
		return
	}
	defer file.Close()

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.handleError(500, err)
		return
	}

	reader := csv.NewReader(file)
	line := 0
	for  {
		line++
		lineStr := strconv.Itoa(line)
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			this.handleError(500, err)
			return
		}
		if len(row) < 1 {
			this.handleError(500, errors.New("invalid line " + lineStr))
			return
		}

		txId, err := ccs.ChaincodeUpdate(models.ChaincodeId, "addArea", [][]byte{[]byte(lineStr), []byte(row[0])})
		if err != nil {
			this.handleError(500, err)
			return
		}
		beego.Info(string(txId))
	}

	this.handleSuccess(nil)
}

func (this *AreaController) Get()  {
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.handleError(500, err)
		return
	}

	data, err := ccs.ChaincodeQuery(models.ChaincodeId, "getAreaList", [][]byte{})
	if err != nil {
		this.handleError(500, err)
		return
	}

	var areaList []map[string]interface{}
	err = json.Unmarshal(data, &areaList)
	if err != nil {
		this.handleError(500, err)
		return
	}

	this.handleSuccess(areaList)
}