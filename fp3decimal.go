package fpdecimal

// FP3Decimal is fixed-point decimal with 3 fractions.
// Actual value scaled by 1000x.
// Values fit in ~9 quadrillion.
// Fractions lower than that are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type FP3Decimal int64

func FP3DecimalFromInt[T Integer](v T) FP3Decimal { return FP3Decimal(v) * 1000 }

func FP3DecimalFromDecimal[T Float](v T) FP3Decimal { return FP3Decimal(v) * 1000 }

func (v FP3Decimal) Decimal32() float32 { return float32(v) / 1000 }

func (v FP3Decimal) Decimal64() float64 { return float64(v) / 1000 }

func (v FP3Decimal) String() string { return FixedPointDecimalToString(int64(v), 3) }

// FP3DecimalFromString is based on hot-path for strconv.Atoi. Rounded down.
func FP3DecimalFromString(s string) (FP3Decimal, error) {
	v, err := ParseFixedPointDecimal(s, 3)
	return FP3Decimal(v), err
}

// UnmarshalJSON parses from JSON float number by treating it as bytes.
func (v *FP3Decimal) UnmarshalJSON(b []byte) (err error) {
	*v, err = FP3DecimalFromString(string(b))
	return err
}
