package dao

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Dao struct {
	logger zap.Logger
}

var DB *gorm.DB

func (dao *Dao) initDao() {
	db, err := gorm.Open("mysql", "root:123456@/cloud?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		dao.logger.Error("连接数据失败!", zap.Error(err))
		return
	}
	db.AutoMigrate(&File{}, &UploadTask{}, &User{})
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			dao.logger.Error("创建User表失败")
			return
		}
	}
	if !db.HasTable(&File{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&File{}).Error; err != nil {
			dao.logger.Error("创建File表失败")
		}
	}
}
