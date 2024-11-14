package main

import (
	"fmt"
	"github.com/gofiber/swagger"
	"github.com/michalK00/sg-qr/internal/api"
	"github.com/michalK00/sg-qr/internal/config"
	"github.com/michalK00/sg-qr/internal/middleware"
	"github.com/michalK00/sg-qr/platform/database/mongodb"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/michalK00/sg-qr/docs"
	"github.com/michalK00/sg-qr/internal/shutdown"
)

// @title Image library
// @version 0.1
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

	a := api.NewApi(db)
	a.Routes(app, env)
	// Add routes
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	return app, func() {
		mongodb.CloseMongo(db)
	}, nil
}
