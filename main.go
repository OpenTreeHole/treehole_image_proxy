package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fiber.New()
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(logger.New(logger.Config{TimeZone: "Asia/Shanghai"}))
	router := app.Group("/api")
	router.Get("/upload", GetAuthToken)
	router.Post("/json", UploadImage)

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			log.Println(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
