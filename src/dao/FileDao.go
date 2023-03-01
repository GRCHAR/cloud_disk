package dao

type File struct {
	Id         string
	Name       string
	Parent     Dir
	CreateUser string
	FileType   string
	FileLength int64
	CreateTime int64
	UpdateTime int64
}

type FileDao struct {
}

func NewFile() *File {
	return new(File)
}

func (*FileDao) CreateFile() {

}

func (*FileDao) FindFile(id string) File {
	return *NewFile()
}
