package s3lib

import (
	"bytes"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

func (s3client *S3Client) PutObjectWithOpt(bucketName, key, value string, options ...PutObjectOption) (err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(value)),
	}
	for _, o := range options {
		o(params)
	}
	if _, err = s3client.Client.PutObject(params); err != nil {
		return err
	}
	return
}

type PutObjectOption func(p *s3.PutObjectInput)

func WithContentType(contentType string) PutObjectOption {
	return func(p *s3.PutObjectInput) {
		p.ContentType = aws.String(contentType)
	}
}

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
	if err != nil {
		return err
	}
	return
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
