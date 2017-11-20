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

func BenchmarkGetBit(b *testing.B) {
	bitKeys := make([][]byte, b.N)
	keys := util.GenerateStringSlice(b.N, 10, 100)
	for i := 0; i < len(bitKeys); i++ {
		bitKeys[i] = []byte(keys[i].(string))
	}

	for n := 0; n < b.N; n++ {
		getBit(bitKeys[n], 50)
	}
}
