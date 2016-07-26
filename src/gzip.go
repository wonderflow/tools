package common

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"sync"
)

// GzipPool manages a pool of gzip.Writer.
// The pool uses sync.Pool internally.
type GzipPool struct {
	writerPool sync.Pool
	buffPool   sync.Pool
}

// GetWriter returns gzip.Writer from the pool, or creates a new one
// with gzip.BestCompression if the pool is empty.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.writerPool.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = gzip.NewWriterLevel(dst, 1) // NewWriterLevel(dst, gzip.BestSpeed)
	}
	return writer
}

// ReturnWriter returns a gzip.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *GzipPool) ReturnWriter(writer *gzip.Writer) {
	writer.Close()
	pool.writerPool.Put(writer)
}

func (pool *GzipPool) GetBuffer() (buff *bytes.Buffer) {
	if b := pool.buffPool.Get(); buff != nil {
		buff = b.(*bytes.Buffer)
		buff.Reset()
	} else {
		buff = bytes.NewBuffer([]byte{})
	}
	return buff
}

func (pool *GzipPool) ReturnBuffer(buff *bytes.Buffer) {
	pool.buffPool.Put(buff)
}

// getwriter from gzipPool before get gzip content
// and return it back after
func (pool *GzipPool) Gzip(value string) ([]byte, error) {
	buff := pool.GetBuffer()
	writer := pool.GetWriter(buff)
	_, err := writer.Write([]byte(value))
	// 先return，会flush writer
	pool.ReturnWriter(writer)
	// 然后获取内容
	ret := buff.Bytes()
	pool.ReturnBuffer(buff)
	return ret, err
}

func (pool *GzipPool) UnGzip(value string) ([]byte, error) {
	b, err := gzip.NewReader(strings.NewReader(value))
	if err != nil {
		return nil, err
	}
	defer b.Close()
	buff := pool.GetBuffer()
	defer pool.ReturnBuffer(buff)
	_, err = buff.ReadFrom(b)
	if err != nil {
		return nil, err
	}
	ret := buff.Bytes()
	return ret, nil
}
