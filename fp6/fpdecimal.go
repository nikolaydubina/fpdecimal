package fp6

import "github.com/nikolaydubina/fpdecimal"

// Decimal with 6 fractional digits.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854.775807
// Min: -9223372036854.775808
type Decimal struct{ v int64 }

var Zero = Decimal{}

type integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

const (
	fractionDigits = 6
	multiplier     = 1_000_000
)

func FromInt[T integer](v T) Decimal { return Decimal{int64(v) * multiplier} }

func FromFloat[T float32 | float64](v T) Decimal {
	return Decimal{int64(float64(v) * float64(multiplier))}
}

// FromIntScaled expects value already scaled to minor units
func FromIntScaled[T integer](v T) Decimal { return Decimal{int64(v)} }

func FromString(s string) (Decimal, error) {
	v, err := fpdecimal.ParseFixedPointDecimal([]byte(s), fractionDigits)
	return Decimal{v}, err
}

func (v *Decimal) UnmarshalJSON(b []byte) (err error) {
	v.v, err = fpdecimal.ParseFixedPointDecimal(b, fractionDigits)
	return err
}

func (v Decimal) MarshalJSON() ([]byte, error) { return []byte(v.String()), nil }

func (v *Decimal) UnmarshalText(b []byte) (err error) {
	v.v, err = fpdecimal.ParseFixedPointDecimal(b, fractionDigits)
	return err
}

func (v Decimal) MarshalText() ([]byte, error) { return []byte(v.String()), nil }

func (a Decimal) Scaled() int64 { return a.v }

func (a Decimal) Float32() float32 { return float32(a.v) / float32(multiplier) }

func (a Decimal) Float64() float64 { return float64(a.v) / float64(multiplier) }

func (a Decimal) String() string { return fpdecimal.FixedPointDecimalToString(a.v, fractionDigits) }

func (a Decimal) Add(b Decimal) Decimal { return Decimal{v: a.v + b.v} }

func (a Decimal) Sub(b Decimal) Decimal { return Decimal{v: a.v - b.v} }

func (a Decimal) Mul(b Decimal) Decimal { return Decimal{v: a.v * b.v / multiplier} }

func (a Decimal) Div(b Decimal) Decimal { return Decimal{v: a.v * multiplier / b.v} }

func (a Decimal) Mod(b Decimal) Decimal { return Decimal{v: a.v % (b.v / multiplier)} }

func (a Decimal) DivMod(b Decimal) (part, remainder Decimal) { return a.Div(b), a.Mod(b) }

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

func Min(vs ...Decimal) Decimal {
	if len(vs) == 0 {
		panic("min of empty set is undefined")
	}
	var v Decimal = vs[0]
	for _, q := range vs {
		if q.LessThan(v) {
			v = q
		}
	}
	return v
}

func Max(vs ...Decimal) Decimal {
	if len(vs) == 0 {
		panic("max of empty set is undefined")
	}
	var v Decimal = vs[0]
	for _, q := range vs {
		if q.GreaterThan(v) {
			v = q
		}
	}
	return v
}
