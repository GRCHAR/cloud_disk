package controller

import (
	"cloud_disk/src/service"
	"cloud_disk/src/vo"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var fileService *service.FileService

type fileController struct {
	logger *logrus.Logger
}

func init() {
	fileService = service.GetFileService()
}

func NewFileController(logger *logrus.Logger) *fileController {
	controller := new(fileController)
	controller.logger = logger
	return controller
}

func (controller *fileController) CreateUploadFileHandler(c *gin.Context) {
	controller.logger.Info("create")
	userId, err := c.Cookie("u_id")
	if err != nil {
		controller.logger.Error("CreateUploadFileHandler cookie err:", err)
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
		controller.logger.Error("CreateUploadFileHandler ParseInt err:")
		serverError(c)
		return
	}
	taskId, err := fileService.CreateUploadTask(fileName, fileSizeInt, userId)
	if err != nil {
		serverError(c)
		return
	}
	vo.ResponseDataSuccess(gin.H{"taskId": taskId}, c)

}

func (controller *fileController) MergeUploadFileHandler(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		notFindQuery("taskId", c)
	}
	//userId, err := c.Cookie("u_id")
	//if err != nil {
	//	controller.logger.Error("cookie not found cookie", )
	//	serverError(c)
	//	return
	//}
	taskIdInt, _ := strconv.Atoi(taskId)
	err := fileService.MergeFile(taskIdInt)
	if err != nil {
		serverError(c)
		return
	}
	vo.ResponseDataSuccess(gin.H{}, c)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "开始合并",
	})

}

func (controller *fileController) StartUploadFileTaskHandler(c *gin.Context) {
	userId, err := c.Cookie("u_id")
	if err != nil {
		controller.logger.Error("未在cookie中找到u_id")
		notFindCookie("u_id", c)
		return
	}
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		controller.logger.Error("未找到taskId参数")
		notFindQuery("taskId", c)
		return
	}
	userIdInt, _ := strconv.Atoi(userId)
	taskIdInt, _ := strconv.Atoi(taskId)
	if err = fileService.RunUploadTask(taskIdInt, userIdInt); err != nil {
		vo.ResponseDataFail(c, err)
		return
	}
	fileService.RunMergeTask(taskIdInt)
	vo.ResponseWithoutData(c, http.StatusOK, 0, "开始任务成功")
}

func (controller *fileController) StopUploadFileTaskHandler(c *gin.Context) {

}

func (controller *fileController) GetUploadFileTaskHandler(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		notFindQuery("taskId", c)
		return
	}
	taskIdInt, _ := strconv.Atoi(taskId)
	status, err := fileService.GetTaskStatus(taskIdInt)
	if err != nil {
		serverError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "",
		"data": map[string]string{
			"status": status,
		},
	})
}

func (controller *fileController) UploadFileHandler(c *gin.Context) {
	taskId, ok := c.GetQuery("taskId")
	if !ok {
		serverError(c)
		return
	}
	partNumber, ok := c.GetQuery("part_number")
	if !ok {
		serverError(c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		controller.logger.Error("FormFile file")
		serverError(c)
		return
	}
	taskIdInt, _ := strconv.Atoi(taskId)
	pieceNumber, _ := strconv.Atoi(partNumber)
	fileService.AppendUploadFileChannel(c, file, taskIdInt, pieceNumber)

	if err != nil {
		controller.logger.Error("SaveFile err")
		serverError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "done",
		"data": map[string]int{
			"pieceNumber": pieceNumber,
		},
	})
}

func (controller *fileController) DownloadFileHandler(c *gin.Context) {

}

func (controller *fileController) MoveFileHandler(c *gin.Context) {

}

func (controller *fileController) DeleteFileHandler(c *gin.Context) {

}

func (controller *fileController) GetFileTreeHandler(c *gin.Context) {

}

func notFindCookie(key string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    1,
		"message": "not found cookie" + key,
	})
}

func notFindQuery(key string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    1,
		"message": "not found query" + key,
	})
}

func serverError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    1,
		"message": "server error",
	})
}
