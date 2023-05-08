package dao

import "testing"

var fileDao *FileDao
var file File

func init() {
	fileDao = NewFileDao()
}

func TestCreateFile(t *testing.T) {
	file, err := fileDao.CreateFile("ghr-test", 0, 0, "", 0)
	if err != nil {
		t.Error("CreateFile failed:", err)
		t.Fail()
		return
	}
	if file.Id >= 0 {
		t.Log("File created:", file)
	}
}

func TestFindFile(t *testing.T) {
	findFile, err := fileDao.FindFile(file.Id)
	if err != nil {
		t.Error("FindFile failed:", err)
		t.Fail()
		return
	}
	if findFile.Name == file.Name && findFile.ParentId == file.ParentId && findFile.CreateUser == file.CreateUser {
		t.Log("File found:", findFile)
	}
}

func TestFindAllFilesByDirId(t *testing.T) {
	findFiles, err := fileDao.FindAllFilesByDirId(0)
	if err != nil {
		t.Error("FindAllFilesByDirId failed:", err)
		t.Fail()
		return
	}
	for _, value := range findFiles {
		if value.ParentId == 0 {
			t.Log("FindAllFilesByDirId found")
			break
		}
	}
}

func TestDeleteFile(t *testing.T) {
	_, err := fileDao.DeleteFile(file.Id)
	if err != nil {
		t.Error("DeleteFile failed:", err)
		t.Fail()
		return
	}
	t.Log("File deleted")
}
