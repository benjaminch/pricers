package testshelpers

import "math"

func FloatEquals(a, b float64) bool {
	var eps float64 = 0.00000001
	if math.Abs(a-b) < eps {
		return true
	}
	return false
}
