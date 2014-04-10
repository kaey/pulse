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

func Example_wav() {
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
		Format:   pulse.FormatS16le,
		Rate:     rd.SampleRate,
		Channels: rd.Channels,
	}

	w, err := pulse.NewWriter(&s, "player", file.Name())
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()
	defer w.Drain()

	buf := make([]byte, 2048)
	for {
		n, err := rd.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		if n == 0 {
			break
		}

		err = w.Write(buf[:n])
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func Example_echo() {
	s := pulse.Sample{
		Format:   pulse.FormatS16le,
		Rate:     44100,
		Channels: 2,
	}

	r, err := pulse.NewReader(&s, "echo", "mic")
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()
	defer r.Drain()
	
	rl, err := r.Latency()
	if err != nil {
		log.Fatalln(err)
	}
	
	w, err := pulse.NewWriter(&s, "echo", "play")
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()
	defer w.Drain()
	
	wl, err := w.Latency()
	if err != nil {
		log.Fatalln(err)
	}
	
	log.Println("Reader latency", rl)
	log.Println("Writer latency", wl)

	buf := make([]byte, 2048)
	for {
		err := r.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}

		err = w.Write(buf)
		if err != nil {
			log.Fatalln(err)
		}
	}
}