package s3lib

import (
	"fmt"
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/aws/awserr"
	"github.com/journeymidnight/aws-sdk-go/aws/awsutil"
	"github.com/journeymidnight/aws-sdk-go/aws/client"
	"github.com/journeymidnight/aws-sdk-go/aws/request"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
	"github.com/journeymidnight/aws-sdk-go/service/s3/s3iface"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"xBrowser/app/db"
	. "xBrowser/app/models"
)

// DefaultDownloadPartSize is the default range of bytes to get at a time when
// using Download().
const DefaultDownloadPartSize = 1024 * 1024 * 5

// DefaultDownloadConcurrency is the default number of goroutines to spin up
// when using Download().
const DefaultDownloadConcurrency = 5

// The Downloader structure that calls Download(). It is safe to call Download()
// on this structure for multiple objects and across concurrent goroutines.
// Mutating the Downloader's properties is not safe to be done concurrently.
type Downloader struct {
	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultDownloadPartSize value will be used.
	//
	// PartSize is ignored if the Range input parameter is provided.
	PartSize int64

	// The number of goroutines to spin up in parallel when sending parts.
	// If this is set to zero, the DefaultDownloadConcurrency value will be used.
	//
	// Concurrency of 1 will download the parts sequentially.
	//
	// Concurrency is ignored if the Range input parameter is provided.
	Concurrency int

	// An S3 client to use when performing downloads.
	S3 s3iface.S3API

	// List of request options that will be passed down to individual API
	// operation requests made by the downloader.
	RequestOptions []request.Option
}

// WithDownloaderRequestOptions appends to the Downloader's API request options.
func WithDownloaderRequestOptions(opts ...request.Option) func(*Downloader) {
	return func(d *Downloader) {
		d.RequestOptions = append(d.RequestOptions, opts...)
	}
}

// NewDownloader creates a new Downloader instance to downloads objects from
// S3 in concurrent chunks. Pass in additional functional options  to customize
// the downloader behavior. Requires a client.ConfigProvider in order to create
// a S3 service client. The session.Session satisfies the client.ConfigProvider
// interface.
//
// Example:
//     // The session the S3 Downloader will use
//     sess := session.Must(session.NewSession())
//
//     // Create a downloader with the session and default options
//     downloader := s3manager.NewDownloader(sess)
//
//     // Create a downloader with the session and custom options
//     downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
//          d.PartSize = 64 * 1024 * 1024 // 64MB per part
//     })
func NewDownloader(c client.ConfigProvider, options ...func(*Downloader)) *Downloader {
	d := &Downloader{
		S3:          s3.New(c),
		PartSize:    DefaultDownloadPartSize,
		Concurrency: DefaultDownloadConcurrency,
	}
	for _, option := range options {
		option(d)
	}

	return d
}

type maxRetrier interface {
	MaxRetries() int
}

type GetObjectInput struct {
	_ struct{} `type:"structure"`

	DownloadTask *DownloadTask

	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	// Return the object only if its entity tag (ETag) is the same as the one specified,
	// otherwise return a 412 (precondition failed).
	IfMatch *string `location:"header" locationName:"If-Match" type:"string"`

	// Return the object only if it has been modified since the specified time,
	// otherwise return a 304 (not modified).
	IfModifiedSince *time.Time `location:"header" locationName:"If-Modified-Since" type:"timestamp"`

	// Return the object only if its entity tag (ETag) is different from the one
	// specified, otherwise return a 304 (not modified).
	IfNoneMatch *string `location:"header" locationName:"If-None-Match" type:"string"`

	// Return the object only if it has not been modified since the specified time,
	// otherwise return a 412 (precondition failed).
	IfUnmodifiedSince *time.Time `location:"header" locationName:"If-Unmodified-Since" type:"timestamp"`

	// Key is a required field
	Key *string `location:"uri" locationName:"Key" min:"1" type:"string" required:"true"`

	// Part number of the object being read. This is a positive integer between
	// 1 and 10,000. Effectively performs a 'ranged' GET request for the part specified.
	// Useful for downloading just a part of an object.
	PartNumber *int64 `location:"querystring" locationName:"partNumber" type:"integer"`

	// Downloads the specified range bytes of an object. For more information about
	// the HTTP Range header, go to http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35.
	Range *string `location:"header" locationName:"Range" type:"string"`

	// Confirms that the requester knows that she or he will be charged for the
	// request. Bucket owners need not specify this parameter in their requests.
	// Documentation on downloading objects from requester pays buckets can be found
	// at http://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
	RequestPayer *string `location:"header" locationName:"x-amz-request-payer" type:"string" enum:"RequestPayer"`

	// Sets the Cache-Control header of the response.
	ResponseCacheControl *string `location:"querystring" locationName:"response-cache-control" type:"string"`

	// Sets the Content-Disposition header of the response
	ResponseContentDisposition *string `location:"querystring" locationName:"response-content-disposition" type:"string"`

	// Sets the Content-Encoding header of the response.
	ResponseContentEncoding *string `location:"querystring" locationName:"response-content-encoding" type:"string"`

	// Sets the Content-Language header of the response.
	ResponseContentLanguage *string `location:"querystring" locationName:"response-content-language" type:"string"`

	// Sets the Content-Name header of the response.
	ResponseContentType *string `location:"querystring" locationName:"response-content-type" type:"string"`

	// Sets the Expires header of the response.
	ResponseExpires *time.Time `location:"querystring" locationName:"response-expires" type:"timestamp"`

	// Specifies the algorithm to use to when encrypting the object (e.g., AES256).
	SSECustomerAlgorithm *string `location:"header" locationName:"x-amz-server-side-encryption-customer-algorithm" type:"string"`

	// Specifies the customer-provided encryption key for Amazon S3 to use in encrypting
	// data. This value is used to store the object and then it is discarded; Amazon
	// does not store the encryption key. The key must be appropriate for use with
	// the algorithm specified in the x-amz-server-side​-encryption​-customer-algorithm
	// header.
	SSECustomerKey *string `location:"header" locationName:"x-amz-server-side-encryption-customer-key" type:"string" sensitive:"true"`

	// Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321.
	// Amazon S3 uses this header for a message integrity check to ensure the encryption
	// key was transmitted without error.
	SSECustomerKeyMD5 *string `location:"header" locationName:"x-amz-server-side-encryption-customer-key-MD5" type:"string"`

	// VersionId used to reference a specific version of the object.
	VersionId *string `location:"querystring" locationName:"versionId" type:"string"`
}

// Download downloads an object in S3 and writes the payload into w using
// concurrent GET requests.
//
// Additional functional options can be provided to configure the individual
// download. These options are copies of the Downloader instance Download is called from.
// Modifying the options will not impact the original Downloader instance.
//
// It is safe to call this method concurrently across goroutines.
//
// The w io.WriterAt can be satisfied by an os.File to do multipart concurrent
// downloads, or in memory []byte wrapper using aws.WriteAtBuffer.
//
// Specifying a Downloader.Concurrency of 1 will cause the Downloader to
// download the parts from S3 sequentially.
//
// If the GetObjectInput's Range value is provided that will cause the downloader
// to perform a single GetObjectInput request for that object's range. This will
// caused the part size, and concurrency configurations to be ignored.
func (d Downloader) Download(w io.WriterAt, input *GetObjectInput, options ...func(*Downloader)) (n int64, err error) {
	return d.DownloadWithContext(aws.BackgroundContext(), w, input, options...)
}

// DownloadWithContext downloads an object in S3 and writes the payload into w
// using concurrent GET requests.
//
// DownloadWithContext is the same as Download with the additional support for
// Context input parameters. The Context must not be nil. A nil Context will
// cause a panic. Use the Context to add deadlining, timeouts, etc. The
// DownloadWithContext may create sub-contexts for individual underlying
// requests.
//
// Additional functional options can be provided to configure the individual
// download. These options are copies of the Downloader instance Download is
// called from. Modifying the options will not impact the original Downloader
// instance. Use the WithDownloaderRequestOptions helper function to pass in request
// options that will be applied to all API operations made with this downloader.
//
// The w io.WriterAt can be satisfied by an os.File to do multipart concurrent
// downloads, or in memory []byte wrapper using aws.WriteAtBuffer.
//
// Specifying a Downloader.Concurrency of 1 will cause the Downloader to
// download the parts from S3 sequentially.
//
// It is safe to call this method concurrently across goroutines.
//
// If the GetObjectInput's Range value is provided that will cause the downloader
// to perform a single GetObjectInput request for that object's range. This will
// caused the part size, and concurrency configurations to be ignored.
func (d Downloader) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *GetObjectInput, options ...func(*Downloader)) (n int64, err error) {
	impl := downloader{w: w, in: input, cfg: d, ctx: ctx}

	for _, option := range options {
		option(&impl.cfg)
	}
	impl.cfg.RequestOptions = append(impl.cfg.RequestOptions, request.WithAppendUserAgent("S3Manager"))

	if s, ok := d.S3.(maxRetrier); ok {
		impl.partBodyMaxRetries = s.MaxRetries()
	}

	impl.totalBytes = -1
	if impl.cfg.Concurrency == 0 {
		impl.cfg.Concurrency = DefaultDownloadConcurrency
	}

	if impl.cfg.PartSize == 0 {
		impl.cfg.PartSize = DefaultDownloadPartSize
	}

	return impl.download()
}

// downloader is the implementation structure used internally by Downloader.
type downloader struct {
	ctx aws.Context
	cfg Downloader

	in *GetObjectInput
	w  io.WriterAt

	wg sync.WaitGroup
	m  sync.Mutex

	pos        int64
	totalBytes int64
	written    int64
	err        error

	partBodyMaxRetries int
}

// download performs the implementation of the object download across ranged
// GETs.
func (d *downloader) download() (n int64, err error) {
	// If range is specified fall back to single download of that range
	// this enables the functionality of ranged gets with the downloader but
	// at the cost of no multipart downloads.
	if rng := aws.StringValue(d.in.Range); len(rng) > 0 {
		d.downloadRange(rng)
		return d.written, d.err
	}

	if len(d.in.DownloadTask.CompletedPart) != 0 {
		d.totalBytes = d.in.DownloadTask.Size
		for _, v := range d.in.DownloadTask.CompletedPart {
			d.incrWritten(v.PartSize)
			d.pos += v.PartSize
		}
	} else {
		err = db.GlobalAppDB.CreateDownloadTask(d.in.DownloadTask)
		if err != nil {
			return 0, fmt.Errorf("CreateDownloadTask err: %s", err.Error())
		}
		// Spin off first worker to check additional header information
		d.getChunk()
	}

	if total := d.getTotalBytes(); total >= 0 {
		// Spin up workers
		ch := make(chan dlchunk, d.cfg.Concurrency)

		for i := 0; i < d.cfg.Concurrency; i++ {
			d.wg.Add(1)
			go d.downloadPart(ch)
		}

		// Assign work
		for d.getErr() == nil {
			if d.pos >= total {
				break // We're finished queuing chunks
			}

			// Queue the next range of bytes to read.
			ch <- dlchunk{w: d.w, start: d.pos, size: d.cfg.PartSize}

			d.pos += d.cfg.PartSize
		}

		// Wait for completion
		close(ch)
		d.wg.Wait()
	} else {
		// Checking if we read anything new
		for d.err == nil {
			d.getChunk()
		}

		// We expect a 416 error letting us know we are done downloading the
		// total bytes. Since we do not know the content's length, this will
		// keep grabbing chunks of data until the range of bytes specified in
		// the request is out of range of the content. Once, this happens, a
		// 416 should occur.
		e, ok := d.err.(awserr.RequestFailure)
		if ok && e.StatusCode() == http.StatusRequestedRangeNotSatisfiable {
			d.err = nil
		}
	}

	if d.err == nil {
		db.GlobalAppDB.DeleteDownloadTask(d.in.DownloadTask.AccountId, d.in.DownloadTask.TaskId)
	}
	// Return error
	return d.written, d.err
}

// downloadPart is an individual goroutine worker reading from the ch channel
// and performing a GetObject request on the data with a given byte range.
//
// If this is the first worker, this operation also resolves the total number
// of bytes to be read so that the worker manager knows when it is finished.
func (d *downloader) downloadPart(ch chan dlchunk) {
	defer d.wg.Done()
	for {
		chunk, ok := <-ch
		if !ok {
			break
		}
		if d.getErr() != nil {
			// Drain the channel if there is an error, to prevent deadlocking
			// of download producer.
			continue
		}

		if err := d.downloadChunk(chunk); err != nil {
			d.setErr(err)
		} else {
			err = db.GlobalAppDB.CreateDownloadPart(d.in.DownloadTask.TaskId, chunk.start, chunk.size)
			if err != nil {
				d.setErr(err)
			}
		}
	}
}

// getChunk grabs a chunk of data from the body.
// Not thread safe. Should only used when grabbing data on a single thread.
func (d *downloader) getChunk() {
	if d.getErr() != nil {
		return
	}

	chunk := dlchunk{w: d.w, start: d.pos, size: d.cfg.PartSize}
	d.pos += d.cfg.PartSize

	if err := d.downloadChunk(chunk); err != nil {
		d.setErr(err)
	}
}

// downloadRange downloads an Object given the passed in Byte-Range value.
// The chunk used down download the range will be configured for that range.
func (d *downloader) downloadRange(rng string) {
	if d.getErr() != nil {
		return
	}

	chunk := dlchunk{w: d.w, start: d.pos}
	// Ranges specified will short circuit the multipart download
	chunk.withRange = rng

	if err := d.downloadChunk(chunk); err != nil {
		d.setErr(err)
	}

	// Update the position based on the amount of data received.
	d.pos = d.written
}

// downloadChunk downloads the chunk from s3
func (d *downloader) downloadChunk(chunk dlchunk) error {
	in := &s3.GetObjectInput{}
	awsutil.Copy(in, d.in)

	// Get the next byte range of data
	in.Range = aws.String(chunk.ByteRange())

	var n int64
	var err error
	for retry := 0; retry <= d.partBodyMaxRetries; retry++ {
		var resp *s3.GetObjectOutput
		resp, err = d.cfg.S3.GetObjectWithContext(d.ctx, in, d.cfg.RequestOptions...)
		if err != nil {
			return err
		}
		d.setTotalBytes(resp) // Set total if not yet set.

		n, err = io.Copy(&chunk, resp.Body)
		resp.Body.Close()
		if err == nil {
			break
		}

		chunk.cur = 0
		logMessage(d.cfg.S3, aws.LogDebugWithRequestRetries,
			fmt.Sprintf("DEBUG: object part body download interrupted %s, err, %v, retrying attempt %d",
				aws.StringValue(in.Key), err, retry))
	}

	d.incrWritten(n)

	return err
}

// getTotalBytes is a thread-safe getter for retrieving the total byte status.
func (d *downloader) getTotalBytes() int64 {
	d.m.Lock()
	defer d.m.Unlock()

	return d.totalBytes
}

// setTotalBytes is a thread-safe setter for setting the total byte status.
// Will extract the object's total bytes from the Content-Range if the file
// will be chunked, or Content-Length. Content-Length is used when the response
// does not include a Content-Range. Meaning the object was not chunked. This
// occurs when the full file fits within the PartSize directive.
func (d *downloader) setTotalBytes(resp *s3.GetObjectOutput) {
	d.m.Lock()
	defer d.m.Unlock()

	if d.totalBytes >= 0 {
		return
	}

	if resp.ContentRange == nil {
		// ContentRange is nil when the full file contents is provided, and
		// is not chunked. Use ContentLength instead.
		if resp.ContentLength != nil {
			d.totalBytes = *resp.ContentLength
			return
		}
	} else {
		parts := strings.Split(*resp.ContentRange, "/")

		total := int64(-1)
		var err error
		// Checking for whether or not a numbered total exists
		// If one does not exist, we will assume the total to be -1, undefined,
		// and sequentially download each chunk until hitting a 416 error
		totalStr := parts[len(parts)-1]
		if totalStr != "*" {
			total, err = strconv.ParseInt(totalStr, 10, 64)
			if err != nil {
				d.err = err
				return
			}
		}

		d.totalBytes = total
	}
}

func (d *downloader) incrWritten(n int64) {
	d.m.Lock()
	defer d.m.Unlock()

	d.written += n
}

// getErr is a thread-safe getter for the error object
func (d *downloader) getErr() error {
	d.m.Lock()
	defer d.m.Unlock()

	return d.err
}

// setErr is a thread-safe setter for the error object
func (d *downloader) setErr(e error) {
	d.m.Lock()
	defer d.m.Unlock()

	d.err = e
}

// dlchunk represents a single chunk of data to write by the worker routine.
// This structure also implements an io.SectionReader style interface for
// io.WriterAt, effectively making it an io.SectionWriter (which does not
// exist).
type dlchunk struct {
	w     io.WriterAt
	start int64
	size  int64
	cur   int64

	// specifies the byte range the chunk should be downloaded with.
	withRange string
}

// Write wraps io.WriterAt for the dlchunk, writing from the dlchunk's start
// position to its end (or EOF).
//
// If a range is specified on the dlchunk the size will be ignored when writing.
// as the total size may not of be known ahead of time.
func (c *dlchunk) Write(p []byte) (n int, err error) {
	if c.cur >= c.size && len(c.withRange) == 0 {
		return 0, io.EOF
	}

	n, err = c.w.WriteAt(p, c.start+c.cur)
	c.cur += int64(n)

	return
}

// ByteRange returns a HTTP Byte-Range header value that should be used by the
// client to request the chunk's range.
func (c *dlchunk) ByteRange() string {
	if len(c.withRange) != 0 {
		return c.withRange
	}

	return fmt.Sprintf("bytes=%d-%d", c.start, c.start+c.size-1)
}
