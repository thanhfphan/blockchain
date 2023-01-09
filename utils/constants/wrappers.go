package constants

import "math"

const (
	ByteLen  = 1
	ShortLen = 2
	IntLen   = 4
	LongLen  = 8
	BoolLen  = 1

	MaxStringLen = math.MaxUint16
	IPLen        = 16 + ShortLen
)
