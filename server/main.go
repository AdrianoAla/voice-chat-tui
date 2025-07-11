package main

import (
	"adriano/vc/shared"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/v2"
)

type Streamer struct {
	Floats []float64
}

func (t *Streamer) Stream(samples [][2]float64) (n int, ok bool) {
	for i := 0; i < len(samples); i++ {
		samples[i][0] = t.Floats[i]
		samples[i][1] = t.Floats[i]
	}
	return len(samples), true
}

func (t *Streamer) Err() error {
	return nil
}

func read_data(pc net.PacketConn, buf []byte) {

	for {
		_, _, err := pc.ReadFrom(buf)
		if err != nil {
			return
		}

		fmt.Print(buf[:100])
	}
}

func main() {
	pc, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		fmt.Println("hi")
	}
	defer pc.Close()

	buf := make([]byte, 960)
	go read_data(pc, buf)

	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		<-sig
		speaker.Close()
		os.Exit(0)
	}()

	speaker.Init(44100, 960)

	streamer := &Streamer{shared.BytesToFloat64Slice(buf)}

	done := make(chan bool)
	s := beep.Seq(streamer, beep.Callback(func() {
		done <- true
	}))

	speaker.Play(s)
	<-done

}
