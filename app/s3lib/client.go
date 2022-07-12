package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/aws/credentials"
	"github.com/journeymidnight/aws-sdk-go/aws/session"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"github.com/journeymidnight/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	Client   *s3.S3
	Uploader *s3manager.Uploader
}

func NewS3(Endpoint, AccessKey, SecretKey string) *S3Client {
	creds := credentials.NewStaticCredentials(AccessKey, SecretKey, "")
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			DisableSSL:  aws.Bool(true),
			Endpoint:    aws.String(Endpoint),
			Region:      aws.String("none"),
		},
	),
	)
	// By default make sure a region is specified
	return &S3Client{
		Client:   s3.New(sess),
		Uploader: s3manager.NewUploader(sess),
	}
}
