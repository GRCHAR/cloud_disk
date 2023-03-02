package controller

import (
	"cloud_disk/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var fileService service.FileService

type fileController struct {
	logger *zap.Logger
}

func NewFileController() *fileController {
	return new(fileController)
}

func (controller *fileController) CreateUploadFileHandler(c *gin.Context) {
	userId, err := c.Cookie("u_id")
	if err != nil {
		controller.logger.Error("CreateUploadFileHandler cookie err:", zap.Error(err))
		return
	}
	fileName, ok := c.GetQuery("fileName")
	if !ok {
		notFindQuery("fileName", c)
		return
	}
	fileSize, ok := c.GetQuery("fileSize")

	if !ok {
		notFindQuery("fileSize", c)
		return
	}
	fileSizeInt, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil {
		controller.logger.Error("CreateUploadFileHandler ParseInt err:", zap.Error(err))
		serverError(c)
		return
	}
	taskId, err := fileService.CreateUploadTask(fileName, fileSizeInt, userId)
	if err != nil {
		serverError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "",
		"data": map[string]string{
			"taskId": taskId,
		},
	})
}

func (controller *fileController) StartUploadFileTaskHandler(c *gin.Context) {
	//userId, err := c.Cookie("u_id")
	//taskId, ok := c.GetQuery("taskId")

}

func (controller *fileController) StopUploadFileTaskHandler(c *gin.Context) {

}

func (controller *fileController) UploadFileHandler(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		return
	}
	partNumber, ok := c.GetQuery("part_number")
	if !ok {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		controller.logger.Error("FormFile file", zap.Error(err))
		serverError(c)
		return
	}
	pieceNumber, _ := strconv.Atoi(partNumber)
	err = fileService.SaveFile(c, file, taskId, pieceNumber)
	if err != nil {
		controller.logger.Error("SaveFile err", zap.Error(err))
		serverError(c)
		return
	}
}

func (controller *fileController) DownloadFileHandler(c *gin.Context) {

}

func (controller *fileController) MoveFileHandler(c *gin.Context) {

}

func (controller *fileController) DeleteFileHandler(c *gin.Context) {

}

func notFindQuery(key string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    1,
		"message": "not found " + key,
	})
}

func serverError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    1,
		"message": "server error",
	})
}
