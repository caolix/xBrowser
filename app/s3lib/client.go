package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/aws/credentials"
	"github.com/journeymidnight/aws-sdk-go/aws/session"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client     *s3.S3
	Uploader   *Uploader
	Downloader *Downloader
}

func NewS3(Endpoint, AccessKey, SecretKey string, UseSSL bool) *S3Client {
	creds := credentials.NewStaticCredentials(AccessKey, SecretKey, "")
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			DisableSSL:  aws.Bool(!UseSSL),
			Endpoint:    aws.String(Endpoint),
			Region:      aws.String("none"),
		},
	),
	)
	// By default make sure a region is specified
	return &S3Client{
		Client:     s3.New(sess),
		Uploader:   NewUploader(sess),
		Downloader: NewDownloader(sess),
	}
}
