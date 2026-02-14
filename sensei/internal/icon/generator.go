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

	// Draw ninja silhouette
	drawNinjaSilhouette(img, ninjaColor)

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

func drawNinjaSilhouette(img *image.NRGBA, c color.NRGBA) {
	// Draw a simple ninja figure:
	// - Circle head at top
	// - Triangular body
	// - Two small eye dots

	// Head: filled circle centered at (11, 6) with radius 4
	drawFilledCircle(img, 11, 6, 4, c)

	// Body: triangle from (6, 11) to (16, 11) to (11, 20)
	drawFilledTriangle(img, 6, 11, 16, 11, 11, 20, c)

	// Eyes: two small dots (use a contrasting color or white)
	eyeColor := color.NRGBA{255, 255, 255, 255} // White eyes
	if c.R > 200 && c.G > 200 && c.B > 200 {
		// If the ninja color is very light, use black eyes
		eyeColor = color.NRGBA{0, 0, 0, 255}
	}
	img.Set(9, 6, eyeColor)
	img.Set(13, 6, eyeColor)
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

