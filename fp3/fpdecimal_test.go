package fp3_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"testing"
	"unsafe"

	fp "github.com/nikolaydubina/fpdecimal/fp3"
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
		fa := fp.FromIntScaled(a)
		fb := fp.FromIntScaled(b)

		v := []bool{
			// sum commutativity
			fa.Add(fb) == fb.Add(fa),

			// sum associativity
			fp.Zero.Add(fa).Add(fb).Add(fa) == fp.Zero.Add(fb).Add(fa).Add(fa),

			// sum zero
			fa == fa.Add(fb).Sub(fb),
			fa == fa.Sub(fb).Add(fb),
			fp.Zero == fp.Zero.Add(fa).Sub(fa),

			// product identity
			fa == fa.Mul(fp.FromInt(1)),

			// product zero
			fp.Zero == fa.Mul(fp.FromInt(0)),

			// match number
			(a == b) == (fa == fb),
			(a == b) == fa.Equal(fb),
			a < b == fa.LessThan(fb),
			a > b == fa.GreaterThan(fb),
			a <= b == fa.LessThanOrEqual(fb),
			a >= b == fa.GreaterThanOrEqual(fb),

			// match number convert
			fp.FromIntScaled(a+b) == fa.Add(fb),
			fp.FromIntScaled(a-b) == fa.Sub(fb),
		}
		for i, q := range v {
			if !q {
				t.Error(i, a, b, fa, fb)
			}
		}

		if b != 0 {
			pdiv := fa.Div(fp.FromInt(b))
			p, r := fa.DivMod(fp.FromInt(b))
			if p != pdiv {
				t.Error(p, pdiv)
			}
			if p != fp.FromIntScaled(a/b) {
				t.Error(a, b, p, a/b)
			}
			if fr := fp.FromIntScaled(a % b); r != fr {
				t.Error("Mod", "a", a, "b", b, "got", fr, "want", r)
			}
		}
	})
}

func FuzzParse_StringSameAsFloat(f *testing.F) {
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

		v, err := fp.FromString(s)
		if err != nil {
			t.Error(err)
		}

		if s == "-0.000" || s == "0.000" || rs == 0 || rs == -0 || (rs > -0.001 && rs < 0.001) {
			if v.String() != "0" {
				t.Errorf("s('0') != Decimal.String(%#v) of fp3(%#v) float32(%#v) .3f-float32(%#v)", v.String(), v, r, s)
			}
			return
		}

		if s, fs := strconv.FormatFloat(rs, 'f', -1, 64), v.String(); s != fs {
			t.Error(s, fs, r, v)
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
		v, err := fp.FromString(s)
		if err != nil {
			if v != fp.Zero {
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
		a := fp.FromFloat(v)

		if float32(v) != a.Float32() {
			t.Error("a", a, "a.f32", a.Float32(), "f32.v", float32(v))
		}

		if v != a.Float64() {
			t.Error("a", a, "a.f32", a.Float32(), "v", v)
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
		a := fp.FromFloat(v)

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
	var s fp.Decimal
	var err error

	b.Run("fromString", func(b *testing.B) {
		for _, tc := range floatsForTests {
			b.ResetTimer()
			b.Run(tc.name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s, err = fp.FromString(tc.vals[n%len(tc.vals)])
					if err != nil || s == fp.Zero {
						b.Error(s, err)
					}
				}
			})
		}
	})

	b.Run("UnmarshalJSON", func(b *testing.B) {
		for _, tc := range floatsForTests {
			var vals [][]byte
			for i := range tc.vals {
				vals = append(vals, []byte(tc.vals[i]))
			}

			b.ResetTimer()
			b.Run(tc.name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					if err = s.UnmarshalJSON(vals[n%len(vals)]); err != nil || s == fp.Zero {
						b.Error(s, err)
					}
				}
			})
		}
	})
}

func BenchmarkPrint(b *testing.B) {
	var s string
	for _, tc := range floatsForTests {
		tests := make([]fp.Decimal, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := fp.FromString(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
			tests = append(tests, fp.Zero.Sub(v))
		}

		b.Run("String", func(b *testing.B) {
			b.ResetTimer()
			b.Run(tc.name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s = tests[n%len(tc.vals)].String()
					if s == "" {
						b.Error("empty str")
					}
				}
			})
		})

		b.Run("Marshal", func(b *testing.B) {
			b.ResetTimer()
			b.Run(tc.name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s = tests[n%len(tc.vals)].String()
					if s == "" {
						b.Error("empty str")
					}
				}
			})
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type MyType struct {
		TeslaStockPrice fp.Decimal `json:"tesla-stock-price"`
	}

	tests := []struct {
		json string
		v    fp.Decimal
		s    string
	}{
		{
			json: `{"tesla-stock-price": 9000.001}`,
			v:    fp.FromFloat(9000.001),
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
		TeslaStockPrice fp.Decimal `json:"tesla-stock-price"`
	}

	t.Run("when nil struct, then error", func(t *testing.T) {
		var v *MyType
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), v)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when nil value, then error", func(t *testing.T) {
		var v *fp.Decimal
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), &v)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when nil const of type, then error", func(t *testing.T) {
		err := json.Unmarshal([]byte(`{"tesla-stock-price": 9000.001}`), (*fp.Decimal)(nil))
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
		e := MyType{fp.FromIntScaled(9000001)}
		if v != e {
			t.Error(v)
		}
	})
}

func FuzzJSON(f *testing.F) {
	type MyType struct {
		A fp.Decimal `json:"a"`
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
		rs, _ := strconv.ParseFloat(b, 64)
		s := `{"a":` + strconv.FormatFloat(rs, 'f', -1, 64) + `}`

		var x MyType
		err := json.Unmarshal([]byte(s), &x)
		if err != nil {
			t.Error(err, s)
		}

		if b == "-0.000" || b == "0.000" || rs == 0 || rs == -0 || (rs > 0.001 && rs < 0.001) {
			if x.A.String() != "0" {
				t.Error(b, x)
			}
			return
		}

		if a := x.A.String(); a != strconv.FormatFloat(rs, 'f', -1, 64) {
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
	var BuySP500Price = fp.FromInt(9000)

	input := []byte(`{"sp500": 9000.023}`)

	type Stocks struct {
		SP500 fp.Decimal `json:"sp500"`
	}
	var v Stocks
	if err := json.Unmarshal(input, &v); err != nil {
		log.Fatal(err)
	}

	var amountToBuy fp.Decimal
	if v.SP500.GreaterThan(BuySP500Price) {
		amountToBuy = amountToBuy.Add(v.SP500.Mul(fp.FromInt(2)))
	}

	fmt.Println(amountToBuy)
	// Output: 18000.046
}

func ExampleDecimal_skip_whole_fraction() {
	v, _ := fp.FromString("1013.0000")
	fmt.Println(v)
	// Output: 1013
}

func ExampleDecimal_skip_trailing_zeros() {
	v, _ := fp.FromString("102.0020")
	fmt.Println(v)
	// Output: 102.002
}

func ExampleDecimal_Div() {
	x, _ := fp.FromString("1.000")
	p := x.Div(fp.FromInt(3))
	fmt.Print(p)
	// Output: 0.333
}

func ExampleDecimal_Div_whole() {
	x, _ := fp.FromString("1.000")
	p := x.Div(fp.FromInt(5))
	fmt.Print(p)
	// Output: 0.2
}

func ExampleDecimal_Mod() {
	x, _ := fp.FromString("1.000")
	m := x.Mod(fp.FromInt(3))
	fmt.Print(m)
	// Output: 0.001
}

func ExampleDecimal_DivMod() {
	x, _ := fp.FromString("1.000")
	p, m := x.DivMod(fp.FromInt(3))
	fmt.Print(p, m)
	// Output: 0.333 0.001
}

func ExampleDecimal_DivMod_whole() {
	x, _ := fp.FromString("1.000")
	p, m := x.DivMod(fp.FromInt(5))
	fmt.Print(p, m)
	// Output: 0.2 0
}

func ExampleFromInt_uint8() {
	var x uint8 = 100
	v := fp.FromInt(x)
	fmt.Print(v)
	// Output: 100
}

func ExampleFromInt_int8() {
	var x int8 = -100
	v := fp.FromInt(x)
	fmt.Print(v)
	// Output: -100
}

func ExampleFromInt_int() {
	var x int = -100
	v := fp.FromInt(x)
	fmt.Print(v)
	// Output: -100
}

func ExampleFromInt_uint() {
	var x uint = 100
	v := fp.FromInt(x)
	fmt.Print(v)
	// Output: 100
}

func ExampleMin() {
	min := fp.Min(fp.FromInt(100), fp.FromFloat(0.999), fp.FromFloat(100.001))
	fmt.Print(min)
	// Output: 0.999
}

func ExampleMin_empty() {
	defer func() { fmt.Print(recover()) }()
	fp.Min()
	// Output: min of empty set is undefined
}

func ExampleMax() {
	max := fp.Max(fp.FromInt(100), fp.FromFloat(0.999), fp.FromFloat(100.001))
	fmt.Print(max)
	// Output: 100.001
}

func ExampleMax_empty() {
	defer func() { fmt.Print(recover()) }()
	fp.Max()
	// Output: max of empty set is undefined
}

func BenchmarkArithmetic(b *testing.B) {
	x, _ := fp.FromString("251.231")
	y, _ := fp.FromString("21231.001")

	var s, u fp.Decimal

	u = fp.FromInt(5)

	b.Run("add", func(b *testing.B) {
		s = fp.Zero
		for n := 0; n < b.N; n++ {
			s = x.Add(y)
		}
	})

	b.Run("div", func(b *testing.B) {
		s = fp.Zero
		for n := 0; n < b.N; n++ {
			s = x.Div(y)
		}
	})

	b.Run("divmod", func(b *testing.B) {
		s = fp.Zero
		u = fp.Zero
		for n := 0; n < b.N; n++ {
			s, u = x.DivMod(y)
		}
	})

	if s == fp.Zero || u == fp.Zero {
		b.Error()
	}
}

func TestDecimalMemoryLayout(t *testing.T) {
	a, _ := fp.FromString("-1000.123")
	if v := unsafe.Sizeof(a); v != 8 {
		t.Error(a, v)
	}
}

func TestDecimal_Compare(t *testing.T) {
	a, _ := fp.FromString("1.123")

	if b, _ := fp.FromString("1.122"); a.Compare(b) != 1 {
		t.Error(a, ">", b)
	}
	if b, _ := fp.FromString("1.124"); a.Compare(b) != -1 {
		t.Error(a, "<", b)
	}
	if b, _ := fp.FromString("1.123"); a.Compare(b) != 0 {
		t.Error(a, "==", b)
	}
}
