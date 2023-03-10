package router

import "github.com/gin-gonic/gin"
import "cloud_disk/src/controller"

func init() {

}

func InitRouter() {

	r := gin.Default()
	fileController := controller.NewFileController()
	r.Group("/file")
	{
		r.POST("/createUploadTask", fileController.CreateUploadFileHandler)
		r.POST("/")

		r.POST("/download", fileController.DownloadFileHandler)
		r.POST("/move", fileController.MoveFileHandler)
		r.POST("/delete", fileController.DeleteFileHandler)
	}
	r.Group("/info")
	{

	}
	r.Group("/dir")
	{

	}

	r.Run(":8072")
}
