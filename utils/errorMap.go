package utils

import (
	"errors"
)

const (
	SystemInterError    = 100001
	DatabaseError       = 100002
	ParamInvalid        = 100003
	ObjectNotFound      = 100004
	ConfigFileNotExist  = 100005
	ResponseFail        = 100006
)

var (
	ErrSystemInter         = errors.New("系统内部错误")
	ErrDataBase            = errors.New("数据库操作错误")
	ErrParamInvalid        = errors.New("接口参数不合法")
	ErrNotFound            = errors.New("数据对象不存在")
	ErrConfigFileNotExist  = errors.New("配置文件加载失败")
	ErrResponseFail        = errors.New("系统繁忙")
)

var ErrCodeMap = map[error]int{
	ErrSystemInter:         SystemInterError,
	ErrDataBase:            DatabaseError,
	ErrParamInvalid:        ParamInvalid,
	ErrNotFound:            ObjectNotFound,
	ErrConfigFileNotExist:  ConfigFileNotExist,
	ErrResponseFail:        ResponseFail,
}
