package main

import "fmt"
import "log"
import "math"
import "errors"
import "strconv"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/etxt"
import "github.com/tinne26/fonts/liberation/lbrtmono"

var ErrExit error = errors.New("esc to exit the program")
var (
	BackColor  = color.RGBA{0, 24, 28, 255}
	MainColor  = color.RGBA{239, 241, 197, 255}
	FocusColor = color.RGBA{234, 82, 111, 255}
)

type View struct {
	text *etxt.Renderer
	gamepadIds []ebiten.GamepadID
	gamepadName string
	numButtons int
	pressedButtons []ebiten.GamepadButton
	axisValues []float64

	// internal variables
	lastCanvasWidth float64
	lastCanvasHeight float64
	lastDisplayScale float64
	fsKeyPressed bool
}

func (self *View) Layout(int, int) (int, int) { panic("ebitengine >= v2.5.0") }
func (self *View) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	// get display scale and refresh it if necessary
	scale := ebiten.DeviceScaleFactor()
	if scale != self.lastDisplayScale {
		self.text.SetScale(scale)
	}

	// compute canvas dimensions and refresh if necessary
	canvasWidth  := math.Ceil(logicWinWidth*scale)
	canvasHeight := math.Ceil(logicWinHeight*scale)
	if canvasWidth != self.lastCanvasWidth || canvasHeight != self.lastCanvasHeight {
		self.lastCanvasWidth = canvasWidth
		self.lastCanvasHeight = canvasHeight
	}
	return canvasWidth, canvasHeight
}

func (self *View) Update() error {
	// detect game being closed
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrExit
	}

	// handle fullscreening
	fsKeyPressed := ebiten.IsKeyPressed(ebiten.KeyF)
	if fsKeyPressed != self.fsKeyPressed {
		if fsKeyPressed {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
		self.fsKeyPressed = fsKeyPressed
	}

	// get gamepad ids and read buttons and axes
	self.gamepadIds = self.gamepadIds[: 0]
	self.gamepadIds = ebiten.AppendGamepadIDs(self.gamepadIds)

	if len(self.gamepadIds) > 0 {
		gamepadId := self.gamepadIds[0]
		self.gamepadName = ebiten.GamepadName(gamepadId)
		self.numButtons  = ebiten.GamepadButtonCount(gamepadId)
		self.pressedButtons = self.pressedButtons[ : 0]
		for i := 0; i < self.numButtons; i++ {
			button := ebiten.GamepadButton(i)
			if ebiten.IsGamepadButtonPressed(gamepadId, button) {
				self.pressedButtons = append(self.pressedButtons, button)
			}
		}

		axisCount := ebiten.GamepadAxisCount(gamepadId)
		if len(self.axisValues) < axisCount {
			if cap(self.axisValues) >= axisCount {
				self.axisValues = self.axisValues[ : axisCount]
			} else {
				self.axisValues = make([]float64, axisCount)
			}
		}
		for i := 0; i < axisCount; i++ {
			self.axisValues[i] = ebiten.GamepadAxisValue(gamepadId, i)
		}
		self.axisValues = self.axisValues[ : axisCount]
	}

	return nil
}

func (self *View) Draw(canvas *ebiten.Image) {
	canvas.Fill(BackColor)

	self.text.SetColor(MainColor)
	lineAdvance := self.text.Utils().GetLineHeight()
	baseX, tabX := int(lineAdvance), int(lineAdvance*2.0)
	y := lineAdvance*1.6
	if len(self.gamepadIds) == 0 {
		self.text.Draw(canvas, "No gamepads detected\nDid you plug it in already?", baseX, int(y))
	} else { // len(self.gamepadIds) > 0
		// "Detected N gamepads"
		self.text.Draw(canvas, "Detected " + strconv.Itoa(len(self.gamepadIds)) + " gamepad(s)", baseX, int(y))
		y += lineAdvance

		// "Monitoring XXXXX"
		pre := "Monitoring "
		self.text.Draw(canvas, pre, baseX, int(y))
		
		self.text.SetColor(FocusColor)
		x := lineAdvance + self.text.Measure(pre).Width().ToFloat64()
		self.text.Draw(canvas, self.gamepadName, int(x), int(y))
		y += lineAdvance*1.5
		self.text.SetColor(MainColor)

		// "Axis values:"
		self.text.Draw(canvas, "Axis values:", baseX, int(y))
		y += lineAdvance
		self.text.SetColor(FocusColor)
		if len(self.axisValues) == 0 {
			self.text.Draw(canvas, "(No axes detected)", int(lineAdvance*2), int(y))
			y += lineAdvance
		} else {
			hint := []string{
				"left joystick horz", "left joystick vert",
				"right joystick horz", "right joystick vert",
				"left trigger", "right trigger",
			}
			for i, value := range self.axisValues {
				floatStr := strconv.FormatFloat(value, 'f', 2, 64)
				if floatStr == "-0.00" { floatStr = "0.00" }
				if floatStr[0] != '-' { floatStr = " " + floatStr }
				self.text.Draw(canvas, floatStr, int(lineAdvance*2), int(y))
				if len(hint) > i {
					offset := self.text.Measure(floatStr).Width().ToFloat64() + lineAdvance
					self.text.SetColor(reAlpha(MainColor, 144))
					self.text.Draw(canvas, "(" + hint[i] + ")", int(lineAdvance*2 + offset), int(y))
					self.text.SetColor(FocusColor)
				}
				y += lineAdvance
			}
		}
		self.text.SetColor(MainColor)
		y += lineAdvance*0.5

		// "Pressed buttons:"
		self.text.Draw(canvas, "Pressed buttons [0.." + strconv.Itoa(self.numButtons - 1) + "]:", baseX, int(y))
		y += lineAdvance
		self.text.SetColor(FocusColor)
		if len(self.pressedButtons) == 0 {
			self.text.Draw(canvas, "(no buttons pressed)", tabX, int(y))
			y += lineAdvance
		} else {
			var buttons string
			for i, value := range self.pressedButtons {
				buttons += "#" + strconv.Itoa(int(value))
				if i != len(self.pressedButtons) - 1 {
					buttons += ", "
				}
			}
			self.text.Draw(canvas, buttons, tabX, int(y))
			y += lineAdvance
		}
		self.text.SetColor(MainColor)
	}
}

func main() {
	renderer := etxt.NewRenderer()
	renderer.Utils().SetCache8MiB()
	renderer.SetFont(lbrtmono.Font())
	renderer.SetSize(16)
	renderer.SetAlign(etxt.Left | etxt.TopBaseline)

	fmt.Printf("Axis labels are only hints, they may not match your controller.\n")
	fmt.Printf("Press F to fullscreen, ESC to close the program.\n")

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("tinne/gamepad-quicktest")
	ebiten.SetScreenClearedEveryFrame(false)
	err := ebiten.RunGame(&View{ text: renderer })
	if err != nil && err != ErrExit { log.Fatal(err) }
}

// Rescale the given color to the given alpha.
func reAlpha(clr color.RGBA, alpha uint8) color.RGBA {
	scalingFactor := float64(alpha)/float64(clr.A)
	return color.RGBA{
		R: uint8(scalingFactor*float64(clr.R)),
		G: uint8(scalingFactor*float64(clr.G)),
		B: uint8(scalingFactor*float64(clr.B)),
		A: uint8(alpha),
	}
}
