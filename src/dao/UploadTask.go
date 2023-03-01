package dao

import "math"

type UploadTask struct {
	Id           string
	FileId       string
	UploadPath   string
	FileName     string
	FileSize     int64
	PartNumber   int
	UserId       string
	UploadedPart []int
	Status       string
}

type UploadTaskDao struct {
}

func (*UploadTaskDao) CreateUploadTask(file File) (string, error) {
	fileName := file.Name
	fileSize := file.FileLength
	partTotalNumber := math.Ceil(float64(fileSize / (2 * 1024)))
	uploadTask := UploadTask{FileName: fileName, PartNumber: int(partTotalNumber)}
	return uploadTask.Id, nil
}

func (*UploadTaskDao) FindUploadTaskById(id string) (UploadTask, error) {
	uploadTask := UploadTask{
		Id:           "",
		FileName:     "",
		FileSize:     0,
		PartNumber:   0,
		UserId:       "",
		UploadedPart: nil,
	}
	return uploadTask, nil
}

func (*UploadTaskDao) UpdateUploadTask(task UploadTask) (UploadTask, error) {
	uploadTask := UploadTask{
		Id:           "",
		FileId:       "",
		UploadPath:   "",
		FileName:     "",
		FileSize:     0,
		PartNumber:   0,
		UserId:       "",
		UploadedPart: nil,
		Status:       "",
	}
	return uploadTask, nil
}
