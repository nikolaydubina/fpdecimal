package fpdecimal

import (
	"strconv"
)

const zeroPrefix = "0.000000000000000000000000000000000000"

// FixedPointDecimalToString formats fixed-point decimal to string.
// strconv.AppendInt is very efficient.
// Efficient converting int64 to ASCII is not as trivial.
func FixedPointDecimalToString(v int64, p int) string {
	if v == 0 {
		return "0"
	}

	if p == 0 {
		return strconv.FormatInt(v, 10)
	}

	// max int64: +9223372036854775.807
	// min int64: -9223372036854775.808
	// max bytes int64: 21
	b := make([]byte, 1, 21)

	if v < 0 {
		v = -v
		b[0] = '-'
	}

	b = strconv.AppendInt(b, int64(v), 10)

	if len(b) > (p + 1) {
		i := len(b) - p
		b = append(b, 0)
		copy(b[i+1:], b[i:])
		b[i] = '.'
	} else {
		i := 3 + p - len(b)
		for j := 0; j < i; j++ {
			b = append(b, 0)
		}
		copy(b[i:], b)
		copy(b[1:], []byte(zeroPrefix[:i]))
	}

	if b[0] != '-' {
		b = b[1:]
	}

	return string(b)
}
