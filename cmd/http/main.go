package main

import (
	"fmt"
	"github.com/gofiber/swagger"
	collections "github.com/michalK00/sg-qr/app/collections"
	galleries "github.com/michalK00/sg-qr/app/galleries"
	"github.com/michalK00/sg-qr/internal/config"
	"github.com/michalK00/sg-qr/internal/middleware"
	"github.com/michalK00/sg-qr/platform/database/mongodb"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/michalK00/sg-qr/docs"
	"github.com/michalK00/sg-qr/internal/shutdown"
)

// @title Studio Ginger - QR code generator
// @version 0.1
// @description Image galleries that provides uploading images
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
	db, err := mongodb.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	app := fiber.New()

	// middleware
	middleware.FiberMiddleware(app)

	// Create user domain
	collectionStore := collections.NewCollectionStorage(db)
	collectionController := collections.NewCollectionController(collectionStore)
	galleryStore := galleries.NewGalleryStorage(db)
	galleryService := galleries.NewGalleryService(galleryStore, collectionStore)
	galleryController := galleries.NewGalleryController(galleryService)

	// Add routes
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	collections.AddRoutes(app, collectionController, env)
	galleries.AddRoutes(app, galleryController, env)

	return app, func() {
		mongodb.CloseMongo(db)
	}, nil
}
