package handlers

import (
	"github.com/gofiber/fiber/v2"

	filesystemmanager "github.com/mehdiseddik.com/pkg/services/fileSystemManager"
)

type CreateFolderRequest struct {
	Name           string `json:"name"`
	ParentFolderId string `json:"parentFolderId"`
}

// fiber handler
func CreateFolderHandler(c *fiber.Ctx) error {
	var body CreateFolderRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	name := body.Name
	parentFolderId := body.ParentFolderId
	created, err := filesystemmanager.CreateFolder(name, parentFolderId)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	returned := map[string]interface{}{
		"createdFolder": created,
		"msg":           "folder created successfully",
	}
	Broadcast(filesystemmanager.CurrentFolder)
	c.Status(201)
	return c.JSON(returned)
}
