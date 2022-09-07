package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"time"
)

type Progress struct {
	ctx           context.Context
	eventProgress string

	CurrentBytes int64
	TotalBytes   int64

	// TODO: implement speed
	IsCalc           bool
	LastCheckedBytes int64
	SpeedBytes       int64
	closeCh          chan struct{}
}

func NewProgress(ctx context.Context, eventProgress string, total int64) *Progress {
	return &Progress{
		ctx:           ctx,
		eventProgress: eventProgress,
		TotalBytes:    total,
		closeCh:       make(chan struct{}),
	}
}

func (p *Progress) Add(num int64) {
	p.CurrentBytes += num
}

func (p *Progress) Set(num int64) {
	p.CurrentBytes = num
}

func (p *Progress) Ratio() int64 {
	if p.TotalBytes == 0 {
		return 100
	}
	return p.CurrentBytes * 100 / p.TotalBytes
}

func (p *Progress) Speed() int64 {
	return p.SpeedBytes
}

func (p *Progress) Current() int64 {
	return p.CurrentBytes
}

func (p *Progress) calc() {
	t := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-t.C:
			p.SpeedBytes = (p.CurrentBytes - p.LastCheckedBytes) / 2 // Bytes per 1s
			p.LastCheckedBytes = p.CurrentBytes

			if p.Ratio() < 100 && p.eventProgress != "" {
				runtime.EventsEmit(p.ctx, p.eventProgress, p.Ratio())
			}
		case <-p.closeCh:
			runtime.EventsEmit(p.ctx, p.eventProgress, p.Ratio())
			return
		}
	}
}

// ReadSeeker is the progressbar io.Reader struct
type ReadSeekCloser struct {
	io.ReadSeekCloser
	p *Progress
}

func NewUploadProgressReader(rs io.ReadSeekCloser, p *Progress) *ReadSeekCloser {
	return &ReadSeekCloser{
		ReadSeekCloser: rs,
		p:              p,
	}
}

// Read will read the data and add the number of bytes to the progressbar
func (r *ReadSeekCloser) Read(p []byte) (n int, err error) {
	if !r.p.IsCalc {
		go r.p.calc()
		r.p.IsCalc = true
	}
	n, err = r.ReadSeekCloser.Read(p)
	if err != nil {
		r.Close()
		r.p.closeCh <- struct{}{}
	}
	r.p.Add(int64(n))
	return
}

func (r *ReadSeekCloser) Seek(offset int64, whence int) (n int64, err error) {
	n, err = r.ReadSeekCloser.Seek(offset, whence)
	r.p.Set(offset)
	return n, err
}

func (r *ReadSeekCloser) Close() error {
	return r.ReadSeekCloser.Close()
}

type WriterAt struct {
	io.WriterAt
	p *Progress
}

func NewDownloadProgressWriterAt(w io.WriterAt, p *Progress) *WriterAt {
	return &WriterAt{
		WriterAt: w,
		p:        p,
	}
}

func (w *WriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	if !w.p.IsCalc {
		go w.p.calc()
		w.p.IsCalc = true
	}
	n, err = w.WriterAt.WriteAt(p, offset)
	if err != nil {
		w.p.closeCh <- struct{}{}
	}
	w.p.Add(int64(n))
	return
}
