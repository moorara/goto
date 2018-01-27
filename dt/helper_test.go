package dt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntComparator(t *testing.T) {
	tests := []struct {
		a                  int
		b                  int
		compare            Compare
		expectedComparison int
	}{
		{27, 27, CompareInt, 0},
		{88, 27, CompareInt, 1},
		{77, 99, CompareInt, -1},
	}

	for _, tc := range tests {
		cmp := tc.compare(tc.a, tc.b)

		assert.Equal(t, tc.expectedComparison, cmp)
	}
}

func TestIntBitStringer(t *testing.T) {
	tests := []struct {
		a                 int
		bitString         BitString
		expectedBitString []byte
	}{
		{27, BitStringInt, []byte{0x32, 0x37}},
		{69, BitStringInt, []byte{0x36, 0x39}},
		{88, BitStringInt, []byte{0x38, 0x38}},
	}

	for _, tc := range tests {
		bitString := tc.bitString(tc.a)

		assert.Equal(t, tc.expectedBitString, bitString)
	}
}

func TestStringComparator(t *testing.T) {
	tests := []struct {
		a                  string
		b                  string
		compare            Compare
		expectedComparison int
	}{
		{"Same", "Same", CompareString, 0},
		{"Milad", "Jackie", CompareString, 1},
		{"Alice", "Bob", CompareString, -1},
	}

	for _, tc := range tests {
		cmp := tc.compare(tc.a, tc.b)

		assert.Equal(t, tc.expectedComparison, cmp)
	}
}

func TestStringBitStringer(t *testing.T) {
	tests := []struct {
		a                 string
		bitString         BitString
		expectedBitString []byte
	}{
		{"Barak", BitStringString, []byte{0x42, 0x61, 0x72, 0x61, 0x6b}},
		{"Justin", BitStringString, []byte{0x4a, 0x75, 0x73, 0x74, 0x69, 0x6e}},
		{"Milad", BitStringString, []byte{0x4d, 0x69, 0x6c, 0x61, 0x64}},
	}

	for _, tc := range tests {
		bitString := tc.bitString(tc.a)

		assert.Equal(t, tc.expectedBitString, bitString)
	}
}
