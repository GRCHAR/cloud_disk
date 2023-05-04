package dao

import (
	"cloud_disk/src/logger"
)

type File struct {
	Id         int64  `gorm:"primary_key;index:id_idx"`
	Name       string `gorm:"type:varchar(255);not null;index:name_index"`
	Parent     Dir
	ParentId   int64  `gorm:"type:bigint;not null"`
	CreateUser int64  `gorm:"type:bigint;not null"`
	FileType   string `gorm:"type:varchar(255);not null"`
	FileLength int64  `gorm:"type:bigint"`
	CreateTime int64  `gorm:"type:bigint"`
	UpdateTime int64  `gorm:"type:bigint"`
}

type FileDao struct {
	Dao
}

func NewFileDao() *FileDao {
	fileDao := new(FileDao)
	fileDao.logger = logger.GetLogger()
	return fileDao
}

func (fileDao *FileDao) CreateFile(name string, parentId int, userId int, fileType string, fileLength int64) (File, error) {
	file := File{
		Name: name,
		//Parent:     Dir{},
		ParentId:   int64(parentId),
		CreateUser: int64(userId),
		FileType:   fileType,
		FileLength: int64(fileLength),
		CreateTime: 0,
		UpdateTime: 0,
	}
	err := db.Create(&file).Error
	if err != nil {
		fileDao.logger.Error()
		return File{}, err
	}
	return file, nil
}

func (fileDao *FileDao) FindFile(id int) (File, error) {
	file := File{}
	if err := db.Model(&file).Where("id=?", int64(id)).Error; err != nil {
		fileDao.logger.Error(err)
		return file, err
	}
	return file, nil
}

func (fileDao *FileDao) DeleteFile(id int64) (File, error) {
	file := File{}
	if err := db.Delete(&file, id).Error; err != nil {
		fileDao.logger.Error(err)
		return file, err
	}
	return file, nil
}

func (*FileDao) FindAllFilesByUserId(id string) []File {
	return []File{}
}

func (fileDao *FileDao) FindAllFilesByDirId(id string) ([]File, error) {
	var files []File
	if err := db.Model(&files).Where("ParentId", id).Error; err != nil {
		fileDao.logger.Error(err)
		return files, nil
	}

	return files, nil
}
