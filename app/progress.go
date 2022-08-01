package app

import (
	"io"
	"time"
)

type Progress struct {
	CurrentBytes int64
	TotalBytes   int64

	// TODO: implement speed
	LastCheckedBytes int64
	SpeedBytes       int64
}

func (p *Progress) Add(num int64) {
	p.CurrentBytes += num
}

func (p *Progress) Set(num int64) {
	p.CurrentBytes = num
}

func (p *Progress) State() int64 {
	if p.TotalBytes == 0 {
		return 0
	}
	return p.CurrentBytes * 100 / p.TotalBytes
}

func (p *Progress) calcSpeed(close chan struct{}) {
	t := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-t.C:
			p.SpeedBytes = (p.CurrentBytes - p.LastCheckedBytes) / 10 // Bytes per 1s
			p.LastCheckedBytes = p.CurrentBytes
		case <-close:
			return
		}
	}
}

func (p *Progress) Speed() int64 {
	return p.SpeedBytes
}

// Reader is the progressbar io.Reader struct
type Reader struct {
	io.Reader
	p *Progress
}

func NewProgressReader(r io.Reader, p *Progress) *Reader {
	return &Reader{
		Reader: r,
		p:      p,
	}
}

// Read will read the data and add the number of bytes to the progressbar
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.p.CurrentBytes == 0 {
		c := make(chan struct{})
		go r.p.calcSpeed(c)
	}
	n, err = r.Reader.Read(p)
	r.p.Add(int64(n))
	return
}
