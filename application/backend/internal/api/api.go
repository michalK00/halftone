package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/michalK00/halftone/internal/domain"
	"github.com/michalK00/halftone/internal/fcm"
	"github.com/michalK00/halftone/internal/middleware"
	"github.com/michalK00/halftone/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type api struct {
	collectionRepo domain.CollectionRepository
	galleryRepo    domain.GalleryRepository
	photoRepo      domain.PhotoRepository
	orderRepo      domain.OrderRepository
	jobRepo        domain.JobRepository
	fcmService     fcm.Service
}

func NewApi(db *mongo.Database) *api {
	collectionRepo := repository.NewMongoCollection(db)
	galleryRepo := repository.NewMongoGallery(db)
	photoRepo := repository.NewMongoPhoto(db)
	orderRepo := repository.NewMongoOrder(db)
	jobRepo := repository.NewMongoJob(db)
	jsonCredentials, err := fcm.GetCredentialsJSON()
	if err != nil {
		fmt.Println("Error reading Firebase credentials:", err)
	}
	fcmService, err := fcm.NewService(os.Getenv("FCM_PROJECT_ID"), jsonCredentials)
	if err != nil {
		fmt.Println("Error initializing FCM service:", err)
	}

	return &api{
		collectionRepo: collectionRepo,
		galleryRepo:    galleryRepo,
		photoRepo:      photoRepo,
		orderRepo:      orderRepo,
		jobRepo:        jobRepo,
		fcmService:     *fcmService,
	}
}

func (a *api) Routes(app *fiber.App) {

	auth := app.Group("/api/v1/auth")
	auth.Post("/signup", a.SignUp)
	auth.Post("/signin", a.SignIn)
	auth.Post("/verify", a.VerifyAccount)
	auth.Post("/refresh-token", a.RefreshToken)

	public := app.Group("/api/v1")
	public.Get("/qr", a.generateQrHandler)

	push := app.Group("/push")
	push.Post("/subscribe", a.SubscribeToPush)
	push.Post("/send", a.SendPushMessage)

	//client endpoints protected by middleware that checks if an access token was sent and if it matches the one stored in the accessed db
	client := app.Group("/api/v1/client/galleries/:galleryId", middleware.AuthenticateClient(a.galleryRepo))
	client.Get("", a.clientGetGalleryHandler)
	client.Post("", a.clientCreateOrderHandler)
	client.Get("/photos", a.clientGetGalleryPhotosHandler)
	client.Get("/photos/:photoId", a.clientGetPhotoHandler)

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

	// user endpoints to browse and handle client orders
	protected.Get("/orders", a.getOrdersHandler)
	protected.Get("/orders/:orderId", a.getOrderHandler)
	protected.Put("/orders/:orderId", a.updateOrderHandler)
	protected.Delete("/orders/:orderId", a.deleteOrderHandler)

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
