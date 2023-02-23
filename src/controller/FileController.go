package controller

import (
	"cloud_disk/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var fileService service.FileService

type fileController struct {
	logger *zap.Logger
}

func NewFileController() *fileController {
	return new(fileController)
}

func (controller *fileController) CreateUploadFileHandler(c *gin.Context) {
	//fileName, ok := c.GetQuery("fileName")
	//if !ok {
	//	return
	//}
	//
	//fileType, ok := c.GetQuery("fileType")
	//if !ok {
	//	return
	//}
	//userId, err := c.Cookie("u_id")
	//if err != nil {
	//	controller.logger.Error("CreateUploadFileHandler cookie err:", zap.Error(err))
	//	return
	//}

}

func (controller *fileController) UploadFileHandler(c *gin.Context) {

}

func (controller *fileController) DownloadFileHandler(c *gin.Context) {

}

func (controller *fileController) MoveFileHandler(c *gin.Context) {

}

func (controller *fileController) DeleteFileHandler(c *gin.Context) {

}
