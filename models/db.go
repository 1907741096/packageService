package models

import (
	lib "AS/commonlib-go"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	InitDB()
}

func InitDB() {
	// 初始化DB
	DB = lib.RegisterDB(lib.DBConfig)
}

type TransFun func(tx *gorm.DB) error

func Transaction(closures ...TransFun) (err error) {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			switch r.(type) {
			case error:
				err = r.(error)
			case string:
				err = errors.New(r.(string))
			default:
				err = errors.New("system internal error")
			}
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	for _, closure := range closures {
		if err := closure(tx); err != nil {
			tx.Rollback()
			return err
		}
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func MySQLClose() {
	if DB != nil {
		_ = DB.Close()
	}
}