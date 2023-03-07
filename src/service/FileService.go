package service

import (
	"cloud_disk/src/config"
	"cloud_disk/src/dao"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"strconv"
)

type FileService struct {
	rawService
}

type uploadTask struct {
}

type uploadTaskPool struct {
}

var uploadTaskDao *dao.UploadTaskDao
var userDao *dao.UserDao
var fileDao *dao.FileDao

func (service *FileService) CreateUploadTask(fileName string, fileSize int64, userId string) (taskId string, err error) {

	uploadUser, err := userDao.FindUserById(userId)
	if err != nil {
		service.logger.Error("CreateUploadTask FindUserById error:", zap.Error(err))
		return "", err
	}
	uploadFile := dao.File{FileLength: fileSize, Name: fileName, CreateUser: uploadUser.Id}
	taskId, err = uploadTaskDao.CreateUploadTask(uploadFile)
	if err != nil {
		service.logger.Error("CreateUploadTask CreateUploadTask error:", zap.Error(err))
		return "", err
	}
	if !service.StartUploadTask(taskId) {
		return "", errors.New("start upload task failed")
	}
	return taskId, nil
}

func (service *FileService) StartUploadTask(taskId string) bool {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("FindUploadTaskById err", zap.Error(err))
		return false
	}
	task.Status = "uploading"
	task, err = uploadTaskDao.UpdateUploadTask(task)
	if err != nil {
		service.logger.Error("UpdateUploadTask err", zap.Error(err))
		return false
	}
	go service.runTask(task.Id)
	return true
}

func (service *FileService) StopUploadTask() {

}

func (service *FileService) ContinueUploadTask() {

}

func (service *FileService) DeleteUploadTask() {

}

func (service *FileService) GetTaskInfo(userId string) {

}

func (service *FileService) SaveFile(c *gin.Context, header *multipart.FileHeader, taskId string, partNumber int) error {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		return err
	}
	uploadTempDir := config.GetConfig().Server.TempDirPath

	err = c.SaveUploadedFile(header, uploadTempDir+"/"+task.UserId+"/"+task.Id+"/"+strconv.Itoa(partNumber))

	if err != nil {
		service.logger.Error("SaveUploadedFile err", zap.Error(err))
		return err
	}
	return nil
}

func (service *FileService) runTask(taskId string) error {
	return nil
}

func (service *FileService) MergeFile(taskId string) error {
	go func() {
		uploadTempDir := config.GetConfig().Server.TempDirPath
		resultDir := config.GetConfig().Server.MergeDirPath
		task, err := uploadTaskDao.FindUploadTaskById(taskId)
		if err != nil {
			service.logger.Error("FindUploadTaskById err", zap.Error(err))
			return
		}
		fs := os.DirFS(uploadTempDir + "/" + task.UserId + task.Id)
		resultFile, err := os.Create(resultDir + "/" + task.UserId + "/" + task.FileId)
		if err != nil {
			service.logger.Error("os.Create err", zap.Error(err))

			return
		}
		var tempByte []byte
		for i := 0; i < task.PartNumber; i++ {
			tempByte = []byte{}
			file, err := fs.Open(strconv.Itoa(i))
			if err != nil {
				service.logger.Error("fs.Open err", zap.Error(err))
				return
			}
			_, err = file.Read(tempByte)
			if err != nil {
				service.logger.Error("file.Read err", zap.Error(err))
				return
			}
			resultFile.WriteAt(tempByte, int64(1024*i))
		}
		return
	}()
	return nil
}

func (service *FileService) GetTaskStatus(taskId string) (status string, err error) {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("uploadTaskDao.FindUploadTaskById err,", zap.Error(err))
		return "", err
	}
	return task.Status, nil
}
