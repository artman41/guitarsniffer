package main

import (
	"encoding/hex"
	"fmt"

	"github.com/artman41/guitarsniffer/sniffer"
)

const XboxHeaderLength = 22

func main() {
	sniffer, err := sniffer.Start()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case packet := <-sniffer.Packets:
			handle_packet(packet)
		}
	}
}

func handle_packet(packet sniffer.Packet) {
	if packet.CaptureInfo.Length != 40 {
		return
	}
	createGuitarPacket(packet.Data[XboxHeaderLength:])
}

func createGuitarPacket(packet []byte) GuitarPacket {
	fmt.Printf("%d %s\n", len(packet), hex.EncodeToString(packet))
	upperFrets := Frets{
		Green: packet[11],
	}
	lowerFrets := Frets{}
	return GuitarPacket{
		UpperFrets: upperFrets,
		LowerFrets: lowerFrets,
	}
}

getFrets(fretBitMask byte)

type Frets struct {
	Green, Red, Yellow, Blue, Orange byte
}

type GuitarPacket struct {
	UpperFrets, LowerFrets Frets
	Dpad                   byte
	Slider                 byte
	MenuButton             bool
	OptionsButton          bool
	XboxButton             bool
	Whammy                 byte
	Tilt                   byte
}
