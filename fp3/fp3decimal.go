package fp3

import (
	"github.com/nikolaydubina/fpdecimal"
	"github.com/nikolaydubina/fpdecimal/constraints"
)

// Decimal is fixed-point decimal with 3 fractions.
// Actual value scaled by 1000x.
// Values fit in ~9 quadrillion.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type Decimal struct{ v int64 }

var Zero = Decimal{}

func FromInt[T constraints.Integer](v T) Decimal { return Decimal{int64(v) * 1000} }

func FromFloat[T constraints.Float](v T) Decimal {
	return Decimal{int64(float64(v) * 1000)}
}

func (a Decimal) Float32() float32 { return float32(a.v) / 1000 }

func (a Decimal) Float64() float64 { return float64(a.v) / 1000 }

func (a Decimal) String() string { return fpdecimal.FixedPointDecimalToString(int64(a.v), 3) }

func FromString(s string) (Decimal, error) {
	v, err := fpdecimal.ParseFixedPointDecimal(s, 3)
	return Decimal{v}, err
}

func (v *Decimal) UnmarshalJSON(b []byte) (err error) {
	*v, err = FromString(string(b))
	return err
}

func (a Decimal) Add(b Decimal) Decimal { return Decimal{v: a.v + b.v} }

func (a Decimal) Sub(b Decimal) Decimal { return Decimal{v: a.v - b.v} }

func (a Decimal) Mul(b int) Decimal { return Decimal{v: a.v * int64(b)} }

func (a Decimal) HigherThan(b Decimal) bool { return a.v > b.v }

func (a Decimal) LessThan(b Decimal) bool { return a.v < b.v }