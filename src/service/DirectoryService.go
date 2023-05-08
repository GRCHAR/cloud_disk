package service

import (
	"cloud_disk/src/dao"
	"cloud_disk/src/logger"
	"cloud_disk/src/vo"
)

type DirectoryService struct {
	rawService
}

var dirDao *dao.DirDao

func init() {
	dirDao = dao.NewDirDao()
}

func NewDirectoryService() *DirectoryService {
	return &DirectoryService{rawService{
		logger: logger.GetLogger(),
	}}
}

func (service *DirectoryService) CreateDir(parentId int64, dirName string) error {
	_, err := dirDao.Create(parentId, dirName)
	if err != nil {
		service.logger.Error("Create dir error:%+v", err)
		return err
	}
	return nil
}

func (service *DirectoryService) GetDirContent(parentId int64) (vo.DirContentVo, error) {
	dirs, err := service.ListDir(parentId)
	if err != nil {
		service.logger.Error("List dir error:%+v", err)
		return vo.DirContentVo{}, err
	}
	files, err := service.ListFiles(parentId)
	if err != nil {
		service.logger.Error("List files error:%+v", err)
		return vo.DirContentVo{}, err
	}

	return vo.DirContentVo{Dirs: dirs, Files: files}, nil

}

func (service *DirectoryService) DeleteDir(id int64) error {
	err := fileDao.DeleteAllFilesByDirId(id)
	if err != nil {
		service.logger.Error("Delete dir files error:", err)
	}
	err = dirDao.Delete(id)
	if err != nil {
		service.logger.Error("Delete dir error:", err)
		return err
	}
	return nil
}

func (service *DirectoryService) ListDir(parentId int64) ([]dao.Dir, error) {

	dirs, err := dirDao.FindAllDirByDirId(parentId)
	if err != nil {
		service.logger.Error("ListDir FindAllDirByDirId error:", err)
		return nil, err
	}
	return dirs, nil
}

func (service *DirectoryService) ListFiles(parentId int64) ([]dao.File, error) {
	files, err := fileDao.FindAllFilesByDirId(parentId)
	if err != nil {
		service.logger.Error("ListDir FindAllFilesByDirId error:", err)
		return []dao.File{}, err
	}
	return files, nil
}
