package config

import "cloud_disk/src/dao"

type configDB struct {
	DownloadTaskThreadNumber int
	UploadTaskThreadNumber   int
	DownloadTaskNumber       int
	UploadTaskNumber         int
}

var db = dao.GetDB()

var configInfo configDB

func init() {
	configInfo = getConfigInfo()
}

func getConfigDB() *configDB {
	return new(configDB)
}

func SaveConfig(config configDB) {
}

func getConfigInfo() configDB {
	configDB := configDB{
		DownloadTaskThreadNumber: 4,
		UploadTaskThreadNumber:   4,
	}
	return configDB
}

func GetConfigInfo() *configDB {
	return &configInfo
}
