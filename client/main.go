package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/MarkKremer/microphone"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type TapStreamer struct {
	Streamer       beep.Streamer
	Buffer         *bytes.Buffer
	InternetBuffer []byte
	File           *os.File
}

func (t *TapStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = t.Streamer.Stream(samples)
	for i := 0; i < n; i++ {
		left := int16(samples[i][0] * 32767)

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(left))

		t.InternetBuffer = append(t.InternetBuffer, b...)

		t.Buffer.Write(b)
		if t.File != nil {
			_, err := t.File.Write(b)
			if err != nil {
				log.Printf("Failed to write to file: %v\n", err)
			}
		}

	}
	return n, ok
}

func (t *TapStreamer) Err() error {
	return t.Streamer.Err()
}

// "github.com/gopxl/beep/wav"
func main() {

	fmt.Println("Recording. Press Ctrl-C to stop.")

	err := microphone.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer microphone.Terminate()

	stream, _, err := microphone.OpenDefaultStream(44100, 1)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("recording.pcm")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Close the stream at the end if it hasn't already been
	// closed explicitly.

	defer stream.Close()

	stream.Start()
	var buffer bytes.Buffer
	tap := &TapStreamer{
		Streamer: stream,
		Buffer:   &buffer,
		File:     file,
	}

	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		<-sig
		stream.Stop()
		stream.Close()
		os.Exit(0)
		fmt.Println(buffer.Bytes())
	}()

	speaker.Init(44100, 960)

	done := make(chan bool)
	s := beep.Seq(tap, beep.Callback(func() {
		done <- true
	}))

	// samples := make([][2]float64, 960)

	// fmt.Println("ENDING UFFER!!!")
	// for {
	// 	_, _ = tap.Stream(samples)

	// 	bytes := shared.Float64SliceToBytes(samples)

	// 	udpAddr, _ := net.ResolveUDPAddr("udp", ":1053")
	// 	c, _ := net.DialUDP("udp", nil, udpAddr)

	// 	fmt.Println("SENDING BUFFER!!!")
	// 	c.Write(bytes)
	// }

	speaker.Play(s)
	<-done

}
