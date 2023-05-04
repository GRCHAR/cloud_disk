package router

import "github.com/gin-gonic/gin"
import "cloud_disk/src/controller"
import "cloud_disk/src/logger"

func init() {

}

func InitRouter() {

	r := gin.Default()

	fileController := controller.NewFileController(logger.GetLogger())
	fileGroup := r.Group("/file")
	{
		fileGroup.POST("/createUploadTask", fileController.CreateUploadFileHandler)
		fileGroup.POST("/upload", fileController.UploadFileHandler)
		fileGroup.POST("/download", fileController.DownloadFileHandler)
		fileGroup.POST("/move", fileController.MoveFileHandler)
		fileGroup.POST("/delete", fileController.DeleteFileHandler)
		fileGroup.GET("/start", fileController.StartUploadFileTaskHandler)
	}
	r.Group("/info")
	{

	}
	r.Group("/dir")
	{

	}

	r.Run("0.0.0.0:8075")
}

func checkFile(c *gin.Context) {

}
