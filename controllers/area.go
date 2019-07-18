package controllers

import (
		"190702m/models"
				"encoding/json"
)

type AreaController struct {
	BaseController
}

func (this *AreaController) GetList()  {
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId,models.ConfigFile)
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	_data, err := ccs.ChaincodeQuery("getAreaList", [][]byte{})
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal(_data, &data)
	if err != nil {
		this.error(1006, err.Error())
		return
	}

	this.success(data)
}
