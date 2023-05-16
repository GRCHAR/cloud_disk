package service

import "cloud_disk/src/dao"

type downloadTaskService struct {
	rawService
}

var downloadTaskDao = dao.NewDownloadTaskDao()

func (service *downloadTaskService) CreateDownloadTask(fileId int64, userId int64) (int64, bool) {
	taskId, err := downloadTaskDao.CreateDownloadTask(fileId, userId)
	if err != nil {
		service.logger.Error("Create download task error:%+v", err)
		return -1, false
	}
	return taskId, true
}

func (service *downloadTaskService) StartDownloadTask(taskId int64) bool {
	err := downloadTaskDao.UpdateDownloadTaskStatus(taskId, "run")
	if err != nil {
		service.logger.Error("start downloadTask error:", err)
		return false
	}
	return true
}

func (service *downloadTaskService) StopDownloadTask(taskId int64) bool {
	err := downloadTaskDao.UpdateDownloadTaskStatus(taskId, "stop")
	if err != nil {
		service.logger.Error("start downloadTask error:", err)
		return false
	}
	return true
}

func (service *downloadTaskService) DeleteDownloadTask(taskId int64) bool {
	err := downloadTaskDao.DeleteDownloadTask(taskId)
	if err != nil {
		service.logger.Error("delete downloadTask error:", err)
		return false
	}
	return true
}
