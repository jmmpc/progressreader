package progressreader

import (
	"context"
	"io"
	"sync"
)

// ProgressReader implements io.Reader and has aditional posibility
// to show bytes read number of ongoing read operation.
type ProgressReader interface {
	// Read implements the io.Reader interface.
	Read(b []byte) (n int, err error)
	// Loaded returns the number of bytes that have already been read.
	Loaded() int64
}

type progressReader struct {
	reader io.Reader
	mu     sync.Mutex
	loaded int64
}

func (p *progressReader) Read(b []byte) (n int, err error) {
	n, err = p.reader.Read(b)
	p.mu.Lock()
	p.loaded += int64(n)
	p.mu.Unlock()
	return n, err
}

func (p *progressReader) Loaded() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.loaded
}

// New returns a new ProgressReader that uses r as the underlying reader.
func New(r io.Reader) ProgressReader {
	return &progressReader{reader: r}
}

type progressReaderWithContext struct {
	*progressReader
	ctx context.Context
}

func (p *progressReaderWithContext) Read(b []byte) (n int, err error) {
	select {
	case <-p.ctx.Done():
		return 0, p.ctx.Err()
	default:
		return p.progressReader.Read(b)
	}
}

// WithContext returns a new ProgressReader that uses ctx as context and r as the underlying reader.
func WithContext(ctx context.Context, r io.Reader) ProgressReader {
	return &progressReaderWithContext{&progressReader{reader: r}, ctx}
}
