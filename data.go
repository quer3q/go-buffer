package buffer

type (
	Data interface {
		Push(v interface{}) int // Push value and return new data size
		Empty() bool
	}

	DataConstructor func() Data
)
