package app

import (
	"io"
)

type Progress struct {
	CurrentBytes int64
	TotalBytes   int64

	// TODO: implement speed
	// Speed          string
}

func (p *Progress) Add(num int64) {
	p.CurrentBytes += num
}

func (p *Progress) State() int64 {
	if p.TotalBytes == 0 {
		return 0
	}
	return p.CurrentBytes * 100 / p.TotalBytes
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
	n, err = r.Reader.Read(p)
	r.p.Add(int64(n))
	return
}
