package fpdecimal

import (
	"strconv"
)

const zeroPrefix = "0.000000000000000000000000000000000000"

// FixedPointDecimalToString formats fixed-point decimal to string.
func FixedPointDecimalToString(v int64, p int) string {
	if v == 0 {
		return "0"
	}

	n := false
	if v < 0 {
		v = -v
		n = true
	}

	s := strconv.FormatInt(int64(v), 10)

	if len(s) > p {
		p := len(s) - p
		s = s[:p] + "." + s[p:]
	} else {
		s = zeroPrefix[:(2+p-len(s))] + s
	}

	if n {
		s = "-" + s
	}

	return s
}
