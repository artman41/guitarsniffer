package guitarjoypad

import (
	"github.com/artman41/guitarsniffer/guitarpacket"
	"github.com/artman41/vjoy"
)

const (
	minJoyID = 1
	maxJoyID = 16
)

// Button IDs
const (
	upperGreen uint = iota
	upperRed
	upperYellow
	upperBlue
	upperOrange
	lowerGreen
	lowerRed
	lowerYellow
	lowerBlue
	lowerOrange
	dpadUp
	dpadDown
	dpadLeft
	dpadRight
	buttonMenu
	buttonOptions
)

// Axis IDs
const (
	slider = vjoy.AxisX
	whammy = vjoy.AxisY
	tilt   = vjoy.AxisZ
)

// GuitarJoypad is a Container for the JoypadDevice
// with utility functions baked in to retrieve the
// specific Buttons & Axes
type GuitarJoypad struct {
	joypad *vjoy.Device
	rID    uint
}

// UpperGreen retrieves the Upper Green Fret
func (guitarJoypad GuitarJoypad) UpperGreen() *vjoy.Button {
	return guitarJoypad.joypad.Button(upperGreen)
}

// UpperRed retrieves the Upper Red Fret
func (guitarJoypad GuitarJoypad) UpperRed() *vjoy.Button {
	return guitarJoypad.joypad.Button(upperRed)
}

// UpperYellow retrieves the Upper Yellow Fret
func (guitarJoypad GuitarJoypad) UpperYellow() *vjoy.Button {
	return guitarJoypad.joypad.Button(upperYellow)
}

// UpperBlue retrieves the Upper Blue Fret
func (guitarJoypad GuitarJoypad) UpperBlue() *vjoy.Button {
	return guitarJoypad.joypad.Button(upperBlue)
}

// UpperOrange retrieves the Upper Orange Fret
func (guitarJoypad GuitarJoypad) UpperOrange() *vjoy.Button {
	return guitarJoypad.joypad.Button(upperOrange)
}

// LowerGreen retrieves the Lower Green Fret
func (guitarJoypad GuitarJoypad) LowerGreen() *vjoy.Button {
	return guitarJoypad.joypad.Button(lowerGreen)
}

// LowerRed retrieves the Lower Red Fret
func (guitarJoypad GuitarJoypad) LowerRed() *vjoy.Button {
	return guitarJoypad.joypad.Button(lowerRed)
}

// LowerYellow retrieves the Lower Yellow Fret
func (guitarJoypad GuitarJoypad) LowerYellow() *vjoy.Button {
	return guitarJoypad.joypad.Button(lowerYellow)
}

// LowerBlue retrieves the Lower Blue Fret
func (guitarJoypad GuitarJoypad) LowerBlue() *vjoy.Button {
	return guitarJoypad.joypad.Button(lowerBlue)
}

// LowerOrange retrieves the Lower Orange Fret
func (guitarJoypad GuitarJoypad) LowerOrange() *vjoy.Button {
	return guitarJoypad.joypad.Button(lowerOrange)
}

// DpadUp retrieves the Upper Dpad button
func (guitarJoypad GuitarJoypad) DpadUp() *vjoy.Button {
	return guitarJoypad.joypad.Button(dpadUp)
}

// DpadDown retrieves the Down Dpad button
func (guitarJoypad GuitarJoypad) DpadDown() *vjoy.Button {
	return guitarJoypad.joypad.Button(dpadDown)
}

// DpadLeft retrieves the Left Dpad button
func (guitarJoypad GuitarJoypad) DpadLeft() *vjoy.Button {
	return guitarJoypad.joypad.Button(dpadLeft)
}

// DpadRight retrieves the Right Dpad button
func (guitarJoypad GuitarJoypad) DpadRight() *vjoy.Button {
	return guitarJoypad.joypad.Button(dpadRight)
}

// ButtonMenu retrieves the Menu button
func (guitarJoypad GuitarJoypad) ButtonMenu() *vjoy.Button {
	return guitarJoypad.joypad.Button(buttonMenu)
}

// ButtonOptions retrieves the Options button
func (guitarJoypad GuitarJoypad) ButtonOptions() *vjoy.Button {
	return guitarJoypad.joypad.Button(buttonOptions)
}

// Slider retrieves the Slider Axis
func (guitarJoypad GuitarJoypad) Slider() *vjoy.Axis {
	return guitarJoypad.joypad.Axis(slider)
}

// Tilt retrieves the Tilt Axis
func (guitarJoypad GuitarJoypad) Tilt() *vjoy.Axis {
	return guitarJoypad.joypad.Axis(tilt)
}

// Whammy retrieves the Whammy Axis
func (guitarJoypad GuitarJoypad) Whammy() *vjoy.Axis {
	return guitarJoypad.joypad.Axis(whammy)
}

func (guitarJoypad GuitarJoypad) SetUpperFretValues(frets guitarpacket.Frets) {
	guitarJoypad.UpperGreen().Set(frets.Green)
	guitarJoypad.UpperRed().Set(frets.Red)
	guitarJoypad.UpperYellow().Set(frets.Yellow)
	guitarJoypad.UpperBlue().Set(frets.Blue)
	guitarJoypad.UpperOrange().Set(frets.Orange)
}

func (guitarJoypad GuitarJoypad) SetLowerFretValues(frets guitarpacket.Frets) {
	guitarJoypad.LowerGreen().Set(frets.Green)
	guitarJoypad.LowerRed().Set(frets.Red)
	guitarJoypad.LowerYellow().Set(frets.Yellow)
	guitarJoypad.LowerBlue().Set(frets.Blue)
	guitarJoypad.LowerOrange().Set(frets.Orange)
}

func (guitarJoypad GuitarJoypad) SetDpadValues(dpad guitarpacket.Dpad) {
	guitarJoypad.DpadUp().Set(dpad.Up)
	guitarJoypad.DpadDown().Set(dpad.Down)
	guitarJoypad.DpadLeft().Set(dpad.Left)
	guitarJoypad.DpadRight().Set(dpad.Right)
}

const maxFloat int = 0x7fff

func convertByte(b byte) int {
	fraction := float32(b) / float32(0xFF)
	return int(fraction * float32(maxFloat))
}

func (guitarJoypad GuitarJoypad) SetAxesValues(axes guitarpacket.Axes) {
	sliderVal := axes.Slider / 16
	var fixedSliderVal float32
	if sliderVal == 0 {
		fixedSliderVal = 0
	} else {
		sliderFraction := float32(sliderVal) / 4
		fixedSliderVal = sliderFraction * float32(0xFF)
	}
	guitarJoypad.Slider().Setc(convertByte(byte(fixedSliderVal)))
	guitarJoypad.Whammy().Setc(convertByte(axes.Whammy))
	guitarJoypad.Tilt().Setc(convertByte(axes.Tilt))
}

func (guitarJoypad GuitarJoypad) SetButtonValues(buttons guitarpacket.Buttons) {
	guitarJoypad.ButtonMenu().Set(buttons.Menu)
	guitarJoypad.ButtonOptions().Set(buttons.Options)
}

func (guitarJoypad GuitarJoypad) SetValues(guitarPacket guitarpacket.GuitarPacket) {
	guitarJoypad.SetUpperFretValues(guitarPacket.UpperFrets)
	guitarJoypad.SetLowerFretValues(guitarPacket.LowerFrets)
	guitarJoypad.SetDpadValues(guitarPacket.Dpad)
	guitarJoypad.SetAxesValues(guitarPacket.Axes)
	guitarJoypad.SetButtonValues(guitarPacket.Buttons)
}

// Update the vJoyDevice with the set
// Button & Axis values
func (guitarJoypad GuitarJoypad) Update() error {
	return guitarJoypad.joypad.Update()
}

// Reset centers all Axes & resets all Buttons
func (guitarJoypad GuitarJoypad) Reset() {
	guitarJoypad.joypad.Reset()
}

// Relinquish closes the joypad device
func (guitarJoypad GuitarJoypad) Relinquish() {
	guitarJoypad.joypad.Relinquish()
}

// GetVirtualID returns the rID assigned by vJoy
func (guitarJoypad GuitarJoypad) GetVirtualID() uint {
	return guitarJoypad.rID
}

// GetJoypad attempts to obtain a free Joypad
// with a Virtual Device ID between 1 and 16
func GetJoypad() (*GuitarJoypad, error) {
	if !vjoy.Available() {
		return nil, ErrUnavailable
	}
	dev, rID, err := accquireJoypad()
	if err != nil {
		return nil, err
	}
	return &GuitarJoypad{
		joypad: dev,
		rID:    rID,
	}, nil
}

func accquireJoypad() (dev *vjoy.Device, rID uint, err error) {
	var currentID uint = minJoyID
	dev, err = vjoy.Acquire(currentID)
	if err != nil {
		currentID++
	}
	for err == vjoy.ErrDeviceAlreadyOwned && currentID <= maxJoyID {
		dev, err = vjoy.Acquire(currentID)
		currentID++
	}
	if err != nil {
		return nil, 0, err
	}

	return dev, currentID, nil
}
