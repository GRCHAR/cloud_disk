package dao

import "cloud_disk/src/logger"

type User struct {
	Id       int64  `gorm:"primary_key;index:id_idx"`
	Name     string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}

type UserDao struct {
	Dao
}

func NewUserDao() *UserDao {
	return &UserDao{Dao{logger: logger.GetLogger()}}
}

func (*UserDao) CreateUser(user User) (int64, error) {
	if err := db.Create(&user).Error; err != nil {
		return -1, err
	}
	return user.Id, nil
}

func (*UserDao) FindUserById(userId string) (User, error) {
	user := User{}
	if err := db.Where("Id =?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (*UserDao) FindUserByUserName(userName string) (User, error) {
	user := User{}
	if err := db.Where("Name=?", userName).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
