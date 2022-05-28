package fpdecimal

import "errors"

const sep = '.'

var (
	errEmptyString            = errors.New("empty string")
	errMissingDigitsAfterSign = errors.New("missing digits after sign")
	errBadDigit               = errors.New("bad digit")
	errMultipleDots           = errors.New("multiple dots")
)

// ParseFixedPointDecimal parses fixed-point decimal of p fractions into int64.
func ParseFixedPointDecimal(s string, p int8) (int64, error) {
	if s == "" {
		return 0, errEmptyString
	}

	s0 := s
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
		if len(s) < 1 {
			return 0, errMissingDigitsAfterSign
		}
	}

	var d int8 = -1 // current decimal position
	var n int64     // output
	for _, ch := range []byte(s) {
		if d == p {
			break
		}

		if ch == sep {
			if d != -1 {
				return 0, errMultipleDots
			}
			d = 0
			continue
		}

		ch -= '0'
		if ch > 9 {
			return 0, errBadDigit
		}
		n = n*10 + int64(ch)

		if d != -1 {
			d++
		}
	}

	// fill rest of 0
	if d == -1 {
		d = 0
	}
	for i := d; i < p; i++ {
		n = n * 10
	}

	if s0[0] == '-' {
		n = -n
	}

	return n, nil
}
