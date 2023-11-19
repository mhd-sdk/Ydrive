package main

import (
	"log"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mehdiseddik.com/pkg/handlers"
	"github.com/mehdiseddik.com/pkg/middlewares"
)

func main() {
	app := fiber.New()

	app.Use("/ws", middlewares.WebsocketMiddleware)
	app.Use(logger.New())

	// Custom File Writer
	file, err := os.OpenFile("./logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{
		Output: file,
	}))

	app.Get("/ws/arborescence", websocket.New(handlers.Arborescence))

	app.Post("/api/file", handlers.CreateFileHandler)
	app.Post("/api/file/move", handlers.MoveFileHandler)

	app.Put("/api/file/:fileId", handlers.EditFileNameHandler)
	app.Delete("/api/file/:fileId", handlers.DeleteFileHandler)

	app.Post("/api/folder", handlers.CreateFolderHandler)
	app.Put("/api/folder/:folderId", handlers.EditFolderNameHandler)

	log.Fatal(app.Listen(":3000"))
}
