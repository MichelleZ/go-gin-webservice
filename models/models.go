package models

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/miaozhang/webservice/settings"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	var err error
	db, err = gorm.Open(settings.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DatabaseSetting.User,
		settings.DatabaseSetting.Password,
		settings.DatabaseSetting.Host,
		settings.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return settings.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimestampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimestampForUpdateCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

func updateTimestampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimestampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
