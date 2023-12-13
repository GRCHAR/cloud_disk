package server_tasks

import (
	"cloud_disk/src/config"
	"cloud_disk/src/dao"
	"cloud_disk/src/logger"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
)

type deleteFileTask struct {
	serverTask
}

var fileDao = dao.NewFileDao()
var stm = newServerTaskManager()

func NewDeleteFileTask() *deleteFileTask {
	return &deleteFileTask{serverTask{
		taskId:    "",
		task:      nil,
		stop:      nil,
		loop:      false,
		afterTime: 0,
		logger:    logger.GetLogger(),
	}}
}

func init() {

}

func (*deleteFileTask) deleteUnUseFile() error {
	tempDirPath := config.GetConfig().Server.TempDirPath
	dirs, err := os.ReadDir(tempDirPath)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		isInDataBase, err := InDatabase(dir.Name())
		if err != nil {
			continue
		}
		if dir.IsDir() && !isInDataBase {
			os.RemoveAll(tempDirPath + dir.Name())
		}
	}
	return nil
}

func InDatabase(id string) (bool, error) {
	idInt, _ := strconv.Atoi(id)

	_, err := fileDao.FindFile(int64(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return true, err
		}
	}
	return true, nil
}
