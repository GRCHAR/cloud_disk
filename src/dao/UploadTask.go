package dao

import (
	"cloud_disk/src/logger"
	"github.com/sirupsen/logrus"
	"math"
)

type UploadTask struct {
	Id         int    `gorm:"primary_key;index:id_idx"`
	FileId     string `gorm:"type:varchar(255)"`
	ParentId   int    `gorm:"type:int"`
	FileName   string `gorm:"type:varchar(255)"`
	FileSize   int64  `gorm:"type:bigint"`
	PartNumber int    `gorm:"type:int"`
	UserId     int    `gorm:"type:int"`
	Status     string `gorm:"type:varchar(255)"`

	//UploadedPart []int `gorm:"type:varchar(255)"`
}

type UploadedPart struct {
	Id           string `gorm:"primary_key;index:id_idx"`
	UploadTaskId string `gorm:"primary_key;index:id_idx"`
	PartNo       int    `gorm:"type:int""`
}

type UploadTaskDao struct {
	Dao
}

func GetUploadDao() *UploadTaskDao {
	dao := new(UploadTaskDao)
	dao.logger = logger.GetLogger()
	return dao
}

func (dao *UploadTaskDao) CreateUploadTask(file File) (int, error) {
	fileName := file.Name
	fileSize := file.FileLength
	partTotalNumber := math.Ceil(float64(fileSize / (2 * 1024)))
	uploadTask := UploadTask{FileName: fileName, PartNumber: int(partTotalNumber)}
	if err := db.Create(&uploadTask).Error; err != nil {
		dao.logger.WithFields(logrus.Fields{
			"id":   file.Id,
			"name": file.Name,
			"user": file.CreateUser,
		}).Error("创建上传任务失败!", err)
		return -1, err
	}
	return uploadTask.Id, nil
}

func (dao *UploadTaskDao) FindUploadTaskById(id int) (UploadTask, error) {

	uploadTask := UploadTask{}
	if err := db.Model(&uploadTask).Where("id=?", int64(id)).Error; err != nil {
		dao.logger.Error(err)
		return uploadTask, err
	}
	return uploadTask, nil
}

func (dao *UploadTaskDao) UpdateUploadTask(task UploadTask) (UploadTask, error) {
	uploadTask := &UploadTask{Id: task.Id}
	if err := db.Model(uploadTask).Update("status", task.Status).Error; err != nil {
		dao.logger.Error(err)
		return *uploadTask, err
	}
	return *uploadTask, nil
}
