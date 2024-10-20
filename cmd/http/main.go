package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/michalK00/sg-qr/config"
	_ "github.com/michalK00/sg-qr/docs"
	"github.com/michalK00/sg-qr/internal/collection"
	"github.com/michalK00/sg-qr/internal/config/storage"
	"github.com/michalK00/sg-qr/internal/gallery"
	"github.com/michalK00/sg-qr/internal/qr"
	"github.com/michalK00/sg-qr/pkg/shutdown"
)

// @title Studio Ginger - QR code generator
// @version 0.1
// @description Image gallery that provides uploading images
// @contact.name Michał Klemens
func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run server
	cleanup, err := run(env)

	// run the cleanup after the server is terminated
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	db, err := storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)

	if err != nil {
		return nil, nil, err
	}

	app := fiber.New()

	fmt.Println(env.PORT)

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	// Create user domain

	collectionStore := collection.NewCollectionStorage(db)
	collectionController := collection.NewCollectionController(collectionStore)
	collection.AddCollectionRoutes(app, collectionController, env)

	galleryStore := gallery.NewGalleryStorage(db)
	galleryController := gallery.NewGalleryController(galleryStore)
	gallery.AddGalleryRoutes(app, galleryController, env)

	qrService := qr.QrService{}
	qrController := qr.NewQrController(&qrService)
	qr.AddQrRoutes(app, qrController, env)

	return app, func() {
		storage.CloseMongo(db)
	}, nil
}
