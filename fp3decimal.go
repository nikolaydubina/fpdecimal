package fpdecimal

// FP3Decimal is fixed-point decimal with 3 fractions.
// Actual value scaled by 1000x.
// Values fit in ~9 quadrillion.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type FP3Decimal struct{ v int64 }

var FP3DecimalZero = FP3Decimal{}

func FP3DecimalFromInt[T Integer](v T) FP3Decimal { return FP3Decimal{int64(v) * 1000} }

func FP3DecimalFromFloat[T Float](v T) FP3Decimal { return FP3Decimal{int64(float64(v) * 1000)} }

func (a FP3Decimal) Float32() float32 { return float32(a.v) / 1000 }

func (a FP3Decimal) Float64() float64 { return float64(a.v) / 1000 }

func (a FP3Decimal) String() string { return FixedPointDecimalToString(int64(a.v), 3) }

// FP3DecimalFromString is based on hot-path for strconv.Atoi. Rounded down.
func FP3DecimalFromString(s string) (FP3Decimal, error) {
	v, err := ParseFixedPointDecimal(s, 3)
	return FP3Decimal{v}, err
}

// UnmarshalJSON parses from JSON float number by treating it as bytes.
func (v *FP3Decimal) UnmarshalJSON(b []byte) (err error) {
	*v, err = FP3DecimalFromString(string(b))
	return err
}

func (a FP3Decimal) Add(b FP3Decimal) FP3Decimal { return FP3Decimal{v: a.v + b.v} }

func (a FP3Decimal) Sub(b FP3Decimal) FP3Decimal { return FP3Decimal{v: a.v - b.v} }

func (a FP3Decimal) Mul(b int) FP3Decimal { return FP3Decimal{v: a.v * int64(b)} }

func (a FP3Decimal) HigherThan(b FP3Decimal) bool { return a.v > b.v }

func (a FP3Decimal) LessThan(b FP3Decimal) bool { return a.v < b.v }
