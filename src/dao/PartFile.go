package dao

type PartFile struct {
	Id       string
	TaskId   string
	Number   int
	FileSize string
}

type PartFileDao struct {
}

func (*PartFileDao) CreatePartFile() {

}
