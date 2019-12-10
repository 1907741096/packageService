package lib

import (
	"errors"
)

const (
	SystemInterError = 100000 + iota
	DatabaseError
	ParamInvalid
	ObjectNotFound
	FeatureNotFound
	GetRiskFeatureFail
	GetDsqFeatureFail
	TransTypeFail
	FeatureMutex
	ConfigFileNotExist
)

var (
	ErrSystemInter       = errors.New("系统内部错误")
	ErrDataBase          = errors.New("数据库操作错误")
	ErrParamInvalid      = errors.New("接口参数不合法")
	ErrNotFound		     = errors.New("数据对象不存在")
	ErrFeatureNotFound   = errors.New("用户特征不存在")
	ErrRiskFeatureFail	 = errors.New("风控用户特征查询失败")
	ErrDsqFeatureFail    = errors.New("带上钱用户特征查询失败")
	ErrTypeTransErr      = errors.New("类型转化异常")
	ErrFeatureMutex 	 = errors.New("获取不可转化的特征类型")
	ErrConfigFileNotExist = errors.New("配置文件加载失败")
)

var ErrCodeMap = map[error]int{
	ErrSystemInter:       SystemInterError,
	ErrDataBase:          DatabaseError,
	ErrParamInvalid:      ParamInvalid,
	ErrNotFound:		  ObjectNotFound,
	ErrFeatureNotFound:   FeatureNotFound,
	ErrRiskFeatureFail:   GetRiskFeatureFail,
	ErrDsqFeatureFail:    GetDsqFeatureFail,
	ErrTypeTransErr: 	  TransTypeFail,
	ErrFeatureMutex: 	  FeatureMutex,
	ErrConfigFileNotExist: ConfigFileNotExist,
}
