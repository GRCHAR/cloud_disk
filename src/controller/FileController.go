package controller

import "github.com/gin-gonic/gin"

type fileController struct {
}

func NewFileController() *fileController {
	return new(fileController)
}

func (controller *fileController) UploadFileHandler(c *gin.Context) {

}

func (controller *fileController) DownloadFileHandler(c *gin.Context) {

}

func (controller *fileController) MoveFileHandler(c *gin.Context) {

}

func (controller *fileController) DeleteFileHandler(c *gin.Context) {

}
