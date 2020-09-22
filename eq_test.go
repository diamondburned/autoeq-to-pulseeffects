package main

import (
	"reflect"
	"testing"
)

func TestParseBand(t *testing.T) {
	const b = "Filter 10: ON PK Fc 19808 Hz Gain -6.9 dB Q 0.50"

	band, i, err := ParseBand(b)
	if err != nil {
		t.Fatal("Failed to parse band:", err)
	}

	if i != 9 {
		t.Fatalf("Expected i = %d, got %d", 10, i)
	}

	expectBand := NewBand()
	expectBand.Frequency = 19808
	expectBand.Gain = -6.9
	expectBand.Quality = 0.50

	if !reflect.DeepEqual(band, expectBand) {
		t.Fatalf("Unexpected band returned: %#v", band)
	}
}
