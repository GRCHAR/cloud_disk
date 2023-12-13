package service

import (
	"cloud_disk/src/dao"
	"cloud_disk/src/dao/redis"
	"cloud_disk/src/logger"
	"errors"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type userService struct {
	rawService
}

func NewUserService() *userService {
	myUserService := userService{rawService{logger: logger.GetLogger()}}
	return &myUserService
}

var userDao = dao.NewUserDao()
var rdb = redis.Rdb

// RegisterUser register user
func (service *userService) RegisterUser(username string, password string) (int64, error) {
	user := dao.User{
		Name:     username,
		Password: password,
	}
	createUser, err := userDao.CreateUser(user)
	if err != nil {
		service.logger.Error("CreateUser error:%v\n", err)
		return createUser, err
	}
	return createUser, nil
}

func (service *userService) LoginUserByUserName(userName string, password string) (dao.User, error) {
	user, err := userDao.FindUserByUserName(userName)
	if err != nil {
		service.logger.Error("LoginUserByUserName FindUserByUserName error:%v\n", err)
		return dao.User{}, err
	}
	if user.Password != password {
		service.logger.Error("LoginUserByUserName password is wrong ")
		return dao.User{}, errors.New(`"password" is wrong`)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		logger.GetLogger().Error("UUID error:%v\n", err)
		return dao.User{}, err
	}

	if !redis.CreateSessionValue(id.String(), "expiryTime", string(time.Now().Add(7*24*time.Hour).Unix())) {
		return dao.User{}, errors.New(`"sessionId" is already used`)
	}
	if !redis.CreateSessionValue(id.String(), "userId", strconv.FormatInt(user.Id, 10)) {
		return dao.User{}, errors.New(`"sessionId" is already used`)
	}

	return user, nil
}

func (service *userService) LoginUserBySession(sessionKey string) (dao.User, bool) {
	valueKey, ok := redis.FindSessionBySessionKey(sessionKey)

	if !ok {
		return dao.User{}, false
	}

	expiryTime, ok := redis.GetSessionValue(valueKey, "expiryTime")
	if !ok {
		logger.GetLogger().Info("user logout")
		return dao.User{}, false
	}

	expiryTimeInt, _ := strconv.ParseInt(expiryTime, 10, 64)

	if time.Now().After(time.Unix(expiryTimeInt, 0)) {
		logger.GetLogger().Info("user logout")
		redis.DeleteSession(valueKey)
		return dao.User{}, false
	}

	userId, ok := redis.GetSessionValue(valueKey, "userId")
	if !ok {
		return dao.User{}, false
	}

	user, err := userDao.FindUserById(userId)
	if err != nil {
		logger.GetLogger().Error("LoginUserBySession FindUserById err", err)
		return dao.User{}, false
	}
	return user, true

}

func (service *userService) Logout(sessionKey string) bool {
	return redis.DeleteSession(sessionKey)
}
