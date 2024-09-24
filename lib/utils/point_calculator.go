package utils

import "math"

var (
	POINT_DIVIDER = 10000
)

// Calculate VP Point
// At the moment, POINT was set to 10000
func CalculateUserRedeemPoint(invoiceAmount float64) int {
	return int(math.Floor((float64(invoiceAmount) / float64(POINT_DIVIDER))))
}
