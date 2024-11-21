package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/michalK00/sg-qr/internal/domain"
	"github.com/michalK00/sg-qr/internal/middleware"
	"github.com/michalK00/sg-qr/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type api struct {
	collectionRepo domain.CollectionRepository
	galleryRepo    domain.GalleryRepository
	photoRepo      domain.PhotoRepository
	orderRepo      domain.OrderRepository
	jobRepo        domain.JobRepository
	jobQueue       domain.JobQueue
}

func NewApi(db *mongo.Database, rdb *redis.Client) *api {
	collectionRepo := repository.NewMongoCollection(db)
	galleryRepo := repository.NewMongoGallery(db)
	photoRepo := repository.NewMongoPhoto(db)
	orderRepo := repository.NewMongoOrder(db)
	jobRepo := repository.NewMongoJob(db)

	jobQueue := repository.NewRedisJob(rdb)

	return &api{
		collectionRepo: collectionRepo,
		galleryRepo:    galleryRepo,
		photoRepo:      photoRepo,
		orderRepo:      orderRepo,
		jobRepo:        jobRepo,
		jobQueue:       jobQueue,
	}
}

func (a *api) Routes(app *fiber.App) {
	// authMiddleware := auth.NewAuthMiddleware(config)

	collections := app.Group("/api/v1")
	collections.Get("/collections", a.getCollectionsHandler)
	collections.Post("/collections", a.createCollectionHandler)
	collections.Get("/collections/:collectionId", a.getCollectionHandler)
	collections.Put("/collections/:collectionId", a.updateCollectionHandler)
	collections.Delete("/collections/:collectionId", a.deleteCollectionHandler)

	qr := app.Group("/api/v1")
	qr.Get("/qr", a.generateQrHandler)

	gallery := app.Group("/api/v1")
	gallery.Get("/collections/:collectionId/galleries", a.getGalleriesHandler)
	gallery.Get("/collections/:collectionId/galleryCount", a.getGalleryCountHandler)
	gallery.Post("/collections/:collectionId/galleries", a.createGalleryHandler)
	//gallery.Post("/:collectionId/galleries/batch")
	//gallery.Delete("/:collectionId/galleries/batch")
	gallery.Get("/galleries/:galleryId", a.getGalleryHandler)
	gallery.Put("/galleries/:galleryId", a.updateGalleryHandler)
	gallery.Delete("/galleries/:galleryId", a.deleteGalleryHandler)

	gallery.Post("galleries/:galleryId/share", a.shareGalleryHandler)

	//client := app.Group("/api/v1/client")
	//client.Get("/galleries/:galleryId")
	//client.Post("/galleries/:galleryId")
	//client.Get("/galleries/:galleryId/photos")
	//client.Get("/galleries/:galleryId/photos/:photoId")

	photos := app.Group("/api/v1")
	photos.Get("/galleries/:galleryId/photos", a.getPhotosHandler)
	photos.Post("/galleries/:galleryId/photos", a.uploadPhotosHandler)
	//photos.Delete("/galleries/:galleryId/photos")
	//photos.Get("/photos/:photoId")
	photos.Put("/photos/:photoId/confirm", a.confirmPhotoUploadHandler)
	photos.Delete("/photos/:photoId", a.deletePhotoHandler)

	//order := app.Group("/api/v1")
	//order.Get("/orders/:orderId")
	//order.Put("/orders/:orderId")
	//order.Delete("/orders/:orderId")
}

func (a *api) Server() *fiber.App {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	a.Routes(app)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	return app
}
