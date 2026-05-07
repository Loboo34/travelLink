package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloudinaryClient *cloudinary.Cloudinary

func InitCloudinary(cloudName, apiKey, secretKey string) error {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, secretKey)
	if err != nil {
		return fmt.Errorf("Failed cloudinary init: %w", err)
	}

	cloudinaryClient = cld

	return nil
}

func UploadImage(ctx context.Context, file multipart.File, folder string) (string, error) {
	resp, err := cloudinaryClient.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		ResourceType:   "image",
		AllowedFormats: []string{"jpg", "jpeg", "png", "webp"},
	})

	if err != nil {
		return "", fmt.Errorf("Cloudinary upload failed: %w", err)
	}

	return resp.SecureURL, nil
}
