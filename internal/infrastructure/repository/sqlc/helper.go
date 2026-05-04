package sqlc

import "math"

// safeInt32 strictly clamps an int64 to the valid bounds of an int32.
func safeInt32(value int64) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}
