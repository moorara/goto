package ds

type (
	// Value represents a generic data type
	Generic interface{}

	// Comparator is used for comparing
	Comparator interface {
		Compare(a Generic, b Generic) int
	}

	// BitStringer is used for bit-string representation
	BitStringer interface {
		BitString(a Generic) []byte
	}
)
