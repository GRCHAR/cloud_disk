package dao

import "cloud_disk/src/logger"

type DownloadTask struct {
	id      int64
	fileId  int64
	userId  int64
	status  string
	process float64
}

type DownloadTaskDao struct {
	Dao
}

func NewDownloadTaskDao() *DownloadTaskDao {
	downloadTaskDao := DownloadTaskDao{Dao{logger: logger.GetLogger()}}
	return &downloadTaskDao
}

func (dtd *DownloadTaskDao) CreateDownloadTask(fileId int64, userId int64) (int64, error) {
	downloadTask := DownloadTask{
		status:  "stop",
		process: 0,
		fileId:  fileId,
		userId:  userId,
	}
	err := db.Create(&downloadTask).Error
	if err != nil {
		return -1, err
	}
	return downloadTask.id, nil
}

func (dtd *DownloadTaskDao) UpdateDownloadTaskStatus(id int64, status string) error {
	err := db.Model(&DownloadTask{}).Where("id =?", id).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (dtd *DownloadTaskDao) DeleteDownloadTask(id int64) error {
	err := db.Where("id=?", id).Delete(&DownloadTask{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dtd *DownloadTaskDao) GetDownloadTasks(userId int64, page int) ([]DownloadTask, error) {
	var downloadTasks []DownloadTask
	err := db.Limit(20).Offset(page).Where("userId=?", userId).Find(&downloadTasks).Error
	if err != nil {
		return downloadTasks, err
	}
	return downloadTasks, nil
}
