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
	arborescenceBroadCast(filesystemmanager.RootFolder.ToFolderWithoutFileContent())
	c.Status(201)
	return c.JSON(returned)
}

type EditFolderNameRequest struct {
	Name string `json:"name"`
}

func EditFolderNameHandler(c *fiber.Ctx) error {
	folderId := c.Params("folderId")
	var body EditFolderNameRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	folderName := body.Name
	edited, err := filesystemmanager.UpdateFolderName(folderId, folderName)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	arborescenceBroadCast(filesystemmanager.RootFolder.ToFolderWithoutFileContent())
	returned := fiber.Map{
		"editedFile": edited,
		"msg":        "file edited successfully",
	}
	c.Status(200)
	return c.JSON(returned)
}
