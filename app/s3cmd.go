package app

type ListBucketResult struct {
	Buckets []string `json:"buckets"`
	Err     string   `json:"err"`
}

func (a *App) ListBuckets() ListBucketResult {
	res := ListBucketResult{}
	buckets, err := a.S3Client.ListBuckets()
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Buckets = buckets
	return res
}

func (a *App) MakeBucket(bucket string) string {
	err := a.S3Client.MakeBucket(bucket)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (a *App) DeleteBucket(bucket string) string {
	err := a.S3Client.DeleteBucket(bucket)
	if err != nil {
		return err.Error()
	}
	return ""
}
