package dao

type DownloadTask struct {
	id         string
	fileName   string
	fileSize   string
	partNumber int
}

type DownloadTaskDao struct {
}

func (*DownloadTaskDao) CreateDownloadTask() {

}
