package fpfloat

// FP3Float is fixed-point float with 3 decimal fractions.
// Actual value scaled by 1000x.
// Values fit in ~9 quadrillion.
// Fractions are discarded in operations.
// Max: +9223372036854775.807
// Min: -9223372036854775.808
type FP3Float int64

func FP3FloatFromInt[T Integer](v T) FP3Float { return FP3Float(v) * 1000 }

func FP3FloatFromFloat[T Float](v T) FP3Float { return FP3Float(v) * 1000 }

func (v FP3Float) Float32() float32 { return float32(v) / 1000 }

func (v FP3Float) Float64() float64 { return float64(v) / 1000 }

func (v FP3Float) String() string { return FixedPointFloatToString(int64(v), 3) }

// FP3FloatFromString is based on hot-path for strconv.Atoi. Rounded down.
func FP3FloatFromString(s string) (FP3Float, error) {
	v, err := ParseFixedPointFloat(s, 3)
	return FP3Float(v), err
}

// UnmarshalJSON parses from JSON float number by treating it as bytes.
func (v *FP3Float) UnmarshalJSON(b []byte) (err error) {
	*v, err = FP3FloatFromString(string(b))
	return err
}
