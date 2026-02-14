package icon

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

const (
	iconSize = 22
)

// Generate creates a 22x22 PNG icon for the given state
func Generate(state IconState) []byte {
	return generateIcon(state, false)
}

// GenerateTemplate creates a template icon (black on transparent) for macOS
func GenerateTemplate(state IconState) []byte {
	return generateIcon(state, true)
}

// GenerateAll pre-generates all four icon states
func GenerateAll() map[IconState][]byte {
	icons := make(map[IconState][]byte)
	icons[StateIdle] = Generate(StateIdle)
	icons[StateLow] = Generate(StateLow)
	icons[StateMedium] = Generate(StateMedium)
	icons[StateHigh] = Generate(StateHigh)
	return icons
}

func generateIcon(state IconState, template bool) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, iconSize, iconSize))

	// Get the color for this state
	ninjaColor := getColorForState(state, template)

	// Draw ninja face based on state
	drawNinjaFace(img, ninjaColor, state)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

func getColorForState(state IconState, template bool) color.NRGBA {
	if template {
		// Template icons are black on transparent for macOS light/dark mode
		return color.NRGBA{0, 0, 0, 255}
	}

	switch state {
	case StateIdle:
		return color.NRGBA{80, 80, 80, 255} // Gray
	case StateLow:
		return color.NRGBA{76, 175, 80, 255} // Green
	case StateMedium:
		return color.NRGBA{255, 193, 7, 255} // Amber
	case StateHigh:
		return color.NRGBA{244, 67, 54, 255} // Red
	default:
		return color.NRGBA{128, 128, 128, 255}
	}
}

func drawNinjaFace(img *image.NRGBA, c color.NRGBA, state IconState) {
	// Draw ninja head: circle centered at (11, 11) with radius 9
	drawFilledCircle(img, 11, 11, 9, c)

	// Draw mask covering lower half (creates the ninja look)
	// Mask is a filled rectangle from y=12 to bottom
	for y := 12; y < 20; y++ {
		for x := 4; x < 18; x++ {
			// Only draw if within the head circle
			dx := x - 11
			dy := y - 11
			if dx*dx+dy*dy <= 81 { // radius 9 squared
				img.Set(x, y, c)
			}
		}
	}

	// Cut out transparent area for the face opening (forehead/eyes area)
	transparent := color.NRGBA{0, 0, 0, 0}
	// Clear a horizontal band for the eyes (y = 8-11)
	for y := 8; y <= 11; y++ {
		for x := 5; x < 17; x++ {
			dx := x - 11
			dy := y - 11
			if dx*dx+dy*dy <= 81 { // Only within head circle
				img.Set(x, y, transparent)
			}
		}
	}

	// Draw different eye expressions based on state
	switch state {
	case StateIdle:
		// Chill/relaxed: closed eyes (horizontal lines)
		drawHorizontalLine(img, 6, 10, 4, c)  // Left eye
		drawHorizontalLine(img, 12, 10, 4, c) // Right eye

	case StateLow:
		// Alert: eyes dilate (small dots)
		drawFilledCircle(img, 8, 10, 1, c)  // Left eye
		drawFilledCircle(img, 14, 10, 1, c) // Right eye

	case StateMedium:
		// Ready: normal eyes + headband
		drawFilledCircle(img, 8, 10, 1, c)  // Left eye
		drawFilledCircle(img, 14, 10, 1, c) // Right eye
		// Headband across forehead
		drawHorizontalLine(img, 5, 6, 12, c)
		drawHorizontalLine(img, 5, 7, 12, c) // Make it thicker

	case StateHigh:
		// Angry: sharp angled eyes + headband
		// Left eye: angled up-right
		drawAngledLine(img, 6, 11, 10, 9, c)
		// Right eye: angled up-left
		drawAngledLine(img, 12, 9, 16, 11, c)
		// Thick headband
		drawHorizontalLine(img, 5, 6, 12, c)
		drawHorizontalLine(img, 5, 7, 12, c)
	}
}

func drawHorizontalLine(img *image.NRGBA, x, y, length int, c color.NRGBA) {
	for i := 0; i < length; i++ {
		if x+i >= 0 && x+i < iconSize && y >= 0 && y < iconSize {
			img.Set(x+i, y, c)
		}
	}
}

func drawAngledLine(img *image.NRGBA, x1, y1, x2, y2 int, c color.NRGBA) {
	// Bresenham's line algorithm for drawing angled lines
	dx := x2 - x1
	dy := y2 - y1
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := dx - dy
	x, y := x1, y1

	for {
		if x >= 0 && x < iconSize && y >= 0 && y < iconSize {
			img.Set(x, y, c)
		}

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func drawFilledCircle(img *image.NRGBA, cx, cy, r int, c color.NRGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r*r {
				if x >= 0 && x < iconSize && y >= 0 && y < iconSize {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawFilledTriangle(img *image.NRGBA, x1, y1, x2, y2, x3, y3 int, c color.NRGBA) {
	// Simple scanline fill algorithm
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	for y := minY; y <= maxY; y++ {
		// Find intersections with triangle edges
		intersections := make([]int, 0, 2)

		// Check edge 1-2
		if ix, ok := lineIntersectY(x1, y1, x2, y2, y); ok {
			intersections = append(intersections, ix)
		}
		// Check edge 2-3
		if ix, ok := lineIntersectY(x2, y2, x3, y3, y); ok {
			intersections = append(intersections, ix)
		}
		// Check edge 3-1
		if ix, ok := lineIntersectY(x3, y3, x1, y1, y); ok {
			intersections = append(intersections, ix)
		}

		if len(intersections) >= 2 {
			xMin := min(intersections[0], intersections[1])
			xMax := max(intersections[0], intersections[1])
			for x := xMin; x <= xMax; x++ {
				if x >= 0 && x < iconSize && y >= 0 && y < iconSize {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func lineIntersectY(x1, y1, x2, y2, y int) (int, bool) {
	if y1 == y2 {
		// Horizontal line
		if y == y1 {
			return x1, true
		}
		return 0, false
	}

	if (y < y1 && y < y2) || (y > y1 && y > y2) {
		return 0, false
	}

	// Linear interpolation
	t := float64(y-y1) / float64(y2-y1)
	x := float64(x1) + t*float64(x2-x1)
	return int(x + 0.5), true
}

