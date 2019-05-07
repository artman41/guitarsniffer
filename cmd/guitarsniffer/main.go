package main

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"

	"github.com/artman41/guitarsniffer/guitarjoypad"
	"github.com/artman41/guitarsniffer/guitarpacket"
	"github.com/artman41/guitarsniffer/guitarsniffer"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var currentPacket = guitarpacket.GuitarPacket{}
var currentPacketHex = ""
var guitarJoypad *guitarjoypad.GuitarJoypad
var RunDataThread = true

var threads sync.WaitGroup

var width, height int

func main() {
	go guiThread(&RunDataThread)
	threads.Add(1)
	dataThread(&RunDataThread)

	threads.Wait()
}

func guiThread(runDataThread *bool) {
	defer func() {
		threads.Done()
		*runDataThread = false
	}()

	runtime.LockOSThread()

	setupFontCache()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	width, height := 230, 100

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, "Guitar Sniffer", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	initDisplay(width, height)

	for !window.ShouldClose() {
		display(width, height)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

func dataThread(runDataThread *bool) {
	fmt.Println("Getting Joypad...")
	joypad, err := guitarjoypad.GetJoypad()
	if err != nil {
		panic(err)
	}
	guitarJoypad = joypad
	fmt.Println("Obtained!")
	defer guitarJoypad.Relinquish()

	fmt.Println("Starting Sniffer...")
	sniffer, err := guitarsniffer.Start()
	defer sniffer.Stop()
	if err != nil {
		panic(err)
	}
	for *runDataThread {
		select {
		case packet := <-sniffer.Packets:
			handlePacket(&packet)
		default:
			continue
		}
	}
}

func handlePacket(packet *guitarsniffer.Packet) {
	// The packet returned when pressing the Xbox button is 31
	// bytes long, not 40, meaning that currently we're
	// ignoring that it exists but the code is there for it
	if packet.CaptureInfo.Length != 40 {
		return
	}
	currentPacket = guitarpacket.CreateGuitarPacket(packet.Data[guitarpacket.XboxHeaderLength:])
	currentPacketHex = hex.EncodeToString(packet.Data[guitarpacket.XboxHeaderLength:])
	guitarJoypad.SetValues(currentPacket)
	guitarJoypad.Update()
}
