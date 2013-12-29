package main

import (
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// Add "nl" to stream, after every "len" bytes
type LineBreaker struct {
	len, rem int
	nl       []byte
	w        io.Writer
}

func NewLineBreaker(w io.Writer, len int, nl string) *LineBreaker {
	return &LineBreaker{len: len, rem: 0, nl: []byte(nl), w: w}
}

func (lb *LineBreaker) Write(p []byte) (int, error) {
	var n, wn, count int
	var err error

	count = 0
	for n = len(p); n > 0; n = len(p) {
		if lb.rem == 0 {
			_, err = lb.w.Write(lb.nl)
			if err != nil {
				return count, err
			}
			lb.rem = lb.len
		}
		if n >= lb.rem {
			wn = lb.rem
		} else {
			wn = n
		}
		wn, err = lb.w.Write(p[:wn])
		lb.rem -= wn
		p = p[wn:]
		count += wn
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func emitFileHeader(w io.Writer, fname string, sz int, zip bool) error {
	_, err := fmt.Fprintf(w, FileHeadFormat, fname, sz, zip)
	return err
}

// io.Writer <- bufio.Writer <- LineBreaker <- base64.Encoder
type GoWriter struct {
	w64 io.WriteCloser
	wlb *LineBreaker
	wb  *bufio.Writer
}

func NewGoWriter(w io.Writer, fname string, sz int) (*GoWriter, error) {
	var err error
	var gw *GoWriter
	err = emitFileHeader(w, fname, sz, false)
	if err != nil {
		return nil, err
	}
	gw = &GoWriter{}
	gw.wb = bufio.NewWriter(w)
	gw.wlb = NewLineBreaker(gw.wb, 76, "\n")
	gw.w64 = base64.NewEncoder(base64.StdEncoding, gw.wlb)
	return gw, nil
}

func (gw *GoWriter) Write(p []byte) (int, error) {
	return gw.w64.Write(p)
}

func (gw *GoWriter) Close() error {
	var err error

	err = gw.w64.Close()
	if err != nil {
		_ = gw.wb.Flush()
		return err
	}
	_, err = fmt.Fprintf(gw.wb, FileFootFormat)
	if err != nil {
		_ = gw.wb.Flush()
		return err
	}
	return gw.wb.Flush()
}

// GoWriter <- gzip.Writer :
//   io.Writer <- bufio.Writer <- LineBreaker <-
//       <- base6.Encoder <- gzip.Writer
type GoZipWriter struct {
	zw *gzip.Writer
	gw *GoWriter
}

func NewGoZipWriter(w io.Writer,
	fname string, sz int) (*GoZipWriter, error) {
	var err error
	var gzw *GoZipWriter

	err = emitFileHeader(w, fname, sz, true)
	if err != nil {
		return nil, err
	}
	gzw = &GoZipWriter{}
	gzw.gw = &GoWriter{}
	gzw.gw.wb = bufio.NewWriter(w)
	gzw.gw.wlb = NewLineBreaker(gzw.gw.wb, 76, "\n")
	gzw.gw.w64 = base64.NewEncoder(base64.StdEncoding, gzw.gw.wlb)
	gzw.zw = gzip.NewWriter(gzw.gw)
	return gzw, nil
}

func (gzw *GoZipWriter) Write(p []byte) (int, error) {
	return gzw.zw.Write(p)
}

func (gzw *GoZipWriter) Close() error {
	var err error
	err = gzw.zw.Close()
	if err != nil {
		_ = gzw.gw.Close()
		return err
	}
	return gzw.gw.Close()
}
