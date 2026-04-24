package s3lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s3client *S3Client) MakeBucket(bucketName string) (err error) {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err = s3client.Client.CreateBucket(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) DeleteBucket(bucketName string) (err error) {
	params := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err = s3client.Client.DeleteBucket(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) HeadBucket(bucketName string) (err error) {
	params := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err = s3client.Client.HeadBucket(params); err != nil {
		return err
	}
	return
}

func (s3client *S3Client) GetBucketVersioning(bucketName string) (status string, MFADelete string, err error) {
	params := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	}
	out, err := s3client.Client.GetBucketVersioning(params)
	if err != nil {
		return "", "", err
	}
	status = *out.Status
	MFADelete = *out.MFADelete
	return
}

func (s3client *S3Client) PutBucketVersioning(bucketName string, status string, MFADelete string) (err error) {
	params := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status:    aws.String(status),
			MFADelete: aws.String(MFADelete),
		},
	}
	_, err = s3client.Client.PutBucketVersioning(params)
	if err != nil {
		return err
	}
	return
}

type Bucket struct {
	buckename, createtime string
}

func (s3client *S3Client) ListBuckets() (buckets *s3.ListBucketsOutput, err error) {
	params := &s3.ListBucketsInput{}
	out, err := s3client.Client.ListBuckets(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s3client *S3Client) ListObjects(bucketName, marker, prefix string, maxKeys int64, delemiter string) (*s3.ListObjectsOutput, error) {
	params := &s3.ListObjectsInput{
		Bucket:    aws.String(bucketName),
		Marker:    aws.String(marker),
		MaxKeys:   aws.Int64(maxKeys),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delemiter),
	}
	return s3client.Client.ListObjects(params)
}
