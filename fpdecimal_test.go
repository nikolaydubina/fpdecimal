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
		fa := fpdecimal.FromIntScaled(a)
		fb := fpdecimal.FromIntScaled(b)

		v := []bool{
			// sum commutativity
			fa.Add(fb) == fb.Add(fa),

			// sum associativity
			fpdecimal.Zero.Add(fa).Add(fb).Add(fa) == fpdecimal.Zero.Add(fb).Add(fa).Add(fa),

			// sum zero
			fa == fa.Add(fb).Sub(fb),
			fa == fa.Sub(fb).Add(fb),
			fpdecimal.Zero == fpdecimal.Zero.Add(fa).Sub(fa),

			// product identity
			fa == fa.Mul(1),

			// product zero
			fpdecimal.Zero == fa.Mul(0),

			// match number
			(a == b) == (fa == fb),
			(a == b) == fa.Equal(fb),
			a < b == fa.LessThan(fb),
			a > b == fa.GreaterThan(fb),
			a <= b == fa.LessThanOrEqual(fb),
			a >= b == fa.GreaterThanOrEqual(fb),

			// match number convert
			fpdecimal.FromIntScaled(a+b) == fa.Add(fb),
			fpdecimal.FromIntScaled(a-b) == fa.Sub(fb),
		}
		for i, q := range v {
			if !q {
				t.Error(i, a, b, fa, fb)
			}
		}

		if b != 0 {
			w, r := fa.Div(int(b))
			if w != fpdecimal.FromIntScaled(a/b) {
				t.Error(w, a/b, a, b, fa)
			}
			if r != fpdecimal.FromIntScaled(a%b) {
				t.Error(r, a%b, a, b, fa)
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

		v, err := fpdecimal.FromString(s)
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
		v, err := fpdecimal.FromString(s)
		if err != nil {
			if v != fpdecimal.Zero {
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
		a := fpdecimal.FromFloat(v)

		if float32(v) != a.Float32() {
			t.Error(a, a.Float32(), float32(v))
		}

		if v != a.Float64() {
			t.Error(a, a.Float32(), v)
		}
	})
}

func FuzzScaled(f *testing.F) {
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
		a := fpdecimal.FromFloat(v)

		if int64(v*1000) != a.Scaled() {
			t.Error(a, a.Scaled(), int64(v*1000))
		}
	})
}

var floatsForTests = []struct {
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
	var s fpdecimal.Decimal
	var err error
	for _, tc := range floatsForTests {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s, err = fpdecimal.FromString(tc.vals[n%len(tc.vals)])
				if err != nil || s == fpdecimal.Zero {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkPrint(b *testing.B) {
	var s string
	for _, tc := range floatsForTests {
		tests := make([]fpdecimal.Decimal, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := fpdecimal.FromString(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
			tests = append(tests, fpdecimal.Zero.Sub(v))
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
		TeslaStockPrice fpdecimal.Decimal `json:"tesla-stock-price"`
	}

	tests := []struct {
		json string
		v    fpdecimal.Decimal
		s    string
	}{
		{
			json: `{"tesla-stock-price": 9000.001}`,
			v:    fpdecimal.FromFloat(9000.001),
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

func TestUMarshalJSON(t *testing.T) {
	type MyType struct {
		TeslaStockPrice fpdecimal.Decimal `json:"tesla-stock-price"`
	}

	t.Run("when nil struct, then error", func(t *testing.T) {
		var v *MyType
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), v)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when nil value, then error", func(t *testing.T) {
		var v *fpdecimal.Decimal
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), &v)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when nil const of type, then error", func(t *testing.T) {
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), (*fpdecimal.Decimal)(nil))
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		var v MyType
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), &v)
		if err != nil {
			t.Error(err)
		}
		e := MyType{fpdecimal.FromIntScaled(9000001)}
		if v != e {
			t.Error(v)
		}
	})
}

func FuzzJSON(f *testing.F) {
	type MyType struct {
		A fpdecimal.Decimal `json:"a"`
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
		if v > math.MaxInt64/1000 || v < math.MinInt64/1000 {
			t.Skip()
		}

		b := fmt.Sprintf("%.3f", v)
		s := `{"a":` + b + `}`

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

		ms, err := json.Marshal(x)
		if err != nil {
			t.Error(err)
		}
		if string(ms) != s {
			t.Error(s, string(ms), x)
		}
	})
}

func ExampleDecimal() {
	var BuySP500Price = fpdecimal.FromInt(9000)

	input := []byte(`{"sp500": 9000.023}`)

	type Stocks struct {
		SP500 fpdecimal.Decimal `json:"sp500"`
	}
	var v Stocks
	if err := json.Unmarshal(input, &v); err != nil {
		log.Fatal(err)
	}

	var amountToBuy fpdecimal.Decimal
	if v.SP500.GreaterThan(BuySP500Price) {
		amountToBuy = amountToBuy.Add(v.SP500.Mul(2))
	}

	fmt.Println(amountToBuy)
	// Output: 18000.046
}

func ExampleDecimal_Div_remainder() {
	x, _ := fpdecimal.FromString("1.000")

	a, r := x.Div(3)
	fmt.Println(a, r)
	// Output: 0.333 0.001
}

func ExampleDecimal_Div_whole() {
	x, _ := fpdecimal.FromString("1.000")

	a, r := x.Div(5)
	fmt.Println(a, r)
	// Output: 0.200 0
}

func BenchmarkArithmetic(b *testing.B) {
	x, _ := fpdecimal.FromString("251.231")
	y, _ := fpdecimal.FromString("21231.001")

	var s fpdecimal.Decimal

	b.Run("add_x1", func(b *testing.B) {
		s = fpdecimal.Zero
		for n := 0; n < b.N; n++ {
			s = x.Add(y)
		}
	})

	b.Run("add_x100", func(b *testing.B) {
		s = fpdecimal.Zero
		for n := 0; n < b.N; n++ {
			for i := 0; i < 100; i++ {
				s = x.Add(y)
			}
		}
	})

	if s == fpdecimal.Zero {
		b.Error()
	}
}

func TestDecimalMemoryLayout(t *testing.T) {
	a, _ := fpdecimal.FromString("-1000.123")
	if v := unsafe.Sizeof(a); v != 8 {
		t.Error(a, v)
	}
}

func TestDecimal_Compare(t *testing.T) {
	a, _ := fpdecimal.FromString("1.123")

	if b, _ := fpdecimal.FromString("1.122"); a.Compare(b) != 1 {
		t.Error(a, ">", b)
	}
	if b, _ := fpdecimal.FromString("1.124"); a.Compare(b) != -1 {
		t.Error(a, "<", b)
	}
	if b, _ := fpdecimal.FromString("1.123"); a.Compare(b) != 0 {
		t.Error(a, "==", b)
	}
}

func TestSetFractionDigits(t *testing.T) {
	defer func() { fpdecimal.FractionDigits = 3 }()

	t.Run("default 3", func(t *testing.T) {
		if a, err := fpdecimal.FromString("1.123"); a.String() != "1.123" || err != nil {
			t.Error("SetFractionDigits", a.String())
		}
	})

	t.Run("5", func(t *testing.T) {
		fpdecimal.FractionDigits = 5
		if a, err := fpdecimal.FromString("1.123456"); a.String() != "1.12345" || err != nil {
			t.Error("SetFractionDigits 5", a.String())
		}
	})

	t.Run("10", func(t *testing.T) {
		fpdecimal.FractionDigits = 10
		if a, err := fpdecimal.FromString("1.12345678910"); a.String() != "1.1234567891" || err != nil {
			t.Error("SetFractionDigits 10", a.String())
		}
	})
}
