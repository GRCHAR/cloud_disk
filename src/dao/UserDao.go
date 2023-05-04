package dao

type User struct {
	Id       string `gorm:"primary_key;index:id_idx"`
	Name     string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}

type UserDao struct {
}

func (*UserDao) CreateUser(user User) {

}

func (*UserDao) FindUserById(userId string) (User, error) {
	return User{}, nil
}
