package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

type migrateIns struct {
	migrate *migrate.Migrate
}

var m *migrateIns

func NewMigration(migrationFile, dsn string) *migrateIns {
	if m != nil {
		return m
	}
	if err := initialize(migrationFile, dsn); err != nil {
		panic(err)
	}
	return m
}

func initialize(migrationFile, dsn string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s&multiStatements=true", dsn))
	if err != nil {
		panic("migration failed. message:" + err.Error())
	}
	defer db.Close()

	// 确认数据库连接正常 && 检查migration 表是否存在
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic("Could not start migration! message:" + err.Error())
	}

	filePath := fmt.Sprintf("file://%s", migrationFile)
	// 根据驱动，解析migration 读取方式
	ins, err := migrate.NewWithDatabaseInstance(filePath, "mysql", driver)
	if err != nil {
		panic("migration failed! message:" + err.Error())
	}

	m = &migrateIns{migrate:ins}
	return nil
}

func (ins *migrateIns) Up() (err error) {
	if err := ins.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		panic("An error occurred while syncing the database! message:" + err.Error())
	}

	log.Println("Database migrated!")
	return
}

func (ins *migrateIns) Force() (err error) {
	v, dirty, err := ins.migrate.Version()
	if v == 0 {
		panic("this version not found")
	}

	if err != nil {
		panic("database occur something wrong. message:" + err.Error())
	}

	if dirty {
		if err := ins.migrate.Force(int(v)); err != nil {
			panic("An error occurred while executing the migration. message:" + err.Error())
		}
		log.Printf("force update dirty version %d execution succeeded", v)
	} else {
		log.Printf("the current version %d is not dirty", v)
	}
	return nil
}