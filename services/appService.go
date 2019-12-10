package services

import (
	"AS/commonlib-go"
	"AS/models"
	"AS/utils"
)

type App struct {
	Version     string
	Channel     string
	Package     string
	Lat			string
	Lng			string
}


type FormatDataReturnType map[string]interface{}

// 获取状态
func (b *App) GetStatus(name string, version string, os string, channel string, Lat string, Lng string) (FormatDataReturnType, error) {
	lib.Info("[monitor]http params message", lib.ErrorMessage{
		"name": name,
		"version": version,
		"os": os,
		"channel":channel,
	}, "monitor")
	appPackage := new(models.AppPackage)
	model, err := appPackage.GetOne(name, version, os, channel)
	lib.Info("[monitor]http model message", lib.ErrorMessage{
		"model": model,
	}, "monitor")
	if err != nil {
		return nil, err
	}
	data := FormatDataReturnType{
		"status" : model.AppPackageStatus,
	}
	return data, nil;
}

// 获取版本列表
func (b *App) GetVersionList(page int, limit int, name string, code string, os string) (FormatDataReturnType, error) {
	appVersion := new(models.AppVersion)
	models, count, err := appVersion.GetVersionList(page, limit, name, code, os)
	if err != nil {
		return nil, err
	}

	list := []FormatDataReturnType{}
	for _, item := range models {
		list = append(list, FormatDataReturnType{
			"id" : item.AppVersionId,
			"name" : item.AppVersionName,
			"code" : item.AppVersionCode,
			"os" : item.AppVersionOs,
			"create_at" : item.AppVersionCreateAt.Format(utils.DatetimeFormat),
			"update_at" : item.AppVersionUpdateAt.Format(utils.DatetimeFormat),
		})
	}

	data := FormatDataReturnType{
		"list" : list,
		"count" : count,
	}
	return data, nil;
}

// 获取包列表
func (b *App) GetPackageList(page int, limit int, versionId int, channel string, status int) (FormatDataReturnType, error) {
	appPackage := new(models.AppPackage)
	models, count, err := appPackage.GetList(page, limit, versionId, channel, status)
	if err != nil {
		return nil, err
	}

	list := []FormatDataReturnType{}
	for _, item := range models {
		list = append(list, FormatDataReturnType{
			"id" : item.AppPackageId,
			"channel" : item.AppPackageChannel,
			"status" : item.AppPackageStatus,
			"create_at" : item.AppPackageCreateAt.Format(utils.DatetimeFormat),
			"update_at" : item.AppPackageUpdateAt.Format(utils.DatetimeFormat),
		})
	}

	data := FormatDataReturnType{
		"list" : list,
		"count" : count,
	}
	return data, nil;
}

// 保存版本
func (b *App) SaveVersion(id int, name string, code string, os string) (FormatDataReturnType, error) {
	appVersion := new(models.AppVersion)
	if id != 0 {
		_, err := appVersion.UpdateOne(id, name, code, os)
		if err != nil {
			return nil, err
		}
		return FormatDataReturnType{}, nil;
	} else {
		_, err := appVersion.CreateOne(name, code, os)
		if err != nil {
			return nil, err
		}
		return FormatDataReturnType{}, nil;
	}
}

// 保存包
func (b *App) SavePackage(id int, versionId int, channel string, status int) (FormatDataReturnType, error) {
	appPackage := new(models.AppPackage)
	if id != 0 {
		_, err := appPackage.UpdateOne(id, channel, status)
		if err != nil {
			return nil, err
		}
		return FormatDataReturnType{}, nil;
	} else {
		_, err := appPackage.CreateOne(versionId, channel, status)
		if err != nil {
			return nil, err
		}
		return FormatDataReturnType{}, nil;
	}
}