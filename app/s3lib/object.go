package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"github.com/journeymidnight/aws-sdk-go/service/s3/s3manager"

	//"github.com/journeymidnight/aws-sdk-go/service/s3/s3manager"
	"io"
)

func (s3client *S3Client) PutObjectWithOpt(bucketName, key string, body io.Reader, options ...UploadOption) (err error) {
	params := &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	for _, o := range options {
		o(params)
	}
	_, err = s3client.Uploader.Upload(params)
	return err
}

type UploadOption func(p *s3manager.UploadInput)

func (s3client *S3Client) GetObjectOutPut(bucketName, key string) (out *s3.GetObjectOutput, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	return s3client.Client.GetObject(params)
}

func (s3client *S3Client) DeleteObject(bucketName, key string) (err error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	_, err = s3client.Client.DeleteObject(params)
	return err
}

func (s3client *S3Client) DeleteObjects(bucketName string, keys map[string]string) (err error) {
	var objects []*s3.ObjectIdentifier
	objects = []*s3.ObjectIdentifier{}
	for k, v := range keys {
		object := &s3.ObjectIdentifier{
			Key:       aws.String(k),
			VersionId: aws.String(v),
		}
		objects = append(objects, object)
	}
	delete := &s3.Delete{
		Objects: objects,
	}
	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: delete,
	}
	_, err = s3client.Client.DeleteObjects(params)
	if err != nil {
		return err
	}
	return
}
