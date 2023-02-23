package dao

type dir struct {
	Id         string
	Name       string
	Parent     *dir
	CreateUser user
	CreateTime int64
	UpdateTime int64
}

func NewDir() *dir {
	return new(dir)
}

func (*dir) Create(name string, parent dir) {

}

func (*dir) UpdateName(name string) {

}

func (*dir) UpdateParent(parent dir) {

}
