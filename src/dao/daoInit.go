package dao

import (
	"cloud_disk/src/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type Dao struct {
	logger *logrus.Logger
}

var db *gorm.DB

func GetDao() *Dao {
	dao := new(Dao)
	dao.logger = logger.GetLogger()
	return dao
}

func GetDB() *gorm.DB {
	return db
}

func init() {
	dao := GetDao()
	dao.initMysql()
	dao.initRedis()
}

func (dao *Dao) initMysql() {
	DB, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/cloud?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		dao.logger.Error("连接数据失败!")
		return
	}
	db = DB
	db.AutoMigrate(&File{}, &UploadTask{}, &User{}, &UploadedPart{}, &Dir{}, &DownloadTask{})
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			dao.logger.Error("创建User表失败")
			return
		}
	}
	if !db.HasTable(&UploadTask{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&UploadTask{}).Error; err != nil {
			dao.logger.Error("创建UploadTask表失败")
		}
	}

	if !db.HasTable(&File{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&File{}).Error; err != nil {
			dao.logger.Error("创建File表失败")
		}
	}

	if !db.HasTable(&UploadedPart{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&UploadedPart{}).Error; err != nil {
			dao.logger.Error("创建UploadedPart表失败")
		}
	}
	if !db.HasTable(&Dir{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Dir{}).Error; err != nil {
			dao.logger.Error("创建Dir表失败")
		}
	}
	if !db.HasTable(&DownloadTask{}) {
		if err := db.Set("gorm:table options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&DownloadTask{}).Error; err != nil {
			dao.logger.Error("创建DownloadTask表失败")
		}
	}
	dao.logger.Info("数据库连接成功")
}

func (dao *Dao) initRedis() {
}

func GetGormDB() *gorm.DB {
	return db
}
