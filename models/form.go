package models


/**
业务 form
*/
type GetOneDataForm struct {
	Name  	string 	`form:"name" json:"name" binding:"required"`
	Version	string 	`form:"version" json:"version" binding:"required"`
	Os  	string 	`form:"os" json:"os" binding:"required"`
	Channel string 	`form:"channel" json:"channel" binding:"required"`
	Lat  	string 	`form:"lat" json:"lat"`
	Lng  	string 	`form:"lng" json:"lng"`
}

type GetVersionListDataForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Limit int `form:"limit" json:"limit" binding:"required"`
	Name string `form:"name" json:"name"`
	Code string `form:"code" json:"code"`
	Os string `form:"os" json:"os"`
}

type GetPackageListDataForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Limit int `form:"limit" json:"limit" binding:"required"`
	VersionId int `form:"version_id" json:"version_id" binding:"required"`
	Channel string 	`form:"channel" json:"channel"`
	Status int `form:"status" json:"status"`
}

type SaveVersionDataForm struct {
	Id int `form:"id" json:"id"`
	Name string `form:"name" json:"name" binding:"required"`
	Code string `form:"code" json:"code" binding:"required"`
	Os string `form:"os" json:"os" binding:"required"`
}

type SavePackageDataForm struct {
	Id int `form:"id" json:"id"`
	VersionId int `form:"version_id" json:"version_id" binding:"required"`
	Channel string 	`form:"channel" json:"channel" binding:"required"`
	Status int `form:"status" json:"status" binding:"required"`
}
