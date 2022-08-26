package s3lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"io"
	"xBrowser/app/db"
)

func (s3client *S3Client) UploadObject(ctx aws.Context, body io.ReadSeeker, task *db.UploadTask) (out *UploadOutput, err error) {
	params := &UploadInput{
		Body:       body,
		Bucket:     aws.String(task.Bucket),
		Key:        aws.String(task.Key),
		UploadTask: task,
	}
	return s3client.Uploader.UploadWithContext(ctx, params)
}

func (s3client *S3Client) DownloadObject(ctx aws.Context, w io.WriterAt, task *db.DownloadTask) (int64, error) {
	params := &GetObjectInput{
		Bucket:       aws.String(task.Bucket),
		Key:          aws.String(task.Key),
		DownloadTask: task,
	}
	return s3client.Downloader.DownloadWithContext(ctx, w, params)
}

func (s3client *S3Client) PutObject(bucketName, key string, body io.ReadSeeker) (err error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   body,
	}
	_, err = s3client.Client.PutObject(params)
	return err
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
