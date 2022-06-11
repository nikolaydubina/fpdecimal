package fpdecimal_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/nikolaydubina/fpdecimal"
)

func FuzzFixedPointDecimalToString(f *testing.F) {
	tests := []float32{
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
	f.Fuzz(func(t *testing.T, r float32) {
		if r > math.MaxInt64/1000 || r < math.MinInt64/1000 {
			t.Skip()
		}

		s := fmt.Sprintf("%.3f", r)

		v, err := fpdecimal.ParseFixedPointDecimal(s, 3)
		if err != nil {
			t.Errorf(err.Error())
		}

		if s == "-0.000" || s == "0.000" || r == 0 || r == -0 {
			if q := fpdecimal.FixedPointDecimalToString(v, 3); q != "0" {
				t.Error(r, s, q)
			}
			return
		}

		if fs := fpdecimal.FixedPointDecimalToString(v, 3); s != fs {
			t.Error(s, fs, f, r, v)
		}
	})
}

func FuzzFixedPointDecimalToString_NoFractions(f *testing.F) {
	tests := []int64{
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
