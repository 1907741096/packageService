package models

import (
	"time"
)
type AppVersion struct {
	AppVersionId       int       `orm:"app_version_id"`
	AppVersionName     string    `orm:"app_version_name"`      // 包名
	AppVersionCode     string    `orm:"app_version_code"`      // 版本
	AppVersionOs       string    `orm:"app_version_os"`        // 操作系统类型 android：安卓  iOS：iphone
	AppVersionCreateAt time.Time `orm:"app_version_create_at"` // 创建时间
	AppVersionUpdateAt time.Time `orm:"app_version_update_at"` // 修改时间
}

func (m *AppVersion) TableName() string {
	return "app_version"
}

// 查询列表
func (m *AppVersion) GetVersionList(page int, limit int, name string, code string, os string) ([]*AppVersion, int, error) {
	offset := (page - 1) * limit
	var model = []*AppVersion{}
	conditions := map[string]interface{}{}
	if name != "" {
		conditions["app_version_name"] = name
	}
	if code != "" {
		conditions["app_version_code"] = code
	}
	if os != "" {
		conditions["app_version_os"] = os
	}
	DB.Table(m.TableName()).
		Where(conditions).
		Offset(offset).
		Limit(limit).
		Order("app_version_id desc").
		Find(&model)

	var count int
	DB.Table(m.TableName()).
		Where(conditions).
		Count(&count)

	return model, count, nil;
}

// 增加记录
func  (m *AppVersion) CreateOne(name string, code string, os string) (*AppVersion, error) {
	model := AppVersion{
		AppVersionName: name,
		AppVersionCode: code,
		AppVersionOs: os,
		AppVersionCreateAt: time.Now(),
		AppVersionUpdateAt: time.Now(),
	}
	if err := DB.Table(m.TableName()).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

//修改记录
func  (m *AppVersion) UpdateOne(id int, name string, code string, os string) (*AppVersion, error) {
	conditions := map[string]interface{}{
		"app_version_id": id,
	}
	model := AppVersion{
		AppVersionName: name,
		AppVersionCode: code,
		AppVersionOs: os,
		AppVersionUpdateAt: time.Now(),
	}
	if err := DB.Table(m.TableName()).Where(conditions).Update(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
