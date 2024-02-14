package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lsendoya/Warewise/pkg/logger"
	"mime"
	"mime/multipart"
	"path/filepath"
)

type Service struct {
	S3 *s3.Client
}

func NewAWSService(S3 *s3.Client) Service {
	return Service{S3: S3}
}

func (awsSvc *Service) UploadFile(bName, bKey string, file multipart.File) (string, error) {

	contentType := mime.TypeByExtension(filepath.Ext(bKey))

	object, err := awsSvc.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bName),
		Key:         aws.String(bKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Errorf("error while uploading the file %v", err)
		return "", err
	}
	_ = object

	logger.Info("File uploaded successfully")

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bName, "us-east-1", bKey)
	return url, nil
}
