package st

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestGetBit(t *testing.T) {
	tests := []struct {
		key          string
		bitString    BitString
		expectedBits []byte
	}{
		{
			"",
			BitStringString,
			[]byte{},
		},
		{
			"M",
			BitStringString,
			[]byte{
				0, 1, 0, 0, 1, 1, 0, 1,
			},
		},
		{
			"Milad",
			BitStringString,
			[]byte{
				0, 1, 0, 0, 1, 1, 0, 1,
				0, 1, 1, 0, 1, 0, 0, 1,
				0, 1, 1, 0, 1, 1, 0, 0,
				0, 1, 1, 0, 0, 0, 0, 1,
				0, 1, 1, 0, 0, 1, 0, 0,
			},
		},
		{
			"میلاد",
			BitStringString,
			[]byte{
				1, 1, 0, 1, 1, 0, 0, 1,
				1, 0, 0, 0, 0, 1, 0, 1,
				1, 1, 0, 1, 1, 0, 1, 1,
				1, 0, 0, 0, 1, 1, 0, 0,
				1, 1, 0, 1, 1, 0, 0, 1,
				1, 0, 0, 0, 0, 1, 0, 0,
				1, 1, 0, 1, 1, 0, 0, 0,
				1, 0, 1, 0, 0, 1, 1, 1,
				1, 1, 0, 1, 1, 0, 0, 0,
				1, 0, 1, 0, 1, 1, 1, 1,
			},
		},
	}

	for _, test := range tests {
		bitKey := test.bitString(test.key)
		for i := 0; i < len(bitKey); i++ {
			pos := i + 1
			bit := getBit(bitKey, pos)
			assert.Equal(t, test.expectedBits[i], bit)
		}

		pos := len(bitKey)*8 + 1
		assert.Zero(t, getBit(bitKey, pos))
	}
}

func TestGetDiffBitPos(t *testing.T) {
	tests := []struct {
		k1, k2      string
		bitString   BitString
		expectedPos int
	}{
		{"Key", "", BitStringString, 2},
		{"", "Key", BitStringString, 2},
		{"Milad", "Mona", BitStringString, 14},
		{"Milad", "میلاد", BitStringString, 1},
	}

	for _, test := range tests {
		bk1, bk2 := test.bitString(test.k1), test.bitString(test.k2)
		pos := getDiffBitPos(bk1, bk2)

		assert.Equal(t, test.expectedPos, pos)
	}
}

func TestGetBinaryString(t *testing.T) {
	tests := []struct {
		key                  string
		bitString            BitString
		expectedBinaryString string
	}{
		{"", BitStringString, ""},
		{"M", BitStringString, "1001101"},
		{"Milad", BitStringString, "10011011101001110110011000011100100"},
		{"میلاد", BitStringString, "11011001100001011101101110001100110110011000010011011000101001111101100010101111"},
	}

	for _, test := range tests {
		bitKey := test.bitString(test.key)
		binaryString := getBinaryString(bitKey)

		assert.Equal(t, test.expectedBinaryString, binaryString)
	}
}

func BenchmarkGetBit(b *testing.B) {
	bitKeys := make([][]byte, b.N)
	keys := util.GenerateStringSlice(b.N, 10, 100)
	for i := 0; i < b.N; i++ {
		bitKeys[i] = []byte(keys[i].(string))
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		getBit(bitKeys[n], 50)
	}
}

func BenchmarkGetDiffBitPos(b *testing.B) {
	var k1, k2 string
	bitKeys := make([][2][]byte, b.N)
	for i := 0; i < b.N; {
		k1 = util.GenerateString(10, 100)
		k2 = util.GenerateString(10, 100)
		if k1 != k2 {
			bitKeys[i][0] = []byte(k1)
			bitKeys[i][1] = []byte(k2)
		}
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		getDiffBitPos(bitKeys[n][0], bitKeys[n][1])
	}
}
