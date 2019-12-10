package models

import (
	"AS/utils"
	"time"
)

type AppPackage struct {
	AppPackageId           int       `orm:"app_package_id"`
	AppPackageAppVersionId int    	 `orm:"app_package_app_version_id"` // 版本
	AppPackageChannel      string    `orm:"app_package_channel"`        // 渠道
	AppPackageStatus       int       `orm:"app_package_status"`         // 状态 0-审核中，1-最新版本 2-非最新版本 3-不可用版本
	AppPackageCreateAt     time.Time `orm:"app_package_create_at"`      // 创建时间
	AppPackageUpdateAt     time.Time `orm:"app_package_update_at"`      // 修改时间
}

func (m *AppPackage) TableName() string {
	return "app_package"
}

// 根据条件查询
func (m *AppPackage) GetOne(name string, code string, os string, channel string) (*AppPackage, error) {
	var model = AppPackage{}
	if (DB.
		Model(&AppPackage{}).
		Joins("left join app_version on app_version_id = app_package_app_version_id").
		Where("app_version_name = ?", name).
		Where("app_version_code = ?", code).
		Where("app_version_os = ?", os).
		Where("app_package_channel = ?", channel).
		Order("app_package_id desc").
		Find(&model).RecordNotFound()) {
		return nil, utils.ErrNotFound
	}
	return &model, nil;
}

// 查询列表
func (m *AppPackage) GetList(page int, limit int,versionId int, channel string, status int) ([]*AppPackage, int, error) {
	offset := (page - 1) * limit
	var model = []*AppPackage{}
	conditions := map[string]interface{}{
		"app_package_app_version_id": versionId,
	}
	if channel != "" {
		conditions["app_package_channel"] = channel
	}
	if status != 0 {
		conditions["app_package_status"] = status
	}
	DB.Table(m.TableName()).
		Where(conditions).
		Offset(offset).
		Limit(limit).
		Order("app_package_id desc").
		Find(&model)

	var count int
	DB.Table(m.TableName()).
		Where(conditions).
		Count(&count)

	return model, count, nil;
}

// 增加记录
func  (m *AppPackage) CreateOne(versionId int, channel string, status int) (*AppPackage, error) {
	model := AppPackage{
		AppPackageAppVersionId: versionId,
		AppPackageChannel: channel,
		AppPackageStatus: status,
		AppPackageCreateAt: time.Now(),
		AppPackageUpdateAt: time.Now(),
	}
	if err := DB.Table(m.TableName()).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

//修改记录
func  (m *AppPackage) UpdateOne(id int, channel string, status int) (*AppPackage, error) {
	model := AppPackage{
		AppPackageStatus: status,
		AppPackageChannel: channel,
		AppPackageUpdateAt: time.Now(),
	}
	conditions := map[string]interface{}{
		"app_package_id": id,
	}
	if err := DB.Table(m.TableName()).Where(conditions).Update(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}