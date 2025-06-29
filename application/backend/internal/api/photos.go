package api

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/michalK00/halftone/internal/aws"
	"github.com/michalK00/halftone/internal/domain"
	awsClient "github.com/michalK00/halftone/platform/cloud/aws"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type photoUploadRequest struct {
	OriginalFilename string `json:"originalFilename"`
}

type photoUploadResponse struct {
	Id                   string                  `json:"id"`
	OriginalFilename     string                  `json:"originalFilename"`
	PresignedPostRequest s3.PresignedPostRequest `json:"presignedPostRequest"`
}

var photoPutObjectConditions = []interface{}{
	[]interface{}{"starts-with", "$Content-Type", "image/"},
	[]interface{}{"content-length-range", 1, 10485760},
}

// @Summary Upload photos to a gallery
// @Description Creates new photo entries in a gallery and returns pre-signed URLs for uploading the actual photo files
// @Tags photos
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID" format(objectId)
// @Param request body []photoUploadRequest true "Photo upload requests"
// @Success 201 {object} []photoUploadResponse "Successfully created photo entries with upload URLs"
// @Failure 400 {object} fiber.Map "Invalid request body or gallery ID"
// @Failure 404 {object} fiber.Map "Gallery not found"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/v1/galleries/{galleryId}/photos [post]
func (a *api) uploadPhotosHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	gallery, err := a.galleryRepo.GetGallery(ctx.Context(), galleryId, userId)
	if err != nil {
		return ServerError(ctx, err, "Server error while retrieving gallery")
	}

	var req []photoUploadRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequest(ctx, err)
	}
	if len(req) == 0 {
		return BadRequest(ctx, errors.New("no photos provided"))
	}
	if len(req) > 30 {
		return BadRequest(ctx, errors.New("too many photos in single request"))
	}

	filenames := make([]string, len(req))
	for i, photo := range req {
		if photo.OriginalFilename == "" {
			return BadRequest(ctx, errors.New("empty filename provided"))
		}
		filenames[i] = photo.OriginalFilename
	}

	photoIds, err := a.photoRepo.CreatePhotos(ctx.Context(), gallery.CollectionId, galleryId, filenames, userId)
	if err != nil {
		return ServerError(ctx, err, "Server error while uploading photos")
	}

	res := make([]photoUploadResponse, len(photoIds))
	for i, photoId := range photoIds {
		ext := filepath.Ext(filenames[i])
		if ext == "" {
			ext = ".jpg"
		}

		objectPath := path.Join(gallery.CollectionId.Hex(), gallery.ID.Hex(), "photos", photoId.Hex()+ext)

		postReq, err := aws.PostObjectRequest(objectPath, photoPutObjectConditions)
		if err != nil {
			_ = a.photoRepo.DeletePhotos(ctx.Context(), photoIds, userId)
			return ServerError(ctx, err, "Failed to get presigned request")
		}

		res[i] = photoUploadResponse{
			Id:                   photoId.Hex(),
			OriginalFilename:     filenames[i],
			PresignedPostRequest: *postReq,
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)

}

type PhotoUploadPayload struct {
	GalleryId string `json:"galleryId"`
	PhotoId   string `json:"photoId"`
	ObjectKey string `json:"objectKey"`
	Bucket    string `json:"bucket"`
}

// @Summary Confirm photo upload
// @Description Confirms that a photo has been successfully uploaded by updating its status
// @Tags photos
// @Accept json
// @Produce json
// @Param photoId path string true "Photo ID (MongoDB ObjectID)" format(objectid)
// @Success 200 {object} domain.PhotoDB
// @Failure 404 {object} fiber.Map "Photo not found or invalid ID"
// @Failure 500 {object} fiber.Map "Server error while confirming upload"
// @Router /api/v1/photos/{photoId}/confirm [put]
func (a *api) confirmPhotoUploadHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	photoId, err := primitive.ObjectIDFromHex(ctx.Params("photoId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	photo, err := a.photoRepo.GetPhoto(ctx.Context(), photoId, userId)
	if err != nil {
		return NotFound(ctx, err)
	}

	if _, err := aws.ObjectExists(photo.ObjectKey); err != nil {
		return NotFound(ctx, err)
	}
	photo, err = a.photoRepo.UpdatePhoto(ctx.Context(), photoId, domain.PhotoStatus(1), userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to confirm photo upload")
	}

	client, err := awsClient.GetClient()
	if err != nil {
		return ServerError(ctx, err, "Failed to get AWS client")
	}
	payload := PhotoUploadPayload{
		GalleryId: photo.GalleryId.Hex(),
		PhotoId:   photo.ID.Hex(),
		ObjectKey: photo.ObjectKey,
		Bucket:    os.Getenv("AWS_S3_NAME"),
	}
	lambdaPayload := &awsClient.LambdaPayload{
		EventType: "photo.uploaded",
		Payload:   payload,
		Metadata: map[string]string{
			"environment": os.Getenv("ENV"),
		},
		DelaySeconds: 0,
	}
	message, err := client.SQS.SendLambdaPayload(ctx.Context(), lambdaPayload)
	if err != nil {
		log.Printf("Failed to trigger photo processing: %v", err)
		return err
	}
	log.Printf("Photo processing triggered. MessageID: %s", *message.MessageId)

	return ctx.Status(fiber.StatusOK).JSON(photo)
}

type getPhotoResponse struct {
	Id               string             `json:"id"`
	OriginalFilename string             `json:"originalFilename"`
	Url              string             `json:"url"`
	ThumbnailUrl     string             `json:"thumbnailUrl"`
	Status           domain.PhotoStatus `json:"status"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	CreatedAt        time.Time          `json:"createdAt"`
}

// @Summary Get gallery photos
// @Description Retrieves all photos from a specific gallery, including their signed URLs
// @Tags photos
// @Accept json
// @Produce json
// @Param galleryId path string true "Gallery ID (MongoDB ObjectID)" format(objectid)
// @Success 200 {array} getPhotoResponse "Array of photos with signed URLs"
// @Failure 404 {object} fiber.Map "Gallery not found or invalid ID"
// @Failure 500 {object} fiber.Map "Server error while retrieving photos"
// @Router /api/v1/galleries/{galleryId}/photos [get]
// @Response 200 {object} getPhotoResponse
// TODO add thumbnail logic
func (a *api) getPhotosHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	galleryId, err := primitive.ObjectIDFromHex(ctx.Params("galleryId"))
	if err != nil {
		return NotFound(ctx, err)
	}
	photos, err := a.photoRepo.GetPhotos(ctx.Context(), galleryId, userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to get photos")
	}

	res := make([]getPhotoResponse, len(photos))
	for i, photo := range photos {
		exists, err := aws.ObjectExists(photo.ThumbnailObjectKey)

		url, err := aws.GetObjectUrl(photo.ObjectKey)
		if err != nil {
			return ServerError(ctx, err, "Failed to get photo url")
		}

		var thumbnailUrl string
		if exists {
			thumbnailUrl, err = aws.GetObjectUrl(photo.ThumbnailObjectKey)
			if err != nil {
				return ServerError(ctx, err, "Failed to get thumbnail url")
			}
		}

		res[i] = getPhotoResponse{
			Id:               photo.ID.Hex(),
			OriginalFilename: photo.OriginalFilename,
			Url:              url,
			ThumbnailUrl:     thumbnailUrl,
			Status:           photo.Status,
			UpdatedAt:        photo.UpdatedAt,
			CreatedAt:        photo.CreatedAt,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// @Summary Delete photo
// @Description Soft delete a given photo
// @Tags photos
// @Accept json
// @Produce json
// @Param photoId path string true "Photo ID (MongoDB ObjectID)" format(objectid)
// @Success 200 {object} nil "Photo successfully deleted"
// @Failure 404 {object} fiber.Map "Photo not found or invalid ID"
// @Failure 500 {object} fiber.Map "Server error while deleting photo"
// @Router /api/v1/photos/{photoId} [delete]
func (a *api) deletePhotoHandler(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	photoId, err := primitive.ObjectIDFromHex(ctx.Params("photoId"))
	if err != nil {
		return NotFound(ctx, err)
	}

	err = a.photoRepo.SoftDeletePhoto(ctx.Context(), photoId, userId)
	if err != nil {
		return ServerError(ctx, err, "Failed to delete photo")
	}
	return ctx.Status(fiber.StatusOK).JSON(nil)
}
