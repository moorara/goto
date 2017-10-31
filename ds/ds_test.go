package ds

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	value struct {
		v string
	}

	key struct {
		k string
	}

	bitStringKey struct {
		k []byte
	}
)

func newValueArray(items ...string) []value {
	arr := make([]value, len(items))
	for i, item := range items {
		arr[i] = value{item}
	}
	return arr
}

func (a value) Compare(b Value) int {
	c, _ := b.(value)
	return strings.Compare(a.v, c.v)
}

func (a key) Compare(b Key) int {
	c, _ := b.(key)
	return strings.Compare(a.k, c.k)
}

func (a bitStringKey) BitString() []byte {
	return a.k
}

func (a bitStringKey) Compare(b BitStringKey) int {
	c, _ := b.(bitStringKey)
	return bytes.Compare(a.k, c.k)
}

func TestValue(t *testing.T) {
	tests := []struct {
		a                  Value
		b                  Value
		expectedComparison int
	}{
		{value{"Same"}, value{"Same"}, 0},
		{value{"Alice"}, value{"Bob"}, -1},
		{value{"Milad"}, value{"Jackie"}, 1},
	}

	for _, test := range tests {
		comparison := test.a.Compare(test.b)

		assert.Equal(t, test.expectedComparison, comparison)
	}
}

func TestKey(t *testing.T) {
	tests := []struct {
		a                  Key
		b                  Key
		expectedComparison int
	}{
		{key{"Same"}, key{"Same"}, 0},
		{key{"Alice"}, key{"Bob"}, -1},
		{key{"Milad"}, key{"Jackie"}, 1},
	}

	for _, test := range tests {
		comparison := test.a.Compare(test.b)

		assert.Equal(t, test.expectedComparison, comparison)
	}
}

func TestBitStringKey(t *testing.T) {
	tests := []struct {
		a                  BitStringKey
		b                  BitStringKey
		expectedBitStringA []byte
		expectedBitStringB []byte
		expectedComparison int
	}{
		{bitStringKey{[]byte{0, 1, 2, 3}}, bitStringKey{[]byte{0, 1, 2, 3}}, []byte{0, 1, 2, 3}, []byte{0, 1, 2, 3}, 0},
		{bitStringKey{[]byte{2, 2, 2, 2}}, bitStringKey{[]byte{1, 1, 1, 1}}, []byte{2, 2, 2, 2}, []byte{1, 1, 1, 1}, 1},
		{bitStringKey{[]byte{4, 4, 4, 4}}, bitStringKey{[]byte{8, 8, 8, 8}}, []byte{4, 4, 4, 4}, []byte{8, 8, 8, 8}, -1},
	}

	for _, test := range tests {
		comparison := test.a.Compare(test.b)

		assert.Equal(t, test.expectedBitStringA, test.a.BitString())
		assert.Equal(t, test.expectedBitStringB, test.b.BitString())
		assert.Equal(t, test.expectedComparison, comparison)
	}
}
