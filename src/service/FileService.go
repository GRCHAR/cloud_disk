package service

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type FileService struct {
	rawService
}

type uploadTask struct {
}

type uploadTaskPool struct {
}

func (service *FileService) StartUploadTask(c *gin.Context, id string) {

}

func (service *FileService) StopUploadTask() {

}

func (service *FileService) ContinueUploadTask() {

}

func (service *FileService) DeleteUploadTask() {

}

func (service *FileService) SaveFile(c *gin.Context, file multipart.File, fileHeader *multipart.FileHeader, path string) error {

	return nil
}
