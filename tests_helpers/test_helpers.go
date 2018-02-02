package testshelpers

import "math"

// FloatEquals : Returns two if two floats can be considered as equal.
func FloatEquals(a, b float64) bool {
	var eps = 0.00000001
	if math.Abs(a-b) < eps {
		return true
	}
	return false
}
