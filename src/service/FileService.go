package service

import (
	"cloud_disk/src/config"
	"cloud_disk/src/dao"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"sync"
)

type uploadPieceFile struct {
	context     *gin.Context
	header      *multipart.FileHeader
	pieceNumber int
}

type FileService struct {
	rawService
	MergeTaskChanMap      sync.Map
	UploadTaskFileChanMap map[string]chan uploadPieceFile
}

type mergeTaskFile struct {
	dist             string
	firstPieceNumber int
	lastPieceNumber  int
}

type uploadTask struct {
}

type uploadTaskPool struct {
}

type taskSystem struct {
	mergeChannel chan int
}

var uploadTaskDao *dao.UploadTaskDao
var userDao *dao.UserDao
var fileDao *dao.FileDao

func (service *FileService) CreateUploadTask(fileName string, fileSize int64, userId string) (taskId string, err error) {

	uploadUser, err := userDao.FindUserById(userId)
	if err != nil {
		service.logger.Error("CreateUploadTask FindUserById error:", zap.Error(err))
		return "", err
	}
	uploadFile := dao.File{FileLength: fileSize, Name: fileName, CreateUser: uploadUser.Id}
	taskId, err = uploadTaskDao.CreateUploadTask(uploadFile)
	if err != nil {
		service.logger.Error("CreateUploadTask CreateUploadTask error:", zap.Error(err))
		return "", err
	}
	if !service.startUploadTask(taskId) {
		return "", errors.New("start upload task failed")
	}
	return taskId, nil
}

func (service *FileService) startUploadTask(taskId string) bool {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("FindUploadTaskById err", zap.Error(err))
		return false
	}
	task.Status = "uploading"
	task, err = uploadTaskDao.UpdateUploadTask(task)
	if err != nil {
		service.logger.Error("UpdateUploadTask err", zap.Error(err))
		return false
	}
	go service.runTask(task.Id)
	return true
}

func (service *FileService) StopUploadTask() {

}

func (service *FileService) ContinueUploadTask() {

}

func (service *FileService) DeleteUploadTask() {

}

func (service *FileService) GetTaskInfo(userId string) {

}

func (service *FileService) saveFile(c *gin.Context, header *multipart.FileHeader, taskId string, partNumber int) error {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		return err
	}
	uploadTempDir := config.GetConfig().Server.TempDirPath

	err = c.SaveUploadedFile(header, uploadTempDir+"/"+task.UserId+"/"+task.Id+"/"+strconv.Itoa(partNumber))

	if err != nil {
		service.logger.Error("SaveUploadedFile err", zap.Error(err))
		return err
	}
	return nil
}

func (service *FileService) runTask(taskId string) error {
	mergeChan := make(chan uploadPieceFile, 100)
	service.MergeTaskChanMap.Store(taskId, &mergeChan)
	return nil
}

func (service *FileService) AppendUploadFileChannel(c *gin.Context, file *multipart.FileHeader, taskId string, pieceNumber int) {
	service.UploadTaskFileChanMap[taskId] <- uploadPieceFile{
		context:     c,
		header:      file,
		pieceNumber: pieceNumber,
	}

}

func (service *FileService) runUploadTask(taskId string) {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("runUploadTask FindUploadTaskById err", zap.Error(err))
	}
	fileChan := make(chan uploadPieceFile, 100)
	service.UploadTaskFileChanMap[taskId] = fileChan
	count := 0
	for task.PartNumber > count {
		pieceFile := <-service.UploadTaskFileChanMap[taskId]
		service.saveFile(pieceFile.context, pieceFile.header, taskId, pieceFile.pieceNumber)
		channel, ok := service.MergeTaskChanMap.Load(taskId)
		mergeChan := channel.(*chan mergeTaskFile)
		if !ok {
			service.logger.Error("can't find mergeChan:" + taskId)
			return
		}
		*mergeChan <- mergeTaskFile{
			dist:             "dist",
			firstPieceNumber: pieceFile.pieceNumber,
			lastPieceNumber:  pieceFile.pieceNumber,
		}
		count++
	}
	task.Status = "success"
	uploadTaskDao.UpdateUploadTask(task)
}

func (service *FileService) runMergeTask(taskId string) {
	channel, _ := service.MergeTaskChanMap.Load(taskId)
	mergeChan := channel.(*chan mergeTaskFile)
	var mergeArray []mergeTaskFile
	task, _ := uploadTaskDao.FindUploadTaskById(taskId)
	count := 0
	for {
		if count == task.PartNumber {
			close(*mergeChan)
			break
		}
		mergeFile := <-*mergeChan
		count++
		for key, mergedFile := range mergeArray {

			if mergeFile.firstPieceNumber == mergedFile.firstPieceNumber-1 {
				file, err := os.OpenFile(mergeFile.dist, os.O_RDWR, 0755)
				firstFile, _ := os.ReadFile(mergeFile.dist)
				if err != nil {
					return
				}
				appendFile, _ := os.ReadFile(mergedFile.dist)
				firstFile = append(firstFile, appendFile...)
				file.Write(firstFile)
				mergeArray[key].firstPieceNumber -= 1

			} else if mergeFile.lastPieceNumber == mergedFile.lastPieceNumber+1 {
				file, err := os.OpenFile(mergedFile.dist, os.O_RDWR, 0755)
				firstFile, _ := os.ReadFile(mergedFile.dist)
				if err != nil {
					return
				}
				appendFile, _ := os.ReadFile(mergeFile.dist)
				firstFile = append(firstFile, appendFile...)
				info, _ := file.Stat()
				file.WriteAt(firstFile, info.Size())
				mergeArray[key].lastPieceNumber += 1
			} else {
				mergeArray = append(mergeArray, mergeFile)
			}

		}

	}

	sort.Slice(mergeArray, func(i, j int) bool {
		return mergeArray[i].firstPieceNumber < mergeArray[j].firstPieceNumber
	})
	firstFile, _ := os.OpenFile(mergeArray[0].dist, os.O_RDWR, 0775)
	info, _ := firstFile.Stat()
	size := info.Size()
	for i := 1; i < len(mergeArray); i++ {
		appendFile, _ := os.ReadFile(mergeArray[i].dist)

		firstFile.WriteAt(appendFile, size)
		size += int64(len(appendFile))
	}
	task.Status = "success"
	uploadTaskDao.UpdateUploadTask(task)
}

func (service *FileService) MergeFile(taskId string) error {
	go func() {
		uploadTempDir := config.GetConfig().Server.TempDirPath
		resultDir := config.GetConfig().Server.MergeDirPath
		task, err := uploadTaskDao.FindUploadTaskById(taskId)
		if err != nil {
			service.logger.Error("FindUploadTaskById err", zap.Error(err))
			return
		}
		fs := os.DirFS(uploadTempDir + "/" + task.UserId + task.Id)
		resultFile, err := os.Create(resultDir + "/" + task.UserId + "/" + task.FileId)
		if err != nil {
			service.logger.Error("os.Create err", zap.Error(err))

			return
		}
		var tempByte []byte
		var readAtSize int64
		for i := 0; i < task.PartNumber; i++ {
			tempByte = []byte{}
			file, err := fs.Open(strconv.Itoa(i))
			if err != nil {
				service.logger.Error("fs.Open err", zap.Error(err))
				return
			}
			stat, err := file.Stat()
			if err != nil {
				service.logger.Error("fs.Stat err", zap.Error(err))
				return
			}
			_, err = file.Read(tempByte)
			if err != nil {
				service.logger.Error("file.Read err", zap.Error(err))
				return
			}

			resultFile.WriteAt(tempByte, readAtSize)
			readAtSize += stat.Size()
		}
		return
	}()
	return nil
}

func (service *FileService) startMergeTask(taskId string) {

}

func (service *FileService) GetTaskStatus(taskId string) (status string, err error) {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("uploadTaskDao.FindUploadTaskById err,", zap.Error(err))
		return "", err
	}
	return task.Status, nil
}

func (service *FileService) GetDirFiles(dirId string, userId string) (files []dao.File, err error) {
	files, err = fileDao.FindAllFilesByDirId(dirId)
	return
}

func (service *FileService) CreateDownloadTask(fileId string, userId string) {

}
