package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
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
	
	// For mouse/trackpad interaction
	isDragging    bool
	lastMousePosX float64
	lastMousePosY float64
	canvasObj     *canvas.Raster
	hoverX        float64
	hoverY        float64
	
	// Mode flags
	rotateMode bool
	panMode    bool
	rKeyPressed bool
}

// NewVisualizer creates a new 3D visualizer
func NewVisualizer(space *Space3D) *Visualizer {
	vis := &Visualizer{
		space:      space,
		pointSize:  10,
		width:      800,
		height:     600,
		scale:      50,
		xRotation:  0,
		yRotation:  0,
		zRotation:  0,
		xOffset:    0,
		yOffset:    0,
		rotateMode: true,
		panMode:    false,
		rKeyPressed: false,
		hoverX:     0,
		hoverY:     0,
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

// Custom MouseDown event handler
func (v *Visualizer) handleMouseDown(ev *desktop.MouseEvent) {
	v.isDragging = true
	v.lastMousePosX = float64(ev.Position.X)
	v.lastMousePosY = float64(ev.Position.Y)

	// R key state is tracked by keyboard handlers
	if v.rKeyPressed {
		// R key is pressed - force rotation mode
		v.rotateMode = true
		v.panMode = false
	} else if ev.Button == desktop.MouseButtonSecondary || ev.Modifier == fyne.KeyModifierAlt {
		// Right mouse button or Option/Alt + click for panning
		v.rotateMode = false
		v.panMode = true
	} else {
		// Default is left-click for rotation
		v.rotateMode = true
		v.panMode = false
	}
}

// Custom MouseUp event handler
func (v *Visualizer) handleMouseUp() {
	v.isDragging = false
}

// Custom MouseMoved event handler
func (v *Visualizer) handleMouseMove(ev *desktop.MouseEvent) {
	// Always update hover position
	v.hoverX = float64(ev.Position.X)
	v.hoverY = float64(ev.Position.Y)
	
	// Refresh to update hover effects
	v.canvasObj.Refresh()
	
	if !v.isDragging {
		return
	}

	// Calculate the delta movement
	deltaX := float64(ev.Position.X) - v.lastMousePosX
	deltaY := float64(ev.Position.Y) - v.lastMousePosY

	// Update the last position
	v.lastMousePosX = float64(ev.Position.X)
	v.lastMousePosY = float64(ev.Position.Y)
	
	// Handle based on mode or R key
	if v.rKeyPressed || v.rotateMode {
		// Rotation - adjust the rotation based on mouse movement
		sensitivity := 0.01
		v.yRotation += deltaX * sensitivity
		v.xRotation += deltaY * sensitivity
	} else if v.panMode {
		// Panning - adjust the offset based on mouse movement
		v.xOffset += deltaX
		v.yOffset += deltaY
	}

	// Refresh the canvas
	v.canvasObj.Refresh()
}

// Custom scroll event handler for zooming or rotating
func (v *Visualizer) handleScroll(ev *fyne.ScrollEvent) {
	// Check if R key is pressed for rotation
	if v.rKeyPressed {
		// R key + scroll for rotation
		rotationSpeed := 0.1
		
		// Vertical scroll (DY) changes X rotation (up/down)
		if ev.Scrolled.DY != 0 {
			v.xRotation += float64(ev.Scrolled.DY) * rotationSpeed
		}
		
		// Horizontal scroll (DX) changes Y rotation (left/right)
		if ev.Scrolled.DX != 0 {
			v.yRotation += float64(ev.Scrolled.DX) * rotationSpeed
		}
	} else {
		// Normal zoom mode
		zoomFactor := 1.1
		
		if ev.Scrolled.DY < 0 {
			// Zoom out
			v.scale /= zoomFactor
		} else {
			// Zoom in
			v.scale *= zoomFactor
		}
		
		// Enforce min/max scale values
		if v.scale < 5 {
			v.scale = 5
		} else if v.scale > 500 {
			v.scale = 500
		}
	}
	
	// Refresh the canvas
	v.canvasObj.Refresh()
}

// Run starts the visualizer
func (v *Visualizer) Run() {
	v.app = app.New()
	v.window = v.app.NewWindow("3D Points Visualizer")
	v.window.Resize(fyne.NewSize(1000, 800))

	// Create a canvas to draw on
	v.canvasObj = canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))

		// Update width and height based on current canvas size
		v.width = float32(w)
		v.height = float32(h)

		// Draw a light gray background
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.Set(x, y, color.RGBA{240, 240, 240, 255})
			}
		}

		// Draw 3D grid across all three planes
		gridSize := 5
		gridStep := 1.0
		
		// Colors for the different plane grids
		xzGridColor := color.RGBA{180, 180, 180, 160}  // Floor - light gray
		xyGridColor := color.RGBA{180, 180, 220, 120}  // Side wall - light blue tint
		yzGridColor := color.RGBA{220, 180, 180, 120}  // Side wall - light red tint
		
		// 1. XZ Plane (Floor) grid
		for i := -gridSize; i <= gridSize; i++ {
			for j := -gridSize; j <= gridSize; j++ {
				p1 := NewPoint3D(float64(i)*gridStep, 0, float64(j)*gridStep)
				p2 := NewPoint3D(float64(i+1)*gridStep, 0, float64(j)*gridStep)
				p3 := NewPoint3D(float64(i)*gridStep, 0, float64(j+1)*gridStep)
				
				x1, y1 := v.project3DTo2D(p1)
				x2, y2 := v.project3DTo2D(p2)
				x3, y3 := v.project3DTo2D(p3)
				
				// Only draw if within screen bounds
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x2), int(y2), w, h) {
					drawLine(img, int(x1), int(y1), int(x2), int(y2), xzGridColor)
				}
				
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x3), int(y3), w, h) {
					drawLine(img, int(x1), int(y1), int(x3), int(y3), xzGridColor)
				}
			}
		}
		
		// 2. XY Plane (Vertical wall) grid
		for i := -gridSize; i <= gridSize; i++ {
			for j := -gridSize; j <= gridSize; j++ {
				p1 := NewPoint3D(float64(i)*gridStep, float64(j)*gridStep, 0)
				p2 := NewPoint3D(float64(i+1)*gridStep, float64(j)*gridStep, 0)
				p3 := NewPoint3D(float64(i)*gridStep, float64(j+1)*gridStep, 0)
				
				x1, y1 := v.project3DTo2D(p1)
				x2, y2 := v.project3DTo2D(p2)
				x3, y3 := v.project3DTo2D(p3)
				
				// Only draw if within screen bounds
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x2), int(y2), w, h) {
					drawLine(img, int(x1), int(y1), int(x2), int(y2), xyGridColor)
				}
				
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x3), int(y3), w, h) {
					drawLine(img, int(x1), int(y1), int(x3), int(y3), xyGridColor)
				}
			}
		}
		
		// 3. YZ Plane (Vertical wall) grid
		for i := -gridSize; i <= gridSize; i++ {
			for j := -gridSize; j <= gridSize; j++ {
				p1 := NewPoint3D(0, float64(i)*gridStep, float64(j)*gridStep)
				p2 := NewPoint3D(0, float64(i+1)*gridStep, float64(j)*gridStep)
				p3 := NewPoint3D(0, float64(i)*gridStep, float64(j+1)*gridStep)
				
				x1, y1 := v.project3DTo2D(p1)
				x2, y2 := v.project3DTo2D(p2)
				x3, y3 := v.project3DTo2D(p3)
				
				// Only draw if within screen bounds
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x2), int(y2), w, h) {
					drawLine(img, int(x1), int(y1), int(x2), int(y2), yzGridColor)
				}
				
				if isVisible(int(x1), int(y1), w, h) && isVisible(int(x3), int(y3), w, h) {
					drawLine(img, int(x1), int(y1), int(x3), int(y3), yzGridColor)
				}
			}
		}

		// Draw prominent coordinate axes
		origin := NewPoint3D(0, 0, 0)
		
		// Make axes longer and add axis labels
		axisLength := 2.0
		xAxis := NewPoint3D(axisLength, 0, 0)
		yAxis := NewPoint3D(0, axisLength, 0)
		zAxis := NewPoint3D(0, 0, axisLength)
		
		// Get 2D coordinates
		ox, oy := v.project3DTo2D(origin)
		xx, xy := v.project3DTo2D(xAxis)
		yx, yy := v.project3DTo2D(yAxis)
		zx, zy := v.project3DTo2D(zAxis)
		
		// Draw thicker axes with more vibrant colors
		axisThickness := 3
		
		// X-axis (bright red)
		drawThickLine(img, int(ox), int(oy), int(xx), int(xy), color.RGBA{255, 50, 50, 255}, axisThickness)
		// Y-axis (bright green)
		drawThickLine(img, int(ox), int(oy), int(yx), int(yy), color.RGBA{50, 255, 50, 255}, axisThickness)
		// Z-axis (bright blue)
		drawThickLine(img, int(ox), int(oy), int(zx), int(zy), color.RGBA{50, 50, 255, 255}, axisThickness)
		
		// Add axis labels
		drawString(img, "X", int(xx)+5, int(xy)-5, color.RGBA{255, 0, 0, 255})
		drawString(img, "Y", int(yx)+5, int(yy)-5, color.RGBA{0, 255, 0, 255})
		drawString(img, "Z", int(zx)+5, int(zy)-5, color.RGBA{0, 0, 255, 255})

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

			// Check if mouse is near the point
			mouseX, mouseY := int(v.hoverX), int(v.hoverY)
			pointX, pointY := int(screenX), int(screenY)
			
			// Simple box-based hover detection with large margins
			boxSize := size * 5
			if mouseX >= pointX-boxSize && mouseX <= pointX+boxSize &&
			   mouseY >= pointY-boxSize && mouseY <= pointY+boxSize {
				coordStr := formatCoord(point)
				drawString(img, coordStr, int(screenX)+size+5, int(screenY)-5, color.RGBA{50, 50, 50, 255})
			}
		}


		return img
	})
	
	// Instructions card
	instructionsCard := widget.NewCard("", "Controls",
		widget.NewLabel("• Rotate: Left-click + drag\n• Rotate (alternate): Hold R key + scroll wheel\n• Pan: Right-click + drag (or Option/Alt + drag)\n• Zoom: Scroll wheel or pinch gesture"))

	// Reset button
	resetBtn := widget.NewButton("Reset View", func() {
		v.xRotation = 0
		v.yRotation = 0
		v.zRotation = 0
		v.scale = 50
		v.xOffset = 0
		v.yOffset = 0
		v.canvasObj.Refresh()
	})
	
	// Upload CSV button
	uploadBtn := widget.NewButton("Upload CSV", func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, v.window)
				return
			}
			if reader == nil {
				return // User cancelled
			}
			defer reader.Close()
			
			// Create a temporary file
			tempFile, err := os.CreateTemp("", "upload-*.csv")
			if err != nil {
				dialog.ShowError(err, v.window)
				return
			}
			defer tempFile.Close()
			
			// Copy the uploaded file to the temp file
			_, err = io.Copy(tempFile, reader)
			if err != nil {
				dialog.ShowError(err, v.window)
				return
			}
			
			// Create new space and load points
			newSpace := NewSpace3D()
			err = newSpace.LoadPointsFromCSV(tempFile.Name())
			if err != nil {
				dialog.ShowError(err, v.window)
				return
			}
			
			// Update visualizer with new points
			v.space = newSpace
			
			// Reset view for better visualization
			v.xRotation = 0
			v.yRotation = 0
			v.zRotation = 0
			v.scale = 50
			v.xOffset = 0
			v.yOffset = 0
			v.canvasObj.Refresh()
		}, v.window)
		
		// Set filter for CSV files
		openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		openDialog.Show()
	})
	
	// Layout
	controls := container.New(layout.NewVBoxLayout(),
		instructionsCard,
		uploadBtn,
		resetBtn,
	)

	// Create a custom canvas wrapper for mouse interaction
	canvasWrapper := newCanvasWrapper(v.canvasObj, v)
	
	// Add keyboard event handler for key press
	v.window.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == "R" || ke.Name == "r" {
			v.rKeyPressed = true
		}
	})
	
	// Add a goroutine to simulate key releases since Fyne doesn't provide direct access
	go func() {
		for {
			if v.rKeyPressed {
				// Simulate a key release after a short delay
				time.Sleep(300 * time.Millisecond)
				v.rKeyPressed = false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Split layout with controls on the right
	split := container.NewHSplit(
		canvasWrapper,
		container.NewPadded(controls),
	)
	split.Offset = 0.8 // 80% of space for canvas, 20% for controls

	v.window.SetContent(split)
	v.window.ShowAndRun()
}

// canvasWrapper is a custom widget that wraps a canvas and handles mouse events
type canvasWrapper struct {
	widget.BaseWidget
	canvas fyne.CanvasObject
	vis    *Visualizer
}

// newCanvasWrapper creates a new canvas wrapper
func newCanvasWrapper(canvas fyne.CanvasObject, vis *Visualizer) *canvasWrapper {
	w := &canvasWrapper{
		canvas: canvas,
		vis:    vis,
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer implements the widget interface
func (c *canvasWrapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.canvas)
}

// MouseDown implements desktop.Mouseable
func (c *canvasWrapper) MouseDown(ev *desktop.MouseEvent) {
	c.vis.handleMouseDown(ev)
}

// MouseUp implements desktop.Mouseable
func (c *canvasWrapper) MouseUp(*desktop.MouseEvent) {
	c.vis.handleMouseUp()
}

// MouseMoved implements desktop.Mouseable
func (c *canvasWrapper) MouseMoved(ev *desktop.MouseEvent) {
	c.vis.handleMouseMove(ev)
}

// MouseIn implements desktop.Hoverable
func (c *canvasWrapper) MouseIn(*desktop.MouseEvent) {}

// MouseOut implements desktop.Hoverable
func (c *canvasWrapper) MouseOut() {}

// Scrolled implements fyne.Scrollable
func (c *canvasWrapper) Scrolled(ev *fyne.ScrollEvent) {
	c.vis.handleScroll(ev)
}

// Ensure canvasWrapper implements necessary interfaces
var _ desktop.Mouseable = (*canvasWrapper)(nil)
var _ desktop.Hoverable = (*canvasWrapper)(nil)
var _ fyne.Scrollable = (*canvasWrapper)(nil)

// Helper function to check if a point is visible on screen
func isVisible(x, y, width, height int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
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

// Helper function for drawing thicker lines
func drawThickLine(img *image.RGBA, x1, y1, x2, y2 int, clr color.RGBA, thickness int) {
	// Draw the main line
	drawLine(img, x1, y1, x2, y2, clr)
	
	// Draw additional lines around the main line for thickness
	halfThick := thickness / 2
	for dy := -halfThick; dy <= halfThick; dy++ {
		for dx := -halfThick; dx <= halfThick; dx++ {
			if dx*dx + dy*dy <= halfThick*halfThick {
				drawLine(img, x1+dx, y1+dy, x2+dx, y2+dy, clr)
			}
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
		'X': {
			"#   #",
			" # # ",
			"  #  ",
			" # # ",
			"#   #",
			"     ",
		},
		'Y': {
			"#   #",
			" # # ",
			"  #  ",
			"  #  ",
			"  #  ",
			"     ",
		},
		'Z': {
			"#####",
			"   # ",
			"  #  ",
			" #   ",
			"#####",
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}