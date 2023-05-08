package controller

import (
	"cloud_disk/src/logger"
	"cloud_disk/src/service"
	"cloud_disk/src/vo"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type directoryController struct {
	logger *logrus.Logger
}

var dirService = service.NewDirectoryService()

func init() {

}

func NewDirectoryController() *directoryController {
	directoryController := &directoryController{logger: logger.GetLogger()}
	return directoryController
}

func (controller *directoryController) CreateDirHandler(c *gin.Context) {
	_, err := c.Cookie("u_id")
	if err != nil {

		c.JSON(500, "no u_id in the cookies")
		return
	}

	parentId, ok := c.GetQuery("parentId")
	if !ok {
		controller.logger.Error("no parentId in the request", err)
		c.JSON(500, "no parentId in the request")
		return
	}
	dirName, ok := c.GetQuery("dirName")
	if !ok {
		controller.logger.Error("no dirName in the request", err)
		c.JSON(500, "no dirName in the request")
		return
	}
	parentIdInt, _ := strconv.Atoi(parentId)
	err = dirService.CreateDir(int64(parentIdInt), dirName)
	if err != nil {
		controller.logger.Error("Error while creating directory", err)
		c.JSON(500, "Error while creating directory")
		return
	}
	return
}

func (controller *directoryController) DeleteDirHandler(c *gin.Context) {

}

func (controller *directoryController) GetDirFilesHandler(c *gin.Context) {
	_, err := c.Cookie("u_id")
	if err != nil {

		c.JSON(500, "no u_id in the cookies")
		return
	}
	dirId, ok := c.GetQuery("dirId")
	if !ok {
		controller.logger.Error("no dirId in the request", err)
		c.JSON(500, "no dirId in the request")
		return
	}
	dirIdInt, _ := strconv.Atoi(dirId)
	result, err := dirService.GetDirContent(int64(dirIdInt))
	if err != nil {
		controller.logger.Error("Error while getting directory files", err)
		c.JSON(500, "Error while getting directory files")
		return
	}
	vo.ResponseDataSuccess(result, c)
}
