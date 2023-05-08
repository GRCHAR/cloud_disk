package vo

import "cloud_disk/src/dao"

type DirContentVo struct {
	Files []dao.File
	Dirs  []dao.Dir
}
