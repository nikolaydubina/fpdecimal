package fpdecimal_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"testing"
	"unsafe"

	"github.com/nikolaydubina/fpdecimal"
)

func FuzzFP3Decimal_ParseStringSameAsDecimal(f *testing.F) {
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

		v, err := fpdecimal.FP3DecimalFromString(s)
		if err != nil {
			t.Errorf(err.Error())
		}

		if s == "-0.000" || s == "0.000" || r == 0 || r == -0 {
			if v.String() != "0" {
				t.Errorf("s('0') != FP3Decimal.String(%#v) of fpdecimal(%#v) float32(%#v) .3f-float32(%#v)", v.String(), v, r, s)
			}
			return
		}

		if s != v.String() {
			t.Errorf("s(%#v) != FP3Decimal.String(%#v) of fpdecimal(%#v) float32(%#v)", s, v.String(), v, r)
		}
	})
}

func FuzzFP3Decimal_ParseStringRaw(f *testing.F) {
	tests := []string{
		"123.456",
		"0.123",
		"0.1",
		"0.01",
		"0.001",
		"0.000",
		"0.123.2",
		"0..1",
		"0.1.2",
		"123.1o2",
		"--123",
		"00000.123",
		"-",
		"",
		"123456",
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add("-" + tc)
	}
	f.Fuzz(func(t *testing.T, s string) {
		v, err := fpdecimal.FP3DecimalFromString(s)
		if err != nil {
			if v != fpdecimal.FP3DecimalZero {
				t.Errorf("has to be 0 on error")
			}
			return
		}
	})
}

func FuzzFP3Decimal_ToType(f *testing.F) {
	tests := []float64{
		0,
		0.001,
		1,
		123.456,
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add(-tc)
	}
	f.Fuzz(func(t *testing.T, v float64) {
		a := fpdecimal.FP3DecimalFromFloat(v)

		if float32(v) != a.Float32() {
			t.Error(a, a.Float32(), float32(v))
		}

		if v != a.Float64() {
			t.Error(a, a.Float32(), v)
		}
	})
}

var testsFloats = []struct {
	name string
	vals []string
}{
	{
		name: "small",
		vals: []string{
			"123.456",
			"0.123",
			"0.012",
			"0.001",
			"0.982",
			"0.101",
			"10",
			"11",
			"1",
		},
	},
	{
		name: "large",
		vals: []string{
			"123123123112312.1232",
			"5341320482340234.123",
		},
	},
}

func BenchmarkParse_FP3Decimal(b *testing.B) {
	var s fpdecimal.FP3Decimal
	var err error
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s, err = fpdecimal.FP3DecimalFromString(tc.vals[n%len(tc.vals)])
				if err != nil || s == fpdecimal.FP3DecimalZero {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkPrint_FP3Decimal(b *testing.B) {
	var s string
	for _, tc := range testsFloats {
		tests := make([]fpdecimal.FP3Decimal, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := fpdecimal.FP3DecimalFromString(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
		}

		b.ResetTimer()
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s = tests[n%len(tc.vals)].String()
				if s == "" {
					b.Error("empty str")
				}
			}
		})
	}
}

func TestFP3Decimal_UnmarshalJSON(t *testing.T) {
	type MyType struct {
		TeslaStockPrice fpdecimal.FP3Decimal `json:"tesla-stock-price"`
	}

	tests := []struct {
		json string
		v    fpdecimal.FP3Decimal
		s    string
	}{
		{
			json: `{"tesla-stock-price": 9000.001}`,
			v:    fpdecimal.FP3DecimalFromFloat(9000.001),
			s:    `9000.001`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.json, func(t *testing.T) {
			var v MyType
			err := json.Unmarshal([]byte(tc.json), &v)
			if err != nil {
				t.Error(err)
			}
			if s := v.TeslaStockPrice.String(); s != tc.s {
				t.Errorf("s(%#v) != tc.s(%#v)", s, tc.s)
			}
		})
	}
}

func FuzzFP3Decimal_UnmarshalJSON(f *testing.F) {
	type MyType struct {
		A fpdecimal.FP3Decimal `json:"a"`
	}

	tests := []float32{
		123456,
		0,
		1.1,
		0.123,
		0.0123,
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add(-tc)
	}
	f.Fuzz(func(t *testing.T, v float32) {
		b := fmt.Sprintf("%.3f", v)
		s := `{"a": ` + b + `}`

		var x MyType
		err := json.Unmarshal([]byte(s), &x)
		if err != nil {
			t.Error(err, s)
		}

		if b == "-0.000" || b == "0.000" || v == 0 || v == -0 {
			if x.A.String() != "0" {
				t.Error(b, x)
			}
			return
		}

		if a := x.A.String(); a != b {
			t.Error(a, b)
		}
	})
}

func ExampleFP3Decimal() {
	var BuySP500Price = fpdecimal.FP3DecimalFromInt(9000)

	input := []byte(`{"sp500": 9000.023}`)

	type Stocks struct {
		SP500 fpdecimal.FP3Decimal `json:"sp500"`
	}
	var v Stocks
	if err := json.Unmarshal(input, &v); err != nil {
		log.Fatal(err)
	}

	var amountToBuy fpdecimal.FP3Decimal
	if v.SP500.HigherThan(BuySP500Price) {
		amountToBuy = amountToBuy.Add(v.SP500.Mul(2))
	}

	fmt.Println(amountToBuy)
	// Output: 18000.046
}

func FuzzFP3Decimal_AddSub(f *testing.F) {
	tests := [][2]float32{
		{1, 2},
		{1, -5},
		{1, 0},
		{1.1, -0.002},
	}
	for _, tc := range tests {
		f.Add(tc[0], tc[1])
	}
	f.Fuzz(func(t *testing.T, a, b float32) {
		fa := fpdecimal.FP3DecimalFromFloat(a)
		fb := fpdecimal.FP3DecimalFromFloat(b)

		if fa.Add(fb) != fb.Add(fa) {
			t.Error(a, b)
		}

		if fpdecimal.FP3DecimalZero.Add(fa).Add(fb).Add(fa) != fpdecimal.FP3DecimalZero.Add(fb).Add(fa).Add(fa) {
			t.Error(a, b)
		}

		if v := fa.Add(fb).Sub(fb); v != fa {
			t.Error(a, b, v)
		}

		if fpdecimal.FP3DecimalZero.Add(fa).Sub(fa) != fpdecimal.FP3DecimalZero {
			t.Error(a)
		}

	})
}

func FuzzFP3Decimal_AddSub_Int(f *testing.F) {
	tests := [][2]int{
		{1, 2},
		{2, -5},
		{1, 0},
		{111, -5},
		{0, 0},
	}
	for _, tc := range tests {
		f.Add(tc[0], tc[1])
	}
	f.Fuzz(func(t *testing.T, a, b int) {
		fa := fpdecimal.FP3DecimalFromInt(a)
		fb := fpdecimal.FP3DecimalFromInt(b)

		if a < b {
			if !fa.LessThan(fb) {
				t.Error(a, b, fa, fb)
			}
		}

		if a > b {
			if !fa.HigherThan(fb) {
				t.Error(a, b, fa, fb)
			}
		}

		if a == b {
			if fa != fb {
				t.Error(a, b, fa, fb)
			}
		}
	})
}

func BenchmarkArithmetic_FP3Decimal(b *testing.B) {
	x, _ := fpdecimal.FP3DecimalFromString("251.231")
	y, _ := fpdecimal.FP3DecimalFromString("21231.001")

	var s fpdecimal.FP3Decimal

	b.Run("add_x1", func(b *testing.B) {
		s = fpdecimal.FP3DecimalZero
		for n := 0; n < b.N; n++ {
			s = x.Add(y)
		}
	})

	b.Run("add_x100", func(b *testing.B) {
		s = fpdecimal.FP3DecimalZero
		for n := 0; n < b.N; n++ {
			for i := 0; i < 100; i++ {
				s = x.Add(y)
			}
		}
	})

	if s == fpdecimal.FP3DecimalZero {
		b.Error()
	}
}

func TestFP3Decimal_memlayout(t *testing.T) {
	a, _ := fpdecimal.FP3DecimalFromString("-1000.123")
	if v := unsafe.Sizeof(a); v != 8 {
		t.Error(a, v)
	}
}
