package services

import (
	"ai-content-creation/models"
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadImageToS3(imageURL []byte, contentReq *models.ContentRequest) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			*aws.String(os.Getenv("AWS_ACCESS_KEY_ID")),
			*aws.String(os.Getenv("AWS_SECRET_ACCESS_KEY")),
			"",
		),
	}))

	uploader := s3manager.NewUploader(sess)
	imageReader := bytes.NewReader(imageURL)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(contentReq.RequestID),
		Body:        imageReader,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func GetImageURL(key string) string {
	return fmt.Sprintf("https://%s.s3.eu-central-1.amazonaws.com/%s", os.Getenv("S3_BUCKET"), key)
}
