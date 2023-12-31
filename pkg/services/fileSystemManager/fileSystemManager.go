package filesystemmanager

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mehdiseddik.com/pkg/models"
)

var RootFolder = models.Folder{
	Files:   []*models.File{},
	SubDirs: []*models.Folder{},
	Id:      "rootFolder",
}

func FindFolderById(folder *models.Folder, id string) *models.Folder {
	if folder.Id == id {
		return folder
	}
	for _, subFolder := range folder.SubDirs {
		found := FindFolderById(subFolder, id)
		if found != nil {
			return found
		}
	}
	return nil
}

func CreateFile(name string, parentFolderId string) (*models.File, error) {
	foundFolder := FindFolderById(&RootFolder, parentFolderId)
	if foundFolder == nil {
		return nil, errors.New("parent folder not found")
	}
	if foundFolder.CheckDuplicateName(name) {
		return nil, errors.New("folder name already exists")
	}
	newFile := &models.File{
		Content: "",
		Name:    name,
		Id:      string(uuid.New().String()),
	}
	foundFolder.AddFile(newFile)
	return newFile, nil
}

func FindFileById(id string) (*models.File, error) {
	found, error := RootFolder.FindFileById(id)
	if error != nil {
		return nil, error
	}
	return found, nil
}

func UpdateFileName(id string, name string) (*models.File, error) {
	foundFile, err := RootFolder.FindFileById(id)
	if err != nil {
		return nil, err
	}
	parentFolder := RootFolder.GetFileParentFolder(id)
	if foundFile == nil {
		return nil, errors.New("file not found")
	}
	if parentFolder.CheckDuplicateName(name) {
		return nil, errors.New("file name already exists")
	}
	foundFile.SetName(name)
	return foundFile, nil
}

func UpdateFolderName(id string, name string) (*models.Folder, error) {
	foundFolder := RootFolder.FindFolderById(id)
	parentFolder := RootFolder.GetFolderParentFolder(id)
	if foundFolder == nil {
		return nil, errors.New("folder not found")
	}
	if parentFolder.CheckDuplicateName(name) {
		return nil, errors.New("folder name already exists")
	}
	foundFolder.SetName(name)
	return foundFolder, nil
}

func CreateFolder(name string, parentFolderId string) (*models.Folder, error) {
	foundFolder := FindFolderById(&RootFolder, parentFolderId)
	if foundFolder == nil {
		return nil, errors.New("parent folder not found")
	}
	if foundFolder.CheckDuplicateName(name) {
		return nil, errors.New("folder name already exists")
	}
	newFolder := &models.Folder{
		Files:   []*models.File{},
		SubDirs: []*models.Folder{},
		Name:    name,
		Id:      string(uuid.New().String()),
	}
	foundFolder.AddSubDir(newFolder)

	return newFolder, nil
}

func MoveFile(fileId string, folderId string) (*models.File, error) {
	foundFile, err := RootFolder.FindFileById(fileId)
	if err != nil {
		return nil, err
	}
	if foundFile == nil {
		return nil, errors.New("file not found")
	}
	foundFolder := FindFolderById(&RootFolder, folderId)
	if foundFolder == nil {
		return nil, errors.New("folder not found")
	}
	oldParentFolder := RootFolder.GetFileParentFolder(fileId)

	foundFolder.AddFile(foundFile)
	err = oldParentFolder.RemoveFile(foundFile.Id) // probleme
	if err != nil {
		return nil, err
	}
	return foundFile, nil
}
