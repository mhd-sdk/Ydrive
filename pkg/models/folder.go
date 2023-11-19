package models

import "errors"

type Folder struct {
	Name    string    `json:"name"`
	Id      string    `json:"id"`
	Files   []*File   `json:"files"`
	SubDirs []*Folder `json:"subDirs"`
}

func (f *Folder) AddFile(file *File) *Folder {
	f.Files = append(f.Files, file)
	return f
}

func (f *Folder) SetName(name string) *Folder {
	f.Name = name
	return f
}

func (f *Folder) AddSubDir(folder *Folder) *Folder {
	f.SubDirs = append(f.SubDirs, folder)
	return f
}

func (f *Folder) FindFileById(id string) (*File, error) {

	for _, file := range f.Files {
		if file.Id == id {
			return file, nil
		}
	}
	for _, folder := range f.SubDirs {
		found, err := folder.FindFileById(id)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}
	return nil, errors.New("file not found")
}

func (f *Folder) FindFolderById(id string) *Folder {
	if f.Id == id {
		return f
	}
	for _, folder := range f.SubDirs {
		found := folder.FindFolderById(id)
		if found != nil {
			return found
		}
	}
	return nil
}

// return the subfolder with the given id
func (f *Folder) GetSubDir(id string) *Folder {
	for _, folder := range f.SubDirs {
		if folder.Id == id {
			return folder
		}
	}
	return nil
}

// remove the file with the given id from the files of the current folder
func (f *Folder) RemoveFile(id string) error {
	for i, file := range f.Files {
		if file.Id == id {
			f.Files = append(f.Files[:i], f.Files[i+1:]...)
			return nil
		}
	}
	return errors.New("file not found&&&")
}

// remove the subfolder with the given id from the subfolders of the current folder
func (f *Folder) RemoveSubDir(id string) *Folder {
	for i, folder := range f.SubDirs {
		if folder.Id == id {
			f.SubDirs = append(f.SubDirs[:i], f.SubDirs[i+1:]...)
			return f
		}
	}
	return f
}

// return true if the name already exist in subfolders and files
func (f *Folder) CheckDuplicateName(name string) bool {
	for _, folder := range f.SubDirs {
		if folder.Name == name {
			return true
		}
	}
	for _, file := range f.Files {
		if file.Name == name {
			return true
		}
	}
	return false
}

// Return the folder that contains the file with the given id
func (f *Folder) GetFileParentFolder(fileId string) *Folder {

	for _, file := range f.Files {
		if file.Id == fileId {
			return f
		}
	}
	for _, folder := range f.SubDirs {
		found := folder.GetFileParentFolder(fileId)
		if found != nil {
			return found
		}
	}
	return nil

}

// Return the parent folder of the folder with the given id
func (f *Folder) GetFolderParentFolder(folderId string) *Folder {
	folder := f.FindFolderById(folderId)
	if folder != nil {
		return f
	}
	for _, folder := range f.SubDirs {
		found := folder.GetFolderParentFolder(folderId)
		if found != nil {
			return found
		}
	}
	return nil
}
