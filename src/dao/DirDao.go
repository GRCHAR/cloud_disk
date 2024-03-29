package dao

import "cloud_disk/src/logger"

type Dir struct {
	Id         int64
	Name       string
	Parent     *Dir
	ParentId   int64
	UserId     int64
	CreateTime int64
	UpdateTime int64
}

type DirDao struct {
	Dao
}

func NewDirDao() *DirDao {
	dirDao := &DirDao{Dao{
		logger: logger.GetLogger(),
	}}
	return dirDao
}

func (*DirDao) Create(parentId int64, dirName string) (dir Dir, err error) {
	return Dir{}, nil
}

func (*DirDao) UpdateName(dirId int64, name string) (dir Dir, err error) {
	return Dir{}, nil
}

func (*DirDao) UpdateParent(dirId int64, parentId int64) (dir Dir, err error) {
	return Dir{}, nil
}

func (*DirDao) Delete(id int64) error {
	return nil
}

func (*DirDao) FindAllDirByDirId(id int64) (dirs []Dir, err error) {
	return []Dir{}, nil
}

func (*DirDao) FindDirByDirId(id int64) (Dir, error) {

	return Dir{}, nil
}
