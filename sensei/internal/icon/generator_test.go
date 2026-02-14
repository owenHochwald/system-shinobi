package icon

import (
	"bytes"
	"image/png"
	"testing"
)

func TestGenerateReturnsPNG(t *testing.T) {
	data := Generate(StateIdle)
	if len(data) == 0 {
		t.Fatal("Generate returned empty data")
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 22 || bounds.Dy() != 22 {
		t.Errorf("Expected 22x22 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateAllFourStates(t *testing.T) {
	icons := GenerateAll()

	expectedStates := []IconState{StateIdle, StateLow, StateMedium, StateHigh}
	if len(icons) != len(expectedStates) {
		t.Fatalf("Expected %d icons, got %d", len(expectedStates), len(icons))
	}

	for _, state := range expectedStates {
		data, ok := icons[state]
		if !ok {
			t.Errorf("Missing icon for state %v", state)
			continue
		}

		if len(data) == 0 {
			t.Errorf("Empty data for state %v", state)
			continue
		}

		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Errorf("Failed to decode PNG for state %v: %v", state, err)
			continue
		}

		bounds := img.Bounds()
		if bounds.Dx() != 22 || bounds.Dy() != 22 {
			t.Errorf("State %v: expected 22x22 image, got %dx%d", state, bounds.Dx(), bounds.Dy())
		}
	}
}

func TestDistinctIcons(t *testing.T) {
	icons := GenerateAll()

	// Compare each pair of icons to ensure they're different
	states := []IconState{StateIdle, StateLow, StateMedium, StateHigh}
	for i := 0; i < len(states); i++ {
		for j := i + 1; j < len(states); j++ {
			data1 := icons[states[i]]
			data2 := icons[states[j]]

			if bytes.Equal(data1, data2) {
				t.Errorf("Icons for %v and %v are identical", states[i], states[j])
			}
		}
	}
}
