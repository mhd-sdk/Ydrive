package filemanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var FolderUpdated = make(chan FolderContent)

func WatchDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watcher.Add("public")
	err = filepath.Walk("public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {

		case event, ok := <-watcher.Events:
			if ok {
				switch event.Op {
				case fsnotify.Chmod:
					log.Println("Changement de permission du fichier:", event.Name)
				case fsnotify.Remove:
					log.Println("Suppression du fichier:", event.Name)
				case fsnotify.Write:
					log.Println("Modification du fichier:", event.Name)
				case fsnotify.Create:
					// if is folder
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watcher.Add(event.Name)
						log.Println("Création du dossier:", event.Name)
						break
					}
					log.Println("Création du fichier:", event.Name)
				case fsnotify.Rename:
					log.Println("Renommage du fichier:", event.Name)
				}
				folderContent, _ := ReadFolder("./public")
				select {
				case FolderUpdated <- *folderContent:
				default:
				}
			}
		case err, ok := <-watcher.Errors:
			if ok {
				log.Println("error:", err)
			}
		}

	}
}

func CreateFile(filePath string) error {
	err := CreateFolder(filepath.Dir(filePath))
	if err != nil {
		return err
	}
	_, err = os.Create("public/" + filePath)
	if err != nil {
		return err
	}
	return nil
}

func CreateFolder(folderPath string) error {
	fmt.Println("public/" + folderPath)
	err := os.MkdirAll("public/"+folderPath, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type FolderContent struct {
	AbsPath string           `json:"absPath"`
	Files   []string         `json:"files"`
	SubDirs []*FolderContent `json:"subDirs"`
}

func ReadFolder(path string) (*FolderContent, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	var files []string
	var subDirs []*FolderContent

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(absPath, entry.Name())
		if entry.IsDir() {
			subDirContent, err := ReadFolder(entryPath)
			if err != nil {
				return nil, err
			}
			subDirs = append(subDirs, subDirContent)
		} else {
			files = append(files, entryPath)
		}
	}

	return &FolderContent{
		AbsPath: absPath,
		Files:   files,
		SubDirs: subDirs,
	}, nil
}
