package controller

import (
	"cloud_disk/src/logger"
	"cloud_disk/src/service"
	"cloud_disk/src/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type userController struct {
	*logrus.Logger
}

func NewUserController() *userController {
	uc := userController{
		logger.GetLogger(),
	}
	return &uc
}

var userService = service.NewUserService()

func (uc *userController) Register(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		uc.Logger.Error("username is missing")
		vo.ResponseDataFail(c, errors.New("username is missing"))
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok {
		uc.Logger.Error("password is missing")
		vo.ResponseDataFail(c, errors.New("password is missing"))
		return
	}
	_, err := userService.RegisterUser(username, password)
	if err != nil {
		uc.Logger.Error(err)
		vo.ResponseDataFail(c, err)
		return
	}
	vo.ResponseWithoutData(c, 0, 200, "注册成功")
	return
}

func (uc *userController) LoginByUserName(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		uc.Logger.Error("username is missing")
		vo.ResponseDataFail(c, errors.New("username is missing"))
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok {
		uc.Logger.Error("password is missing")
		vo.ResponseDataFail(c, errors.New("password is missing"))
		return
	}
	user, err := userService.LoginUserByUserName(username, password)
	if err != nil {
		uc.Logger.Error("login failed")
		vo.ResponseDataFail(c, err)
		return
	}
	vo.ResponseDataSuccess(user, c)
}

func (uc *userController) LoginBySession(c *gin.Context) {
	sessionKey := c.GetHeader("session")
	user, ok := userService.LoginUserBySession(sessionKey)
	if !ok {
		uc.Logger.Error("login failed")
		vo.ResponseDataFail(c, errors.New("login failed"))
		return
	}
	vo.ResponseDataSuccess(user, c)

}

func (uc *userController) Logout(c *gin.Context) {
	sessionKey := c.GetHeader("session")
	ok := userService.Logout(sessionKey)
	if !ok {
		uc.Logger.Error("logout failed")
		vo.ResponseDataFail(c, errors.New("logout failed"))
		return
	}
	vo.ResponseDataSuccess(nil, c)
}
