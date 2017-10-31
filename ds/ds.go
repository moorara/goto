package ds

type (
	// Value represents a comparable value
	Value interface {
		Compare(Value) int
	}

	// Key represents a comparable key
	Key interface {
		Compare(Key) int
	}

	// BitStringKey represents a comparable bit-string key
	BitStringKey interface {
		BitString() []byte
		Compare(BitStringKey) int
	}
)
