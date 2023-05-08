package fpdecimal

import (
	"github.com/nikolaydubina/fpdecimal/constraints"
)

// Decimal is a decimal with fixed number of fraction digits.
// By default, uses 3 fractional digits.
// For example, values with 3 fractional digits will fit in ~9 quadrillion.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type Decimal struct{ v int64 }

var Zero = Decimal{}

var multipliers = []int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}

var fractionDigits int8 = 3

func SetFractionDigits(fraction uint8) bool {
	if fraction > 10 {
		return false
	}
	fractionDigits = int8(fraction)
	return true
}

func FromInt[T constraints.Integer](v T) Decimal {
	return Decimal{int64(v) * multipliers[fractionDigits]}
}

func FromFloat[T constraints.Float](v T) Decimal {
	return Decimal{int64(float64(v) * float64(multipliers[fractionDigits]))}
}

func FromIntScaled[T constraints.Integer](v T) Decimal { return Decimal{int64(v)} }

func (a Decimal) Float32() float32 { return float32(a.v) / float32(multipliers[fractionDigits]) }

func (a Decimal) Float64() float64 { return float64(a.v) / float64(multipliers[fractionDigits]) }

func (a Decimal) String() string { return FixedPointDecimalToString(a.v, int(fractionDigits)) }

func (a Decimal) Add(b Decimal) Decimal { return Decimal{v: a.v + b.v} }

func (a Decimal) Sub(b Decimal) Decimal { return Decimal{v: a.v - b.v} }

func (a Decimal) Mul(b int) Decimal { return Decimal{v: a.v * int64(b)} }

func (a Decimal) Div(b int) (part Decimal, remainder Decimal) {
	return Decimal{v: a.v / int64(b)}, Decimal{v: a.v % int64(b)}
}

func (a Decimal) Equal(b Decimal) bool { return a.v == b.v }

func (a Decimal) GreaterThan(b Decimal) bool { return a.v > b.v }

func (a Decimal) LessThan(b Decimal) bool { return a.v < b.v }

func (a Decimal) GreaterThanOrEqual(b Decimal) bool { return a.v >= b.v }

func (a Decimal) LessThanOrEqual(b Decimal) bool { return a.v <= b.v }

func (a Decimal) Compare(b Decimal) int {
	if a.LessThan(b) {
		return -1
	}
	if a.GreaterThan(b) {
		return 1
	}
	return 0
}

func FromString(s string) (Decimal, error) {
	v, err := ParseFixedPointDecimal(s, fractionDigits)
	return Decimal{v}, err
}

func (v *Decimal) UnmarshalJSON(b []byte) (err error) {
	v.v, err = ParseFixedPointDecimal(string(b), fractionDigits)
	return err
}

func (v Decimal) MarshalJSON() ([]byte, error) { return []byte(v.String()), nil }
