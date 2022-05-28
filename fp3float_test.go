package fpfloat_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/nikolaydubina/fpfloat"
)

func FuzzFP3Float_ParseStringSameAsFloat(f *testing.F) {
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

		v, err := fpfloat.FP3FloatFromString(s)
		if err != nil {
			t.Errorf(err.Error())
		}

		if s == "-0.000" || s == "0.000" || r == 0 || r == -0 {
			if v.String() != "0" {
				t.Errorf("s('0') != FP3Float.String(%#v) of fpfloat(%#v) float32(%#v) .3f-float32(%#v)", v.String(), v, r, s)
			}
			return
		}

		if s != v.String() {
			t.Errorf("s(%#v) != FP3Float.String(%#v) of fpfloat(%#v) float32(%#v)", s, v.String(), v, r)
		}
	})
}

func FuzzFP3Float_ParseStringRaw(f *testing.F) {
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
		v, err := fpfloat.FP3FloatFromString(s)
		if err != nil {
			if v != 0 {
				t.Errorf("has to be 0 on error")
			}
			return
		}
	})
}

func FuzzFP3Float_ToType(f *testing.F) {
	tests := []int64{
		0,
		1,
		1000,
		1001,
		123456,
	}
	for _, tc := range tests {
		f.Add(tc)
		f.Add(-tc)
	}
	f.Fuzz(func(t *testing.T, v int64) {
		a := fpfloat.FP3Float(v)

		if b := float32(v) / 1000; b != a.Float32() {
			t.Fatalf("fp3float(%#v) != float32(%#v)", a, b)
		}

		if b := float64(v) / 1000; b != a.Float64() {
			t.Fatalf("fp3float(%#v) != float32(%#v)", a, b)
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

func BenchmarkFP3FloatFromString(b *testing.B) {
	var s fpfloat.FP3Float
	var err error
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s, err = fpfloat.FP3FloatFromString(tc.vals[n%len(tc.vals)])
				if err != nil || s == 0 {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkFP3Float_String(b *testing.B) {
	var s string
	for _, tc := range testsFloats {
		tests := make([]fpfloat.FP3Float, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := fpfloat.FP3FloatFromString(q)
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

func TestFP3Float_UnmarshalJSON(t *testing.T) {
	type MyType struct {
		TeslaStockPrice fpfloat.FP3Float `json:"tesla-stock-price"`
	}

	tests := []struct {
		json string
		v    fpfloat.FP3Float
		s    string
	}{
		{
			json: `{"tesla-stock-price": 9000.001}`,
			v:    fpfloat.FP3Float(9000001),
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

func FuzzFP3Float_UnmarshalJSON(f *testing.F) {
	type MyType struct {
		A fpfloat.FP3Float `json:"a"`
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

func ExampleFP3Float() {
	var sellP500Price = fpfloat.FP3FloatFromFloat(-100.0)
	if sellP500Price == 0 {
		log.Fatal("asdf")
	}

	var BuySP500Price = fpfloat.FP3FloatFromInt(9000)

	input := []byte(`{"sp500": 9000.023}`)

	type Stocks struct {
		SP500 fpfloat.FP3Float `json:"sp500"`
	}
	var v Stocks
	if err := json.Unmarshal(input, &v); err != nil {
		log.Fatal(err)
	}

	var amountToBuy fpfloat.FP3Float
	if v.SP500 > BuySP500Price {
		amountToBuy += v.SP500 * 2
	}

	fmt.Println(amountToBuy)
	// Output: 18000.046
}
