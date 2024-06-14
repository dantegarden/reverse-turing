package sse

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	w http.ResponseWriter
}

func (w *responseWriterWrapper) Read(b []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Close() error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Peek(n int) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Skip(n int) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Release() error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Len() int {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) ReadByte() (byte, error) {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) ReadBinary(n int) (p []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Malloc(n int) (buf []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) SetReadTimeout(t time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) SetWriteTimeout(t time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (w *responseWriterWrapper) Finalize() error {
	return w.Finalize()
}

func (w *responseWriterWrapper) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *responseWriterWrapper) Flush() error {
	if f, ok := w.w.(http.Flusher); ok {
		f.Flush()
		return nil
	}
	return fmt.Errorf("http.ResponseWriter does not implement http.Flusher")
}

func (w *responseWriterWrapper) WriteBinary(p [][]byte) (int, error) {
	totalWritten := 0
	for _, b := range p {
		n, err := w.w.Write(b)
		if err != nil {
			return totalWritten, err
		}
		totalWritten += n
	}
	return totalWritten, nil
}

func (w *responseWriterWrapper) WriteString(s string) (int, error) {
	return w.w.Write([]byte(s))
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) SetHeader(key, value string) {
	w.w.Header().Set(key, value)
}

func (w *responseWriterWrapper) StatusCode() int {
	return 0 // You may need to store and return the actual status code if required
}

func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.w.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("http.ResponseWriter does not implement http.Hijacker")
}
