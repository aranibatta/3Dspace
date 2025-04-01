package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Visualizer represents a 3D points visualizer
type Visualizer struct {
	space     *Space3D
	app       fyne.App
	window    fyne.Window
	pointSize float32
	width     float32
	height    float32
	xRotation float64
	yRotation float64
	zRotation float64
	scale     float64
	xOffset   float64
	yOffset   float64
}

// NewVisualizer creates a new 3D visualizer
func NewVisualizer(space *Space3D) *Visualizer {
	vis := &Visualizer{
		space:     space,
		pointSize: 10,
		width:     800,
		height:    600,
		scale:     50,
		xRotation: 0,
		yRotation: 0,
		zRotation: 0,
		xOffset:   0,
		yOffset:   0,
	}
	return vis
}

// project3DTo2D projects a 3D point onto a 2D plane with simple perspective
func (v *Visualizer) project3DTo2D(point Point3D) (float32, float32) {
	// Apply rotations (very simple rotation around axes)
	// For a real application, you'd want to use a proper 3D matrix library
	x, y, z := point.X, point.Y, point.Z

	// Apply rotations (just a simple example)
	// X rotation
	tempY := y*math.Cos(v.xRotation) - z*math.Sin(v.xRotation)
	tempZ := y*math.Sin(v.xRotation) + z*math.Cos(v.xRotation)
	y, z = tempY, tempZ

	// Y rotation
	tempX := x*math.Cos(v.yRotation) + z*math.Sin(v.yRotation)
	tempZ = -x*math.Sin(v.yRotation) + z*math.Cos(v.yRotation)
	x, z = tempX, tempZ

	// Z rotation
	tempX = x*math.Cos(v.zRotation) - y*math.Sin(v.zRotation)
	tempY = x*math.Sin(v.zRotation) + y*math.Cos(v.zRotation)
	x, y = tempX, tempY

	// Scale, center, and add offset
	// Add a simple perspective effect (farther objects appear smaller)
	perspective := float64(600) / (float64(600) + z)

	screenX := float32(float64(v.width)/2 + (x * v.scale * perspective) + v.xOffset)
	screenY := float32(float64(v.height)/2 + (y * v.scale * perspective) + v.yOffset)

	return screenX, screenY
}

// Run starts the visualizer
func (v *Visualizer) Run() {
	v.app = app.New()
	v.window = v.app.NewWindow("3D Points Visualizer")
	v.window.Resize(fyne.NewSize(1000, 800))

	// Create a canvas to draw on
	canvasObj := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))

		// Draw a light gray background
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.Set(x, y, color.RGBA{240, 240, 240, 255})
			}
		}

		// Draw axes
		drawLine(img, int(v.width/2), 0, int(v.width/2), h, color.RGBA{200, 200, 200, 255})
		drawLine(img, 0, int(v.height/2), w, int(v.height/2), color.RGBA{200, 200, 200, 255})

		// Draw points
		for _, point := range v.space.Points {
			screenX, screenY := v.project3DTo2D(point)

			// Draw a point with border
			size := int(v.pointSize)

			// Draw border (black outline)
			borderSize := size + 2
			for y := -borderSize; y <= borderSize; y++ {
				for x := -borderSize; x <= borderSize; x++ {
					if x*x+y*y <= borderSize*borderSize {
						px, py := int(screenX)+x, int(screenY)+y
						if px >= 0 && px < w && py >= 0 && py < h {
							img.Set(px, py, color.RGBA{0, 0, 0, 255})
						}
					}
				}
			}

			// Draw inner circle (blue)
			for y := -size; y <= size; y++ {
				for x := -size; x <= size; x++ {
					if x*x+y*y <= size*size {
						px, py := int(screenX)+x, int(screenY)+y
						if px >= 0 && px < w && py >= 0 && py < h {
							img.Set(px, py, color.RGBA{30, 144, 255, 255})
						}
					}
				}
			}

			// Draw coordinates
			coordStr := formatCoord(point)
			drawString(img, coordStr, int(screenX)+size+5, int(screenY)-5, color.RGBA{50, 50, 50, 255})
		}

		return img
	})

	// Controls for rotation
	xRotSlider := widget.NewSlider(-math.Pi, math.Pi)
	xRotSlider.OnChanged = func(value float64) {
		v.xRotation = value
		canvasObj.Refresh()
	}

	yRotSlider := widget.NewSlider(-math.Pi, math.Pi)
	yRotSlider.OnChanged = func(value float64) {
		v.yRotation = value
		canvasObj.Refresh()
	}

	zRotSlider := widget.NewSlider(-math.Pi, math.Pi)
	zRotSlider.OnChanged = func(value float64) {
		v.zRotation = value
		canvasObj.Refresh()
	}

	scaleSlider := widget.NewSlider(10, 200)
	scaleSlider.Value = 50
	scaleSlider.OnChanged = func(value float64) {
		v.scale = value
		canvasObj.Refresh()
	}

	// Layout
	rotationCard := widget.NewCard("Rotation Controls", "",
		container.New(layout.NewVBoxLayout(),
			container.NewPadded(container.New(layout.NewFormLayout(), widget.NewLabel("X:"), xRotSlider)),
			container.NewPadded(container.New(layout.NewFormLayout(), widget.NewLabel("Y:"), yRotSlider)),
			container.NewPadded(container.New(layout.NewFormLayout(), widget.NewLabel("Z:"), zRotSlider)),
		),
	)

	scaleCard := widget.NewCard("Display Settings", "",
		container.NewPadded(container.New(layout.NewFormLayout(), widget.NewLabel("Scale:"), scaleSlider)),
	)

	controls := container.New(layout.NewHBoxLayout(),
		container.NewPadded(rotationCard),
		container.NewPadded(scaleCard),
	)

	content := container.New(layout.NewBorderLayout(nil, controls, nil, nil),
		canvasObj, controls)

	v.window.SetContent(content)
	v.window.ShowAndRun()
}

// Helper function to draw a line using Bresenham's algorithm
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, clr color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		if x1 >= 0 && y1 >= 0 && x1 < img.Bounds().Max.X && y1 < img.Bounds().Max.Y {
			img.Set(x1, y1, clr)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// formatCoord returns a formatted string of point coordinates
func formatCoord(p Point3D) string {
	return fmt.Sprintf("(%.1f, %.1f, %.1f)", p.X, p.Y, p.Z)
}

// drawString draws a string on the image
func drawString(img *image.RGBA, s string, x, y int, clr color.RGBA) {
	// Simple font map for basic characters (very basic implementation)
	fontMap := map[rune][]string{
		'0': {
			" ### ",
			"#   #",
			"#   #",
			"#   #",
			"#   #",
			" ### ",
		},
		'1': {
			"  #  ",
			" ##  ",
			"  #  ",
			"  #  ",
			"  #  ",
			" ### ",
		},
		'2': {
			" ### ",
			"#   #",
			"   # ",
			"  #  ",
			" #   ",
			"#####",
		},
		'3': {
			" ### ",
			"#   #",
			"  ## ",
			"    #",
			"#   #",
			" ### ",
		},
		'4': {
			"   # ",
			"  ## ",
			" # # ",
			"#  # ",
			"#####",
			"   # ",
		},
		'5': {
			"#####",
			"#    ",
			"#### ",
			"    #",
			"#   #",
			" ### ",
		},
		'6': {
			" ### ",
			"#    ",
			"#### ",
			"#   #",
			"#   #",
			" ### ",
		},
		'7': {
			"#####",
			"    #",
			"   # ",
			"  #  ",
			" #   ",
			"#    ",
		},
		'8': {
			" ### ",
			"#   #",
			" ### ",
			"#   #",
			"#   #",
			" ### ",
		},
		'9': {
			" ### ",
			"#   #",
			"#   #",
			" ####",
			"    #",
			" ### ",
		},
		'.': {
			"     ",
			"     ",
			"     ",
			"     ",
			"     ",
			"  #  ",
		},
		'-': {
			"     ",
			"     ",
			"#####",
			"     ",
			"     ",
			"     ",
		},
		',': {
			"     ",
			"     ",
			"     ",
			"     ",
			"  #  ",
			" #   ",
		},
		'(': {
			"  #  ",
			" #   ",
			"#    ",
			"#    ",
			" #   ",
			"  #  ",
		},
		')': {
			"  #  ",
			"   # ",
			"    #",
			"    #",
			"   # ",
			"  #  ",
		},
		' ': {
			"     ",
			"     ",
			"     ",
			"     ",
			"     ",
			"     ",
		},
	}

	// Set default patterns for any undefined characters
	defaultPattern := []string{
		"#####",
		"#   #",
		"#   #",
		"#   #",
		"#   #",
		"#####",
	}

	// Character width and height
	charWidth := 5
	charHeight := 6
	spacing := 1

	// Draw background for better readability
	bgPadding := 2
	bgWidth := len(s)*(charWidth+spacing) + bgPadding*2
	// Calculate background dimensions

	// Draw semi-transparent background
	for by := y - bgPadding; by < y+charHeight+bgPadding; by++ {
		for bx := x - bgPadding; bx < x+bgWidth; bx++ {
			if bx >= 0 && bx < img.Bounds().Max.X && by >= 0 && by < img.Bounds().Max.Y {
				img.Set(bx, by, color.RGBA{240, 240, 240, 220})
			}
		}
	}

	// Draw each character
	for i, char := range s {
		pattern, ok := fontMap[char]
		if !ok {
			pattern = defaultPattern
		}

		for dy := 0; dy < charHeight; dy++ {
			if dy < len(pattern) {
				for dx := 0; dx < charWidth; dx++ {
					px := x + i*(charWidth+spacing) + dx
					py := y + dy

					if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
						if dx < len(pattern[dy]) && pattern[dy][dx] == '#' {
							img.Set(px, py, clr)
						}
					}
				}
			}
		}
	}
}

// drawText measures the width of text
func drawText(img *image.RGBA, p Point3D, x, y int) int {
	return len(formatCoord(p)) * 6
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
