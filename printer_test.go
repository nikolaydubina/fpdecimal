package fpdecimal_test

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/nikolaydubina/fpdecimal"
)

func FuzzFixedPointDecimalToString(f *testing.F) {
	tests := []float64{
		0,
		0.100,
		0.101,
		0.010,
		0.001,
		0.0001,
		0.123,
		0.103,
		0.100001,
		12.001,
		12.010,
		12.345,
		1,
		2,
		10,
		12345678,
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add(-tc)
	}
	f.Fuzz(func(t *testing.T, r float64) {
		if r > math.MaxInt64/1000 || r < math.MinInt64/1000 {
			t.Skip()
		}

		s := fmt.Sprintf("%.3f", r)
		rs, _ := strconv.ParseFloat(s, 64)

		v, err := fpdecimal.ParseFixedPointDecimal(s, 3)
		if err != nil {
			t.Errorf(err.Error())
		}

		if s == "-0.000" || s == "0.000" || rs == 0 || rs == -0 || (rs > -0.001 && rs < 0.001) {
			if q := fpdecimal.FixedPointDecimalToString(v, 3); q != "0" {
				t.Error(r, s, q)
			}
			return
		}

		if s, fs := strconv.FormatFloat(rs, 'f', -1, 64), fpdecimal.FixedPointDecimalToString(v, 3); s != fs {
			t.Error(s, fs, r, v)
		}
	})
}

func FuzzFixedPointDecimalToString_NoFractions(f *testing.F) {
	tests := []int64{
		0,
		1,
		2,
		10,
		12345678,
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add(-tc)
	}
	f.Fuzz(func(t *testing.T, r int64) {
		if a, b := fpdecimal.FixedPointDecimalToString(r, 0), strconv.FormatInt(r, 10); a != b {
			t.Error(a, b)
		}
	})
}

func BenchmarkFixedPointDecimalToString(b *testing.B) {
	var s string
	for _, tc := range testsFloats {
		tests := make([]int64, 0, len(tc.vals))
		for range tc.vals {
			tests = append(tests, int64(rand.Int()))
		}

		b.ResetTimer()
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s = fpdecimal.FixedPointDecimalToString(tests[n%len(tests)], 3)
				if s == "" {
					b.Error("empty str")
				}
			}
		})
	}
}

func BenchmarkAppendFixedPointDecimal(b *testing.B) {
	d := make([]byte, 0, 21)
	for _, tc := range testsFloats {
		tests := make([]int64, 0, len(tc.vals))
		for range tc.vals {
			tests = append(tests, int64(rand.Int()))
		}

		b.ResetTimer()
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				d = fpdecimal.AppendFixedPointDecimal(d, tests[n%len(tests)], 3)
				if len(d) == 0 {
					b.Error("empty str")
				}
			}
		})
	}
}
