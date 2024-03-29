package fpdecimal

// Decimal is a decimal with fixed number of fraction digits.
// By default, uses 3 fractional digits.
// For example, values with 3 fractional digits will fit in ~9 quadrillion.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type Decimal struct{ v int64 }

var Zero = Decimal{}

var multipliers = [...]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}

type integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// FractionDigits that operations will use.
// Warning, after change, existing variables are not updated.
// Likely you want to use this once per runtime and in `func init()`.
var FractionDigits uint8 = 3

func FromInt[T integer](v T) Decimal { return Decimal{int64(v) * multipliers[FractionDigits]} }

func FromFloat[T float32 | float64](v T) Decimal {
	return Decimal{int64(float64(v) * float64(multipliers[FractionDigits]))}
}

// FromIntScaled expects value already scaled to minor units
func FromIntScaled[T integer](v T) Decimal { return Decimal{int64(v)} }

func FromString(s string) (Decimal, error) {
	v, err := ParseFixedPointDecimal([]byte(s), FractionDigits)
	return Decimal{v}, err
}

func (v *Decimal) UnmarshalJSON(b []byte) (err error) {
	v.v, err = ParseFixedPointDecimal(b, FractionDigits)
	return err
}

func (v Decimal) MarshalJSON() ([]byte, error) { return []byte(v.String()), nil }

func (a Decimal) Scaled() int64 { return a.v }

func (a Decimal) Float32() float32 { return float32(a.v) / float32(multipliers[FractionDigits]) }

func (a Decimal) Float64() float64 { return float64(a.v) / float64(multipliers[FractionDigits]) }

func (a Decimal) String() string { return FixedPointDecimalToString(a.v, FractionDigits) }

func (a Decimal) Add(b Decimal) Decimal { return Decimal{v: a.v + b.v} }

func (a Decimal) Sub(b Decimal) Decimal { return Decimal{v: a.v - b.v} }

func (a Decimal) Mul(b Decimal) Decimal { return Decimal{v: a.v * b.v / multipliers[FractionDigits]} }

func (a Decimal) Div(b Decimal) Decimal { return Decimal{v: a.v * multipliers[FractionDigits] / b.v} }

func (a Decimal) Mod(b Decimal) Decimal { return Decimal{v: a.v % (b.v / multipliers[FractionDigits])} }

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
