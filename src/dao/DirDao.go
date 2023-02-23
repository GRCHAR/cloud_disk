package dao

type Dir struct {
	Id         string
	Name       string
	Parent     *Dir
	CreateUser user
	CreateTime int64
	UpdateTime int64
}

func NewDir() *Dir {
	return new(Dir)
}

func (*Dir) Create(name string, parent Dir) {

}

func (*Dir) UpdateName(name string) {

}

func (*Dir) UpdateParent(parent Dir) {

}
