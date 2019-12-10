package models

import "time"

type app_package struct {
	AppPackageId           int       `orm:"app_package_id"`
	AppPackageAppVersionId string    `orm:"app_package_app_version_id"` // 版本
	AppPackageChannel      string    `orm:"app_package_channel"`        // 渠道
	AppPackageStatus       int       `orm:"app_package_status"`         // 状态 0-审核中，1-最新版本 2-非最新版本 3-不可用版本
	AppPackageCreateAt     time.Time `orm:"app_package_create_at"`      // 创建时间
	AppPackageUpdateAt     time.Time `orm:"app_package_update_at"`      // 修改时间
}

type app_version struct {
	AppVersionId       int       `orm:"app_version_id"`
	AppVersionName     string    `orm:"app_version_name"`      // 包名
	AppVersionCode     string    `orm:"app_version_code"`      // 版本
	AppVersionOs       string    `orm:"app_version_os"`        // 操作系统类型 android：安卓  iOS：iphone
	AppVersionCreateAt time.Time `orm:"app_version_create_at"` // 创建时间
	AppVersionUpdateAt time.Time `orm:"app_version_update_at"` // 修改时间
}
