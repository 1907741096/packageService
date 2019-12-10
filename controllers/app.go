package controllers

import (
	"AS/services"
	"AS/models"
	"AS/utils"
	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	params := c.MustGet(utils.BindFormKey).(*models.GetOneDataForm)
	result, err := new(services.App).GetStatus(params.Name, params.Version, params.Os, params.Channel, params.Lat, params.Lng)
	if err != nil {
		utils.FormatFail(c, err)
		return
	}
	utils.FormatSuccess(c, result)
}

func GetVersionList(c *gin.Context) {
	params := c.MustGet(utils.BindFormKey).(*models.GetVersionListDataForm)
	result, err := new(services.App).GetVersionList(params.Page, params.Limit, params.Name, params.Code, params.Os)
	if err != nil {
		utils.FormatFail(c, err)
		return
	}
	utils.FormatSuccess(c, result)
}

func GetPackageList(c *gin.Context) {
	params := c.MustGet(utils.BindFormKey).(*models.GetPackageListDataForm)
	result, err := new(services.App).GetPackageList(params.Page, params.Limit, params.VersionId, params.Channel, params.Status)
	if err != nil {
		utils.FormatFail(c, err)
		return
	}
	utils.FormatSuccess(c, result)
}

func SaveVersion(c *gin.Context) {
	params := c.MustGet(utils.BindFormKey).(*models.SaveVersionDataForm)
	result, err := new(services.App).SaveVersion(params.Id, params.Name, params.Code, params.Os)
	if err != nil {
		utils.FormatFail(c, err)
		return
	}
	utils.FormatSuccess(c, result)
}

func SavePackage(c *gin.Context) {
	params := c.MustGet(utils.BindFormKey).(*models.SavePackageDataForm)
	result, err := new(services.App).SavePackage(params.Id, params.VersionId, params.Channel, params.Status)
	if err != nil {
		utils.FormatFail(c, err)
		return
	}
	utils.FormatSuccess(c, result)
}