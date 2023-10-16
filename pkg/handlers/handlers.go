package handlers

import (
	"log"
	"reflect"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdiseddik.com/pkg/services/filemanager"
)

func Arborescence(c *websocket.Conn) {
	// keep the connection alive
	var folderContentp, err = filemanager.ReadFolder("./public")
	folderContent := *folderContentp
	c.WriteJSON(folderContent)

	if err != nil {
		log.Println("Error:", err)
	}
	for {
		NewFolderContent := <-filemanager.FolderUpdated
		if reflect.DeepEqual(folderContent, NewFolderContent) {
			continue
		}
		folderContent = NewFolderContent
		err := c.WriteJSON(folderContent)
		if err != nil {
			log.Println("Error:", err)
		}
	}
}

type CreateFileBody struct {
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "folder"
}

func CreateFile(c *fiber.Ctx) error {
	body := new(CreateFileBody)
	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing body",
		})
	}
	if err := c.BodyParser(body); err != nil {
		return err
	}
	switch body.Type {
	case "file":
		err := filemanager.CreateFile(body.Path)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case "folder":
		err := filemanager.CreateFolder(body.Path)
		if err != nil {
			c.Status(500).JSON(fiber.Map{
				"message": "Error while creating folder",
			})
			return err
		}
	default:
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid type",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "File created successfully",
	})
}
