package models

import "github.com/kr/pretty"

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

func (f *Folder) AddSubDir(folder *Folder) *Folder {
	f.SubDirs = append(f.SubDirs, folder)
	return f
}

func (f *Folder) FindFileById(id string) *File {

	for _, file := range f.Files {
		pretty.Println(file.Id, id)
		pretty.Println(file.Id == id)
		if file.Id == id {
			return file
		}
	}
	for _, folder := range f.SubDirs {
		found := folder.FindFileById(id)
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
func (f *Folder) RemoveFile(id string) *Folder {
	for i, file := range f.Files {
		if file.Id == id {
			f.Files = append(f.Files[:i], f.Files[i+1:]...)
			return f
		}
	}
	return f
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

// Return the parent folder of the file with the given id
func (f *Folder) GetParentFolder(fileId string) *Folder {
	file := f.FindFileById(fileId)
	if file != nil {
		return f
	}
	for _, folder := range f.SubDirs {
		found := folder.GetParentFolder(fileId)
		if found != nil {
			return found
		}
	}
	return nil
}
