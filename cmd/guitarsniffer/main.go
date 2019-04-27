package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/artman41/guitarsniffer/guitarjoypad"
	"github.com/artman41/guitarsniffer/guitarpacket"
	"github.com/artman41/guitarsniffer/guitarsniffer"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

var guitarJoypad *guitarjoypad.GuitarJoypad

const ColorDarkOrange3 = 166

func createParagraphWidget(text string) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = text
	return p
}

func createFretWidget(fret rune, down bool) *ui.Canvas {
	canvas := ui.NewCanvas()
	canvas.Border = true
	var col ui.Color
	switch fret {
	case 'g':
		col = ui.ColorGreen
	case 'r':
		col = ui.ColorRed
	case 'y':
		col = ui.ColorYellow
	case 'b':
		col = ui.ColorBlue
	case 'o':
		col = ColorDarkOrange3
	default:
		col = ui.ColorWhite
	}
	canvas.BorderStyle.Fg = col
	if down {
		rect := canvas.Inner
		for x := 1; x < rect.Dx()-1; x++ {
			for y := 1; y < rect.Dy()-1; y++ {
				canvas.SetPoint(image.Pt(x, y), col)
			}
		}
	}
	canvas.Add(image.Pt(0, 0))
	canvas.Add(image.Pt(20, 20))
	return canvas
}

func createFretsWidget(name string, frets guitarpacket.Frets) *ui.Grid {
	fretGrid := ui.NewGrid()
	fretGrid.Set(
		ui.NewRow(1.0/6,
			createParagraphWidget(name),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, createParagraphWidget("Green")),
				ui.NewRow(1.0/2, createFretWidget('g', frets.Green)),
			),
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, createParagraphWidget("Red")),
				ui.NewRow(1.0/2, createFretWidget('r', frets.Red)),
			),
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, createParagraphWidget("Yellow")),
				ui.NewRow(1.0/2, createFretWidget('y', frets.Yellow)),
			),
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, createParagraphWidget("Blue")),
				ui.NewRow(1.0/2, createFretWidget('b', frets.Blue)),
			),
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, createParagraphWidget("Orange")),
				ui.NewRow(1.0/2, createFretWidget('o', frets.Orange)),
			),
		),
	)
	return fretGrid
}

func createGuitarPacketWidget(packet guitarpacket.GuitarPacket) *ui.Grid {
	upperFrets := createFretsWidget("Upper Frets", packet.UpperFrets)
	lowerFrets := createFretsWidget("Lower Frets", packet.LowerFrets)
	container := ui.NewGrid()
	container.Set(
		ui.NewRow(1.0/2, upperFrets),
		ui.NewRow(1.0/2, lowerFrets),
	)
	return container
}

type UIWidgets struct {
	Log           string
	logTicker     int
	CurrentPacket guitarpacket.GuitarPacket
}

var uiWidgets UIWidgets

func (uiWidgets UIWidgets) GetDrawables() (drawables []ui.Drawable) {
	drawables = []ui.Drawable{
		createParagraphWidget(uiWidgets.Log),
		createGuitarPacketWidget(uiWidgets.CurrentPacket),
	}
	drawables[0].SetRect(0, 0, 150, 5)
	// drawables[1].SetRect(0, 10, 100, 60)
	return drawables
}

func uiLog(str string) {
	uiWidgets.Log = str
	uiWidgets.logTicker = 0
}

func main() {
	go UILoop()

	uiLog("Getting Joypad...")
	joypad, err := guitarjoypad.GetJoypad()
	if err != nil {
		panic(err)
	}
	guitarJoypad = joypad
	uiLog("Obtained!")
	defer guitarJoypad.Relinquish()

	uiLog("Starting Sniffer...")
	sniffer, err := guitarsniffer.Start()
	defer sniffer.Stop()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case packet := <-sniffer.Packets:
			handle_packet(&packet)
		}
	}
}

// UILoop is a goroutine which handles the UI
func UILoop() {
	fmt.Println("Entering UILoop")
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	go UIKeyboardListener()

	duration, _ := time.ParseDuration("3s")
	channel := time.Tick(duration)
	go UILogCleaner(channel)

	for {
		ui.Render(uiWidgets.GetDrawables()...)
	}
}

func UIKeyboardListener() {
	for {
		for e := range ui.PollEvents() {
			if e.Type == ui.KeyboardEvent {
				switch e.ID {
				case "q", "<C-c>":
					os.Exit(0)
				}
			}
		}
	}
}

func UILogCleaner(tick <-chan time.Time) {
	for {
		select {
		case _ = <-tick:
			uiWidgets.Log = ""
		}
	}
}

func handle_packet(packet *guitarsniffer.Packet) {
	// The packet returned when pressing the Xbox button is 31
	// bytes long, not 40, meaning that currently we're
	// ignoring that it exists but the code is there for it
	if packet.CaptureInfo.Length != 40 {
		return
	}
	uiLog(fmt.Sprintf("%d %s\n", packet.CaptureInfo.Length, hex.EncodeToString(packet.Data)))
	uiWidgets.CurrentPacket = guitarpacket.CreateGuitarPacket(packet.Data[guitarpacket.XboxHeaderLength:])
	guitarJoypad.SetValues(uiWidgets.CurrentPacket)
	guitarJoypad.Update()
}
