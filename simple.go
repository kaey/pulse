// Copyright 2013 Pulse Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pulse is a cgo client library for PulseAudio.
package pulse

import (
	"reflect"
	"unsafe"
)

/*
#cgo LDFLAGS: -lpulse-simple -lpulse

#include <stdlib.h>
#include <pulse/simple.h>
#include <pulse/error.h>
*/
import "C"

const (
	FormatS16le    = 0x3
	streamPlayback = 0x1
	streamRecord   = 0x2
)

// Sample specifies sample type for stream.
type Sample struct {
	Format   int32
	Rate     uint32
	Channels uint8
}

// Err is a type returned by all pulse functions.
type Err struct {
	no C.int
}

// Error method calls pa_strerror to get descriptive error message.
func (err Err) Error() string {
	cstr := C.pa_strerror(C.int(err.no))
	return C.GoString(cstr)
}

// Conn is a PulseAudio simple connection.
// Do not use this directly.
type Conn struct {
	simple *C.pa_simple
	err    Err
}

// NewConn creates new connection to PulseAudio server.
func NewConn(s *Sample, appName string, streamName string, streamType int) (*Conn, error) {
	conn := new(Conn)

	cAppName := C.CString(appName)
	defer C.free(unsafe.Pointer(cAppName))

	cStreamName := C.CString(streamName)
	defer C.free(unsafe.Pointer(cStreamName))

	conn.simple = C.pa_simple_new(
		nil,
		cAppName,
		C.pa_stream_direction_t(streamType),
		nil,
		cStreamName,
		(*C.pa_sample_spec)(unsafe.Pointer(s)),
		nil,
		nil,
		&conn.err.no,
	)
	if conn.err.no != C.int(0) {
		return nil, conn.err
	}

	return conn, nil
}

// Latency gets connection latency in usec.
func (conn *Conn) Latency() (uint64, error) {
	clat := C.pa_simple_get_latency(conn.simple, &conn.err.no)
	if conn.err.no != C.int(0) {
		return 0, conn.err
	}
	return uint64(clat), nil
}

// Drain waits until all written data is played by the server.
func (conn *Conn) Drain() error {
	C.pa_simple_drain(
		conn.simple,
		&conn.err.no,
	)
	if conn.err.no != C.int(0) {
		return conn.err
	}

	return nil
}

// Flush discards all data in the server buffer.
func (conn *Conn) Flush() error {
	C.pa_simple_flush(
		conn.simple,
		&conn.err.no,
	)
	if conn.err.no != C.int(0) {
		return conn.err
	}

	return nil
}

// Close closes connection to server.
func (conn *Conn) Close() {
	C.pa_simple_free(conn.simple)
}

type Reader struct {
	*Conn
}

// NewReader creates new connection to server.
func NewReader(s *Sample, appName string, streamName string) (*Reader, error) {
	conn, err := NewConn(s, appName, streamName, streamRecord)
	if err != nil {
		return nil, err
	}
	return &Reader{conn}, nil
}

// Read reads data from server.
func (r *Reader) Read(data []byte) error {
	cdata := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	C.pa_simple_read(
		r.Conn.simple,
		unsafe.Pointer(cdata.Data),
		C.size_t(len(data)),
		&r.Conn.err.no,
	)
	if r.Conn.err.no != C.int(0) {
		return r.Conn.err
	}

	return nil
}

type Writer struct {
	*Conn
}

// NewWriter creates new connection to server.
func NewWriter(s *Sample, appName string, streamName string) (*Writer, error) {
	conn, err := NewConn(s, appName, streamName, streamPlayback)
	if err != nil {
		return nil, err
	}
	return &Writer{conn}, nil
}

// Write writes data to server.
func (w *Writer) Write(data []byte) error {
	cdata := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	C.pa_simple_write(
		w.Conn.simple,
		unsafe.Pointer(cdata.Data),
		C.size_t(len(data)),
		&w.Conn.err.no,
	)
	if w.Conn.err.no != C.int(0) {
		return w.Conn.err
	}

	return nil
}
