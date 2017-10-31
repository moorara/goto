package ds

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	intComparator  struct{}
	intBitStringer struct{}

	stringComparator  struct{}
	stringBitStringer struct{}
)

func (ic *intComparator) Compare(a Generic, b Generic) int {
	intA, _ := a.(int)
	intB, _ := b.(int)
	diff := intA - intB
	switch {
	case diff < 0:
		return -1
	case diff > 0:
		return 1
	default:
		return 0
	}
}

func (ib *intBitStringer) BitString(a Generic) []byte {
	intA, _ := a.(int)
	return []byte(strconv.Itoa(intA))
}

func (sc *stringComparator) Compare(a Generic, b Generic) int {
	strA, _ := a.(string)
	strB, _ := b.(string)
	return strings.Compare(strA, strB)
}

func (sb *stringBitStringer) BitString(a Generic) []byte {
	strA, _ := a.(string)
	return []byte(strA)
}

func TestIntComparator(t *testing.T) {
	tests := []struct {
		a                  int
		b                  int
		comparator         Comparator
		expectedComparison int
	}{
		{27, 27, &intComparator{}, 0},
		{88, 27, &intComparator{}, 1},
		{77, 99, &intComparator{}, -1},
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
		{27, &intBitStringer{}, []byte{0x32, 0x37}},
		{69, &intBitStringer{}, []byte{0x36, 0x39}},
		{88, &intBitStringer{}, []byte{0x38, 0x38}},
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
		{"Same", "Same", &stringComparator{}, 0},
		{"Milad", "Jackie", &stringComparator{}, 1},
		{"Alice", "Bob", &stringComparator{}, -1},
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
		{"Barak", &stringBitStringer{}, []byte{0x42, 0x61, 0x72, 0x61, 0x6b}},
		{"Justin", &stringBitStringer{}, []byte{0x4a, 0x75, 0x73, 0x74, 0x69, 0x6e}},
		{"Milad", &stringBitStringer{}, []byte{0x4d, 0x69, 0x6c, 0x61, 0x64}},
	}

	for _, test := range tests {
		bitString := test.bitStringer.BitString(test.a)

		assert.Equal(t, test.expectedBitString, bitString)
	}
}
