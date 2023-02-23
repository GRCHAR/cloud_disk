package controller

import "github.com/gin-gonic/gin"

type directoryController struct {
}

func NewDirectoryController() *directoryController {
	return new(directoryController)
}

func (controller *directoryController) CreateDirHandler(c *gin.Context) {

}

func (controller *directoryController) DeleteDirHandler(c *gin.Context) {

}
