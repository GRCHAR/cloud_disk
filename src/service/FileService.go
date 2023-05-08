package service

import (
	"cloud_disk/src/config"
	"cloud_disk/src/dao"
	"cloud_disk/src/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	UploadTaskFileChanMap map[int]chan uploadPieceFile
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

func init() {
	uploadTaskDao = dao.GetUploadDao()
	fileDao = dao.NewFileDao()
}

func GetFileService() *FileService {
	service := new(FileService)
	service.logger = logger.GetLogger()
	service.UploadTaskFileChanMap = make(map[int]chan uploadPieceFile, 10000)
	service.MergeTaskChanMap = *new(sync.Map)
	return service
}

func (service *FileService) CreateUploadTask(fileName string, fileSize int64, userId string) (taskId int, err error) {

	//uploadUser, err := userDao.FindUserById(userId)
	//if err != nil {
	//	service.logger.Error("CreateUploadTask FindUserById error:")
	//	return -1, err
	//
	userIdInt, _ := strconv.Atoi(userId)
	uploadFile := dao.File{FileLength: fileSize, Name: fileName, CreateUser: int64(userIdInt)}
	taskId, err = uploadTaskDao.CreateUploadTask(uploadFile)
	if err != nil {
		service.logger.Error("CreateUploadTask CreateUploadTask error:")
		return -1, err
	}
	if !service.startUploadTask(taskId) {
		return -1, errors.New("start upload task failed")
	}
	return taskId, nil
}

func (service *FileService) startUploadTask(taskId int) bool {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("FindUploadTaskById err")
		return false
	}
	task.Status = "uploading"
	task, err = uploadTaskDao.UpdateUploadTask(task)
	if err != nil {
		service.logger.Error("UpdateUploadTask err")
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

func (service *FileService) saveFile(c *gin.Context, header *multipart.FileHeader, taskId int, partNumber int) error {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		return err
	}
	uploadTempDir := config.GetConfig().Server.TempDirPath

	taskIdString := strconv.Itoa(taskId)
	userIdString := strconv.Itoa(task.UserId)
	os.MkdirAll(uploadTempDir+"/"+userIdString+"/"+taskIdString+"/", 0777)
	err = c.SaveUploadedFile(header, uploadTempDir+"/"+userIdString+"/"+taskIdString+"/"+strconv.Itoa(partNumber))
	if err != nil {
		service.logger.Error("SaveUploadedFile err", err)
		return err
	}
	return nil
}

func (service *FileService) runTask(taskId int) error {
	mergeChan := make(chan uploadPieceFile, 100)
	service.MergeTaskChanMap.Store(taskId, &mergeChan)
	return nil
}

func (service *FileService) AppendUploadFileChannel(c *gin.Context, file *multipart.FileHeader, taskId int, pieceNumber int) {
	service.logger.Info(service.UploadTaskFileChanMap)
	service.UploadTaskFileChanMap[taskId] <- uploadPieceFile{
		context:     c,
		header:      file,
		pieceNumber: pieceNumber,
	}

}

func (service *FileService) RunUploadTask(taskId int, userId int) error {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if task.UserId != userId {
		service.logger.WithFields(logrus.Fields{
			"userId": userId,
			"taskId": taskId,
		})
		return errors.New("不匹配的用户ID")
	}
	if err != nil {
		service.logger.Error("runUploadTask FindUploadTaskById err")
	}
	fileChan := make(chan uploadPieceFile, 100)
	mergeChan := make(chan mergeTaskFile, 100)
	service.UploadTaskFileChanMap[taskId] = fileChan
	service.MergeTaskChanMap.Store(taskId, mergeChan)
	count := 0
	go func() {
		defer func() {
			if r := recover(); r != nil {
				service.logger.Error(r)
				return
			}
		}()

		for task.PartNumber > count {
			pieceFile := <-service.UploadTaskFileChanMap[taskId]
			service.saveFile(pieceFile.context, pieceFile.header, taskId, pieceFile.pieceNumber)
			channel, ok := service.MergeTaskChanMap.Load(taskId)
			mergeChan := channel.(chan mergeTaskFile)
			if !ok {
				service.logger.Error("can't find mergeChan:", taskId)
				task.Status = "fail"
				uploadTaskDao.UpdateUploadTask(task)
				break
				//return errors.New("can't find mergeChan")
			}
			uploadTempDir := config.GetConfig().Server.TempDirPath

			mergeChan <- mergeTaskFile{
				dist:             uploadTempDir + "/" + strconv.Itoa(userId) + "/" + strconv.Itoa(taskId) + "/" + strconv.Itoa(pieceFile.pieceNumber),
				firstPieceNumber: pieceFile.pieceNumber,
				lastPieceNumber:  pieceFile.pieceNumber,
			}
			count++
		}
	}()

	task.Status = "run"
	uploadTaskDao.UpdateUploadTask(task)
	return nil
}

func (service *FileService) RunMergeTask(taskId int) {
	channel, _ := service.MergeTaskChanMap.Load(taskId)
	mergeChan := channel.(chan mergeTaskFile)
	var mergeArray []mergeTaskFile
	task, _ := uploadTaskDao.FindUploadTaskById(taskId)
	count := 0
	path := config.GetConfig().Server.MergeDirPath + "/" + strconv.Itoa(task.UserId) + "/" + strconv.Itoa(taskId)
	if err := os.MkdirAll(path, 0777); err != nil {
		service.logger.Error(err)
		task.Status = "fail"
		uploadTaskDao.UpdateUploadTask(task)
		close(mergeChan)
		return
	}
	_, err := os.Create(path + "/" + "result")
	if err != nil {
		service.logger.Error(err)
		task.Status = "fail"
		uploadTaskDao.UpdateUploadTask(task)
		close(mergeChan)
		return
	}

	resultFile, err := fileDao.CreateFile(task.FileName, task.ParentId, task.UserId, "", task.FileSize)
	if err != nil {
		task.Status = "fail"
		uploadTaskDao.UpdateUploadTask(task)
		close(mergeChan)
		return
	}

	go func() {

		for {
			//if count == task.PartNumber {
			//	close(mergeChan)
			//	break
			//}
			if count == 6 {
				close(mergeChan)
				break
			}
			mergeFile := <-mergeChan

			count++
			notAppend := false

			if len(mergeArray) == 0 {
				mergeArray = append(mergeArray, mergeFile)
				notAppend = true
			}
			for key, mergedFile := range mergeArray {

				if mergeFile.firstPieceNumber == mergedFile.firstPieceNumber-1 {

					//file, err := os.OpenFile(path+"/"+"result", os.O_RDWR, 0777)

					file, err := os.OpenFile(mergeFile.dist, os.O_RDWR, 0755)

					firstFile, _ := os.ReadFile(mergeFile.dist)
					if err != nil {
						service.logger.Error(err)
						fileDao.DeleteFile(resultFile.Id)
						return
					}
					appendFile, _ := os.ReadFile(mergedFile.dist)
					firstFile = append(firstFile, appendFile...)
					file.Write(firstFile)
					file.Close()
					mergeArray[key].dist = mergeFile.dist
					mergeArray[key].firstPieceNumber -= 1
					break

				} else if mergeFile.lastPieceNumber == mergedFile.lastPieceNumber+1 {

					//file, err := os.OpenFile(path+"/"+"result", os.O_RDWR, 0777)
					file, err := os.OpenFile(mergedFile.dist, os.O_RDWR, 0755)
					firstFile, _ := os.ReadFile(mergedFile.dist)
					if err != nil {
						service.logger.Error(err)
						fileDao.DeleteFile(resultFile.Id)
						return
					}
					appendFile, _ := os.ReadFile(mergeFile.dist)
					firstFile = append(firstFile, appendFile...)
					//info, _ := file.Stat()
					file.Write(firstFile)
					file.Close()
					mergeArray[key].lastPieceNumber += 1
					break
				} else {
					if !notAppend {
						mergeArray = append(mergeArray, mergeFile)
						break
					}
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
		file, err := os.OpenFile(path+"/"+"result", os.O_RDWR, 0777)
		if err != nil {
			service.logger.Error(err)
			task.Status = "fail"
			uploadTaskDao.UpdateUploadTask(task)
			fileDao.DeleteFile(resultFile.Id)
			return
		}
		info, err = firstFile.Stat()
		result := make([]byte, info.Size())
		_, err = firstFile.Read(result)
		if err != nil {
			service.logger.Error(err)
			task.Status = "fail"
			uploadTaskDao.UpdateUploadTask(task)
			return
		}
		_, err = file.Write(result)
		if err != nil {
			service.logger.Error(err)
			task.Status = "fail"
			uploadTaskDao.UpdateUploadTask(task)
			return
		}
		file.Close()
		task.Status = "success"
		uploadTaskDao.UpdateUploadTask(task)

	}()

}

func (service *FileService) MergeFile(taskId int) error {
	go func() {
		uploadTempDir := config.GetConfig().Server.TempDirPath
		resultDir := config.GetConfig().Server.MergeDirPath
		task, err := uploadTaskDao.FindUploadTaskById(taskId)
		if err != nil {
			service.logger.Error("FindUploadTaskById err")
			return
		}
		taskIdString := strconv.Itoa(taskId)
		userIdString := strconv.Itoa(task.UserId)
		fs := os.DirFS(uploadTempDir + "/" + userIdString + taskIdString)
		resultFile, err := os.Create(resultDir + "/" + userIdString + "/" + task.FileId)
		if err != nil {
			service.logger.Error("os.Create err")

			return
		}
		var tempByte []byte
		var readAtSize int64
		for i := 0; i < task.PartNumber; i++ {
			tempByte = []byte{}
			file, err := fs.Open(strconv.Itoa(i))
			if err != nil {
				service.logger.Error("fs.Open err")
				return
			}
			stat, err := file.Stat()
			if err != nil {
				service.logger.Error("fs.Stat err")
				return
			}
			_, err = file.Read(tempByte)
			if err != nil {
				service.logger.Error("file.Read err")
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

func (service *FileService) GetTaskStatus(taskId int) (status string, err error) {
	task, err := uploadTaskDao.FindUploadTaskById(taskId)
	if err != nil {
		service.logger.Error("uploadTaskDao.FindUploadTaskById err,")
		return "", err
	}
	return task.Status, nil
}

func (service *FileService) GetDirFiles(dirId string, userId string) (files []dao.File, err error) {
	dirIdInt, _ := strconv.Atoi(dirId)
	files, err = fileDao.FindAllFilesByDirId(int64(dirIdInt))
	return
}

func (service *FileService) CreateDownloadTask(fileId string, userId string) {

}
