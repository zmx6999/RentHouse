package controllers

import (
		"190720/models"
		"encoding/json"
)

type AreaController struct {
	BaseController
}

func (this *AreaController) GetList()  {
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	payload, err := ccs.ChaincodeQuery("getAreaList", [][]byte{})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal(payload, &data)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(data)
}
