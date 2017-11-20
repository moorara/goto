package st

// Assuming bit position is always equal greater than 1
func getBit(k []byte, pos int) byte {
	if pos > len(k)*8 {
		return 0
	}

	pos--
	mask := byte(0x80 >> byte(pos%8))
	if k[pos/8]&mask == 0 {
		return 0
	}
	return 1
}

func getDiffBitPos(k1, k2 []byte) int {
	return 0
}

func getBits() string {
	return ""
}
