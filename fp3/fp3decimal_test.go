package fp3_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"testing"
	"unsafe"

	"github.com/nikolaydubina/fpdecimal/fp3"
)

func FuzzArithmetics(f *testing.F) {
	tests := [][2]int64{
		{1, 2},
		{1, -5},
		{1, 0},
		{1100, -2},
	}
	for _, tc := range tests {
		f.Add(tc[0], tc[1])
	}
	f.Fuzz(func(t *testing.T, a, b int64) {
		fa := fp3.FromIntScaled(a)
		fb := fp3.FromIntScaled(b)

		v := []bool{
			// sum commutativity
			fa.Add(fb) == fb.Add(fa),

			// sum associativity
			fp3.Zero.Add(fa).Add(fb).Add(fa) == fp3.Zero.Add(fb).Add(fa).Add(fa),

			// sum zero
			fa == fa.Add(fb).Sub(fb),
			fa == fa.Sub(fb).Add(fb),
			fp3.Zero == fp3.Zero.Add(fa).Sub(fa),

			// product identity
			fa == fa.Mul(1),

			// product zero
			fp3.Zero == fa.Mul(0),

			// match number
			(a == b) == (fa == fb),
			a < b == fa.LessThan(fb),
			a > b == fa.GreaterThan(fb),
			a <= b == fa.LessThanOrEqual(fb),
			a >= b == fa.GreaterThanOrEqual(fb),

			// match number convert
			fp3.FromIntScaled(a+b) == fa.Add(fb),
			fp3.FromIntScaled(a-b) == fa.Sub(fb),
		}
		for i, q := range v {
			if !q {
				t.Error(i, a, b, fa, fb)
			}
		}
	})
}

func FuzzParse_StringSameAsFloat(f *testing.F) {
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

		v, err := fp3.FromString(s)
		if err != nil {
			t.Errorf(err.Error())
		}

		if s == "-0.000" || s == "0.000" || r == 0 || r == -0 {
			if v.String() != "0" {
				t.Errorf("s('0') != Decimal.String(%#v) of fp3(%#v) float32(%#v) .3f-float32(%#v)", v.String(), v, r, s)
			}
			return
		}

		if s != v.String() {
			t.Errorf("s(%#v) != Decimal.String(%#v) of fp3(%#v) float32(%#v)", s, v.String(), v, r)
		}
	})
}

func FuzzParse_StringRaw(f *testing.F) {
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
		v, err := fp3.FromString(s)
		if err != nil {
			if v != fp3.Zero {
				t.Errorf("has to be 0 on error")
			}
			return
		}
	})
}

func FuzzToFloat(f *testing.F) {
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
		a := fp3.FromFloat(v)

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

func BenchmarkParse(b *testing.B) {
	var s fp3.Decimal
	var err error
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s, err = fp3.FromString(tc.vals[n%len(tc.vals)])
				if err != nil || s == fp3.Zero {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkPrint(b *testing.B) {
	var s string
	for _, tc := range testsFloats {
		tests := make([]fp3.Decimal, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := fp3.FromString(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
			tests = append(tests, fp3.Zero.Sub(v))
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

func TestUnmarshalJSON(t *testing.T) {
	type MyType struct {
		TeslaStockPrice fp3.Decimal `json:"tesla-stock-price"`
	}

	tests := []struct {
		json string
		v    fp3.Decimal
		s    string
	}{
		{
			json: `{"tesla-stock-price": 9000.001}`,
			v:    fp3.FromFloat(9000.001),
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

func FuzzUnmarshalJSON(f *testing.F) {
	type MyType struct {
		A fp3.Decimal `json:"a"`
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

func ExampleDecimal() {
	var BuySP500Price = fp3.FromInt(9000)

	input := []byte(`{"sp500": 9000.023}`)

	type Stocks struct {
		SP500 fp3.Decimal `json:"sp500"`
	}
	var v Stocks
	if err := json.Unmarshal(input, &v); err != nil {
		log.Fatal(err)
	}

	var amountToBuy fp3.Decimal
	if v.SP500.GreaterThan(BuySP500Price) {
		amountToBuy = amountToBuy.Add(v.SP500.Mul(2))
	}

	fmt.Println(amountToBuy)
	// Output: 18000.046
}

func BenchmarkArithmetic(b *testing.B) {
	x, _ := fp3.FromString("251.231")
	y, _ := fp3.FromString("21231.001")

	var s fp3.Decimal

	b.Run("add_x1", func(b *testing.B) {
		s = fp3.Zero
		for n := 0; n < b.N; n++ {
			s = x.Add(y)
		}
	})

	b.Run("add_x100", func(b *testing.B) {
		s = fp3.Zero
		for n := 0; n < b.N; n++ {
			for i := 0; i < 100; i++ {
				s = x.Add(y)
			}
		}
	})

	if s == fp3.Zero {
		b.Error()
	}
}

func TestDecimalMemoryLayout(t *testing.T) {
	a, _ := fp3.FromString("-1000.123")
	if v := unsafe.Sizeof(a); v != 8 {
		t.Error(a, v)
	}
}
