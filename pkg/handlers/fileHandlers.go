package handlers

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	filesystemmanager "github.com/mehdiseddik.com/pkg/services/fileSystemManager"
)

// fiber websocket mamager for writeFile handler
var WriteFileClients = make(map[*websocket.Conn]bool)

// RegisterWriteFile a new client
func RegisterWriteFile(c *websocket.Conn) {
	WriteFileClients[c] = true
}

// UnregisterWriteFile a client
func UnregisterWriteFile(c *websocket.Conn) {
	c.Close()
	delete(WriteFileClients, c)
}

// WriteFileBroadCast a message to all clients
func WriteFileBroadCast(message interface{}) {
	for client := range ArborescenceClients {
		client.WriteJSON(message)
	}
}

type CreateFileRequest struct {
	Name           string `json:"name"`
	ParentFolderId string `json:"parentFolderId"`
}

// fiber handler
func CreateFileHandler(c *fiber.Ctx) error {
	var body CreateFileRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	name := body.Name
	parentFolderId := body.ParentFolderId
	created, err := filesystemmanager.CreateFile(name, parentFolderId)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	arborescenceBroadCast(filesystemmanager.RootFolder)
	returned := fiber.Map{
		"createdFile": created,
		"msg":         "file created successfully",
	}
	c.Status(201)
	return c.JSON(returned)
}

type EditFileRequest struct {
	Name string `json:"name"`
}

func TestHandler(c *fiber.Ctx) error {
	var param string = c.Params("id")
	fmt.Println("param: " + param)
	return c.JSON(fiber.Map{"msg": param})
}

func EditFileNameHandler(c *fiber.Ctx) error {
	fileId := c.Params("fileId")
	var body EditFileRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	fileName := body.Name
	edited, err := filesystemmanager.UpdateFileName(fileId, fileName)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	arborescenceBroadCast(filesystemmanager.RootFolder)
	returned := fiber.Map{
		"editedFile": edited,
		"msg":        "file edited successfully",
	}
	c.Status(200)
	return c.JSON(returned)
}

type MoveFileRequest struct {
	FileId   string `json:"fileId"`
	FolderId string `json:"folderId"`
}

func MoveFileHandler(c *fiber.Ctx) error {
	var body MoveFileRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	fileId := body.FileId
	folderId := body.FolderId
	moved, err := filesystemmanager.MoveFile(fileId, folderId)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	arborescenceBroadCast(filesystemmanager.RootFolder)
	returned := fiber.Map{
		"movedFile": moved,
		"msg":       "file moved successfully",
	}
	c.Status(200)
	return c.JSON(returned)
}

func DeleteFileHandler(c *fiber.Ctx) error {
	fileId := c.Params("fileId")
	err := filesystemmanager.RootFolder.RemoveFile(fileId)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	arborescenceBroadCast(filesystemmanager.RootFolder)
	returned := fiber.Map{
		"msg": "file deleted successfully",
	}
	c.Status(200)
	return c.JSON(returned)
}

func WriteFileContentWsHandler(c *websocket.Conn) {
	defer UnregisterWriteFile(c)
	RegisterWriteFile(c)
	fileId := c.Params("fileId")
	file, err := filesystemmanager.FindFileById(fileId)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": err.Error()})
		return
	}
	c.WriteJSON(file.Content)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		file.Content = string(msg)
		WriteFileBroadCast(filesystemmanager.RootFolder)
	}
}
