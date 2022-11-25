package utils

import "math/rand"

const (
	ColorMainNode   = "paleturquoise"
	ColorFromToNode = "lightsalmon"
	ColorSender     = "yellowgreen"
	ColorRecipient  = "mistyrose"
)

func SmoothType() string {
	smooth := []string{
		"dynamic",
		"continuous",
		"discrete",
		"diagonalCross",
		"straightCross",
		"horizontal",
		"vertical",
		"curvedCW",
		"curvedCCW",
		"cubicBezier"}
	return smooth[rand.Intn((len(smooth)-1)-0)+0]
}

func SmoothRoundness() float64 {
	return (float64(rand.Intn((5-1)-1) + 1)) / float64(10)
}
