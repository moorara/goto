package ds

import (
	"strconv"
	"strings"
)

type (
	// IntComparator implements Comparator interface for int numbers
	IntComparator struct{}
	// IntBitStringer implements BitStringer interface for int numbers
	IntBitStringer struct{}

	// StringComparator implements Comparator interface for strings
	StringComparator struct{}
	// StringBitStringer implements BitStringer interface for strings
	StringBitStringer struct{}
)

// Compare compares two int numbers
func (ic *IntComparator) Compare(a Generic, b Generic) int {
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

// BitString returns the bit-string representation of an int number
func (ib *IntBitStringer) BitString(a Generic) []byte {
	intA, _ := a.(int)
	return []byte(strconv.Itoa(intA))
}

// Compare compares two strings
func (sc *StringComparator) Compare(a Generic, b Generic) int {
	strA, _ := a.(string)
	strB, _ := b.(string)
	return strings.Compare(strA, strB)
}

// BitString returns the bit-string representation of a string
func (sb *StringBitStringer) BitString(a Generic) []byte {
	strA, _ := a.(string)
	return []byte(strA)
}
