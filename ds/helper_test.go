package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntComparator(t *testing.T) {
	tests := []struct {
		a                  int
		b                  int
		comparator         Comparator
		expectedComparison int
	}{
		{27, 27, &IntComparator{}, 0},
		{88, 27, &IntComparator{}, 1},
		{77, 99, &IntComparator{}, -1},
	}

	for _, test := range tests {
		comparison := test.comparator.Compare(test.a, test.b)

		assert.Equal(t, test.expectedComparison, comparison)
	}
}

func TestIntBitStringer(t *testing.T) {
	tests := []struct {
		a                 int
		bitStringer       BitStringer
		expectedBitString []byte
	}{
		{27, &IntBitStringer{}, []byte{0x32, 0x37}},
		{69, &IntBitStringer{}, []byte{0x36, 0x39}},
		{88, &IntBitStringer{}, []byte{0x38, 0x38}},
	}

	for _, test := range tests {
		bitString := test.bitStringer.BitString(test.a)

		assert.Equal(t, test.expectedBitString, bitString)
	}
}

func TestStringComparator(t *testing.T) {
	tests := []struct {
		a                  string
		b                  string
		comparator         Comparator
		expectedComparison int
	}{
		{"Same", "Same", &StringComparator{}, 0},
		{"Milad", "Jackie", &StringComparator{}, 1},
		{"Alice", "Bob", &StringComparator{}, -1},
	}

	for _, test := range tests {
		comparison := test.comparator.Compare(test.a, test.b)

		assert.Equal(t, test.expectedComparison, comparison)
	}
}

func TestStringBitStringer(t *testing.T) {
	tests := []struct {
		a                 string
		bitStringer       BitStringer
		expectedBitString []byte
	}{
		{"Barak", &StringBitStringer{}, []byte{0x42, 0x61, 0x72, 0x61, 0x6b}},
		{"Justin", &StringBitStringer{}, []byte{0x4a, 0x75, 0x73, 0x74, 0x69, 0x6e}},
		{"Milad", &StringBitStringer{}, []byte{0x4d, 0x69, 0x6c, 0x61, 0x64}},
	}

	for _, test := range tests {
		bitString := test.bitStringer.BitString(test.a)

		assert.Equal(t, test.expectedBitString, bitString)
	}
}
