package dao

type file struct {
	Id         string
	Name       string
	Parent     dir
	CreateUser user
	FileType   string
	FileLength int64
	CreateTime int64
	UpdateTime int64
}

func NewFile() *file {
	return new(file)
}

func (*file) CreateFile() {

}

func (*file) FindFile(id string) file {
	return *NewFile()
}
