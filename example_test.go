// Copyright 2013 Pulse Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pulse_test

import (
	"io"
	"log"
	"os"

	"github.com/kaey/pulse"
	"github.com/kaey/wav"
)

func Example() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	rd, err := wav.NewReader(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer rd.Close()

	s := pulse.Sample{
		Format:   pulse.SampleS16le,
		Rate:     rd.SampleRate,
		Channels: rd.Channels,
	}

	c, err := pulse.New(&s, "playing", file.Name())
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	defer c.Drain()

	buf := make([]byte, 1024)
	for {
		n, err := rd.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		if n == 0 {
			break
		}

		err = c.Write(buf[:n])
		if err != nil {
			log.Fatalln(err)
		}
	}
}
