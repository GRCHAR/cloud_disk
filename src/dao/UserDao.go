package dao

type User struct {
	Id       string
	Name     string
	Password string
}

type UserDao struct {
}

func (*UserDao) CreateUser(user User) {

}

func (*UserDao) FindUserById(userId string) (User, error) {
	return User{}, nil
}
