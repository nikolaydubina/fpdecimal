package fpdecimal

import (
	"strconv"
)

const zeroPrefix = "0.000000000000000000000000000000000000"

// FixedPointDecimalToString formats fixed-point decimal to string
func FixedPointDecimalToString(v int64, p uint8) string {
	// max int64: +9223372036854775.807
	// min int64: -9223372036854775.808
	// max bytes int64: 21
	b := make([]byte, 0, 21)
	b = AppendFixedPointDecimal(b, v, p)
	return string(b)
}

// AppendFixedPointDecimal appends formatted fixed point decimal to destination buffer.
// Returns appended slice.
// This is efficient for avoiding memory copy.
func AppendFixedPointDecimal(b []byte, v int64, p uint8) []byte {
	if v == 0 {
		return append(b, '0')
	}

	if p == 0 {
		return strconv.AppendInt(b, v, 10)
	}

	if v < 0 {
		v = -v
		b = append(b, '-')
	}

	// strconv.AppendInt is very efficient.
	// Efficient converting int64 to ASCII is not as trivial.
	s := uint8(len(b))
	b = strconv.AppendInt(b, v, 10)

	// has whole?
	if uint8(len(b))-s > p {
		// place decimal point
		i := uint8(len(b)) - p
		b = append(b, 0)
		copy(b[i+1:], b[i:])
		b[i] = '.'
	} else {
		// append zeroes and decimal point
		i := 2 + p - (uint8(len(b)) - s)
		for j := uint8(0); j < i; j++ {
			b = append(b, 0)
		}
		copy(b[s+i:], b[s:])
		copy(b[s:], []byte(zeroPrefix[:i]))
	}

	// remove trailing zeros
	n := 0
	for i, q := range b {
		if q != '0' {
			n = i
		}
	}
	if b[n] == '.' {
		n--
	}
	b = b[:n+1]

	return b
}
