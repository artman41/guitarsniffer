package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/artman41/guitarsniffer/guitarpacket"
	"github.com/artman41/guitarsniffer/sniffer"
)

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
	// The packet returned when pressing the Xbox button is 31
	// bytes long, not 40, meaning that currently we're
	// ignoring that it exists but the code is there for it
	if packet.CaptureInfo.Length != 40 {
		return
	}
	guitarPacket := guitarpacket.CreateGuitarPacket(packet.Data[guitarpacket.XboxHeaderLength:])
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(false)
	err := enc.Encode(guitarPacket)
	if err != nil {
		return
	}
	fmt.Println(string(buf.Bytes()))
}
