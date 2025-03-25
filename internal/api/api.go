package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/michalK00/halftone/internal/domain"
	"github.com/michalK00/halftone/internal/middleware"
	"github.com/michalK00/halftone/internal/repository"
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

	auth := app.Group("/auth")
	auth.Post("/signup", a.SignUp)
	auth.Post("/signin", a.SignIn)
	auth.Post("/verify", a.VerifyAccount)
	auth.Post("/refresh-token", a.RefreshToken)

	public := app.Group("/api/v1")
	public.Get("/qr", a.generateQrHandler)
	protected := app.Group("/api/v1", middleware.Protected())

	protected.Get("/collections", a.getCollectionsHandler)
	protected.Post("/collections", a.createCollectionHandler)
	protected.Get("/collections/:collectionId", a.getCollectionHandler)
	protected.Put("/collections/:collectionId", a.updateCollectionHandler)
	protected.Delete("/collections/:collectionId", a.deleteCollectionHandler)

	protected.Get("/collections/:collectionId/galleries", a.getGalleriesHandler)
	protected.Get("/collections/:collectionId/galleryCount", a.getGalleryCountHandler)
	protected.Post("/collections/:collectionId/galleries", a.createGalleryHandler)
	//protected.Post("/:collectionId/galleries/batch")
	//protected.Delete("/:collectionId/galleries/batch")
	protected.Get("/galleries/:galleryId", a.getGalleryHandler)
	protected.Put("/galleries/:galleryId", a.updateGalleryHandler)
	protected.Delete("/galleries/:galleryId", a.deleteGalleryHandler)

	protected.Post("/galleries/:galleryId/sharing/share", a.shareGalleryHandler)
	protected.Put("/galleries/:galleryId/sharing/reschedule", a.rescheduleGallerySharingHandler)
	protected.Put("/galleries/:galleryId/sharing/stop", a.stopSharingGalleryHandler)

	protected.Get("/galleries/:galleryId/photos", a.getPhotosHandler)
	protected.Post("/galleries/:galleryId/photos", a.uploadPhotosHandler)
	//protected.Delete("/galleries/:galleryId/photos")
	//protected.Get("/photos/:photoId")
	protected.Put("/photos/:photoId/confirm", a.confirmPhotoUploadHandler)
	protected.Delete("/photos/:photoId", a.deletePhotoHandler)

	//protected.Get("/orders/:orderId")
	//protected.Put("/orders/:orderId")
	//protected.Delete("/orders/:orderId")

	//client := public.Group("/api/v1/client")
	//client.Get("/galleries/:galleryId")
	//client.Post("/galleries/:galleryId")
	//client.Get("/galleries/:galleryId/photos")
	//client.Get("/galleries/:galleryId/photos/:photoId")

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
