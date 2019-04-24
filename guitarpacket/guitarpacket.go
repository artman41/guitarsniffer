package guitarpacket

import (
	"encoding/hex"
	"fmt"
)

const XboxHeaderLength = 22

// Frets
const (
	FretGreen  = 0x1
	FretRed    = 0x2
	FretYellow = 0x4
	FretBlue   = 0x8
	FretOrange = 0x10
)

// Buttons
const (
	ButtonXbox    = 0x1
	ButtonMenu    = 0x4
	ButtonOptions = 0x8
)

const (
	DpadDown  = 0x1
	DpadUp    = 0x2
	DpadLeft  = 0x4
	DpadRight = 0x8
)

// Packet Pieces
const (
	PosButtons   = 8
	PosDpad      = 9
	PosTilt      = 10
	PosWhammy    = 11
	PosSlider    = 12
	PosUpperFret = 13
	PosLowerFret = 14
)

// CreateGuitarPacket returns a GuitarPacket struct
// filled with the values of the given packet
//
// Note: the function assumes that the given packet
// has already had the XboxHeader removed from it
func CreateGuitarPacket(packet []byte) GuitarPacket {
	fmt.Printf("(%d) %s\n", len(packet), hex.EncodeToString(packet))
	upperFrets := getFrets(packet[PosUpperFret])
	lowerFrets := getFrets(packet[PosLowerFret])
	dpad := getDpad(packet[PosDpad])
	buttons := getButtons(packet[PosButtons])
	return GuitarPacket{
		UpperFrets: upperFrets,
		LowerFrets: lowerFrets,
		Dpad:       dpad,
		Slider:     int(packet[PosSlider] / 16),
		Buttons:    buttons,
		Whammy:     packet[PosWhammy],
		Tilt:       packet[PosTilt],
	}
}

func getFrets(fretBitMask byte) Frets {
	return Frets{
		Green:  fretBitMask&FretGreen != 0,
		Red:    fretBitMask&FretRed != 0,
		Yellow: fretBitMask&FretYellow != 0,
		Blue:   fretBitMask&FretBlue != 0,
		Orange: fretBitMask&FretOrange != 0,
	}
}

func getButtons(buttonBitMask byte) Buttons {
	return Buttons{
		Menu:    (buttonBitMask&ButtonMenu != 0),
		Options: (buttonBitMask&ButtonOptions != 0),
		Xbox:    (buttonBitMask&ButtonXbox != 0),
	}
}

func getDpad(dpadBitMask byte) Dpad {
	return Dpad{
		Up:    (dpadBitMask&DpadUp != 0),
		Down:  (dpadBitMask&DpadDown != 0),
		Left:  (dpadBitMask&DpadLeft != 0),
		Right: (dpadBitMask&DpadRight != 0),
	}
}

type Frets struct {
	Green, Red, Yellow, Blue, Orange bool
}

type Buttons struct {
	Menu, Options, Xbox bool
}

type Dpad struct {
	Up, Down, Left, Right bool
}

type GuitarPacket struct {
	UpperFrets, LowerFrets Frets
	Dpad                   Dpad
	Slider                 int
	Buttons                Buttons
	Whammy                 byte
	Tilt                   byte
}
