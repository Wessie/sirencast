package icecast

import "testing"

func TestMP3CalculatePadding(t *testing.T) {
	// calculate padding is expected to return at least the
	// padding amount given by atLeast when input is length.
	//
	// (length+padding) % 16 == 0 should be true
	tests := []struct {
		length  int
		atLeast int
	}{
		{16, 0},
		{8, 8},
		{32, 0},
		{64, 0},
		{44, 4},
		{40, 8},
		{60, 4},
		{20, 12},
	}

	for _, h := range tests {
		p := calculatePadding(h.length)
		if p < h.atLeast {
			t.Errorf("failed padding calculation: %d -> %d != %d\n", h.length, p, h.atLeast)
		}

		if (h.length+p)%16 != 0 {
			t.Errorf("length with padding is not multiple of 16: length: %d padding: %d\n",
				h.length, p)
		}
	}
}

func TestMP3ClientMetadata(t *testing.T) {
	// a metadata message is basically of the format
	// "StreamTitle='%s';" the rest of text takes up
	// 15 bytes, so use 15+meta as len expectation
	var titleLength = 15
	tests := []string{
		"test",
		"hello world",
		"world",
		"slightly longer than expected",
		"extra - dashes - included",
	}

	buf := make([]byte, 255*16+1)
	for _, s := range tests {
		// size byte component
		expectedLen := 1
		// metadata component
		expectedLen += titleLength + len(s)
		// padding component
		expectedLen += calculatePadding(expectedLen - 1)

		m := fillMetaBuffer(buf, s)
		if len(m) != expectedLen {
			t.Errorf("filled buffer is not of expected length: (%d != %d) %s\n",
				len(m), expectedLen, string(m))
		}
	}
}
