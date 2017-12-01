package st

import (
	"bytes"
	"fmt"
)

// Assuming bit position is always equal greater than 1
func getBit(k []byte, pos int) byte {
	if pos > len(k)*8 {
		return 0 // padding with 0
	}

	pos--
	mask := byte(0x80 >> byte(pos%8))
	if k[pos/8]&mask == 0 {
		return 0
	}
	return 1
}

// Assuming k1 and k2 are not equal
func getDiffBitPos(k1, k2 []byte) int {
	var i, pos int
	var b1, b2 byte = 0, 0

	for i = 0; b1 == b2; i++ {
		if i < len(k1) {
			b1 = k1[i]
		} else {
			b1 = 0
		}
		if i < len(k2) {
			b2 = k2[i]
		} else {
			b2 = 0
		}
	}

	b1 = b1 ^ b2
	for pos = 0; b1 != 1; pos++ {
		b1 >>= 1
	}
	pos = i*8 - pos

	return pos
}

func getBinaryString(k []byte) string {
	var b bytes.Buffer
	for i := 0; i < len(k); i++ {
		b.WriteString(fmt.Sprintf("%b", k[i]))
	}

	return b.String()
}
