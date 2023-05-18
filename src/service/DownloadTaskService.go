package service

import (
	"cloud_disk/src/config"
	"cloud_disk/src/dao"
	"cloud_disk/src/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

type downloadTaskQueueObject struct {
	fileId  int64
	taskId  int64
	context *gin.Context
}

type downloadTaskService struct {
	rawService
	downloadTaskQueue chan downloadTaskQueueObject
}

var downloadTaskDao = dao.NewDownloadTaskDao()

func newDownloadService() *downloadTaskService {
	service := downloadTaskService{
		rawService:        rawService{logger: logger.GetLogger()},
		downloadTaskQueue: make(chan downloadTaskQueueObject, 100000),
	}
	return &service
}

func init() {
	service := newDownloadService()
	go service.downloadFileQueueHandler()
}

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

func (service *downloadTaskService) downloadFile(fileId int64, taskId int64, c *gin.Context) {
	file, err := fileDao.FindFile(fileId)
	if err != nil {
		service.logger.Error("downloadFile error:", err)
		return
	}
	if !service.StartDownloadTask(taskId) {
		service.logger.Error("downloadFile startDownloadTask error")
		return
	}

	mergePath := config.GetConfig().Server.MergeDirPath
	c.FileAttachment(mergePath+"/"+strconv.Itoa(int(fileId)), file.Name)
}

func (service *downloadTaskService) downloadPartFile(fileId int64, c *gin.Context) {
	//thread_number := config.GetConfigInfo().DownloadTaskThreadNumber
	//file, err := fileDao.FindFile(fileId)
	//if err != nil {
	//	service.logger.Error("downloadPartFile error:", err)
	//	return
	//}
	//tempPath := config.GetConfig().Server.TempDirPath

}

func (service *downloadTaskService) downloadFileQueueHandler() {
	for {
		select {
		case dtqo := <-service.downloadTaskQueue:
			service.downloadFile(dtqo.fileId, dtqo.taskId, dtqo.context)
		default:
			continue
		}
	}

}
