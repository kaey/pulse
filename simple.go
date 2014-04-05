// Copyright 2013 Pulse Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pulse is a client library for PulseAudio.
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
	SampleS16le    = 0x3
	StreamPlayback = 0x1
	StreamRecord   = 0x2
)

type Connection struct {
	simple *C.pa_simple
	errno  Errno
}

type Sample struct {
	Format   int32
	Rate     uint32
	Channels uint8
}

type Errno struct {
	i C.int
}

func (e Errno) Error() string {
	cstr := C.pa_strerror(C.int(e.i))
	return C.GoString(cstr)
}

// New creates new connection to server and returns its identifier.
func New(s *Sample, name string, play string) (*Connection, error) {
	c := new(Connection)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cplay := C.CString(play)
	defer C.free(unsafe.Pointer(cplay))

	c.simple = C.pa_simple_new(
		nil,
		cname,
		StreamPlayback,
		nil,
		cplay,
		(*C.pa_sample_spec)(unsafe.Pointer(s)),
		nil,
		nil,
		&c.errno.i,
	)
	if c.errno.i != C.int(0) {
		return nil, c.errno
	}

	return c, nil
}

// Write writes data to server.
func (c *Connection) Write(data []byte) error {
	cdata := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	C.pa_simple_write(
		c.simple,
		unsafe.Pointer(cdata.Data),
		C.size_t(len(data)),
		&c.errno.i,
	)
	if c.errno.i != C.int(0) {
		return c.errno
	}

	return nil
}

// Read reads data from server.
func (c *Connection) Read(data []byte) error {
	cdata := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	C.pa_simple_read(
		c.simple,
		unsafe.Pointer(cdata.Data),
		C.size_t(len(data)),
		&c.errno.i,
	)
	if c.errno.i != C.int(0) {
		return c.errno
	}

	return nil
}

// Drain waits until all written data is played by the server.
func (c *Connection) Drain() error {
	C.pa_simple_drain(
		c.simple,
		&c.errno.i,
	)
	if c.errno.i != C.int(0) {
		return c.errno
	}

	return nil
}

// Flush discards all data in the buffer.
func (c *Connection) Flush() error {
	C.pa_simple_flush(
		c.simple,
		&c.errno.i,
	)
	if c.errno.i != C.int(0) {
		return c.errno
	}

	return nil
}

// Close closes connection to server.
func (c *Connection) Close() {
	C.pa_simple_free(c.simple)
}
