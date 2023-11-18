package filesystemmanager

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mehdiseddik.com/pkg/models"
)

var CurrentFolder = models.Folder{
	Files:   []*models.File{},
	SubDirs: []*models.Folder{},
	Id:      string(uuid.New().String()),
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
	foundFolder := FindFolderById(&CurrentFolder, parentFolderId)
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

func UpdateFile(id string, name string) (*models.File, error) {
	foundFile := CurrentFolder.FindFileById(id)
	parentFolder := CurrentFolder.GetParentFolder(id)
	if foundFile == nil {
		return nil, errors.New("file not found")
	}
	if parentFolder.CheckDuplicateName(name) {
		return nil, errors.New("file name already exists")
	}
	foundFile.SetName(name)
	return foundFile, nil
}

func CreateFolder(name string, parentFolderId string) (*models.Folder, error) {
	foundFolder := FindFolderById(&CurrentFolder, parentFolderId)
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
