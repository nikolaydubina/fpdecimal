package fpfloat_test

import (
	"fmt"
	"strconv"
	"testing"
)

var testsInts = []struct {
	name string
	vals []string
}{
	{
		name: "small",
		vals: []string{
			"123456",
			"0123",
			"0012",
			"0001",
			"0982",
			"0101",
			"10",
			"11",
			"1",
		},
	},
	{
		name: "large",
		vals: []string{
			"123123123112312",
			"5341320482340234",
		},
	},
}

func BenchmarkParseInt_strconvAtoi(b *testing.B) {
	var s int
	var err error
	for _, tc := range testsInts {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s, err = strconv.Atoi(tc.vals[n%len(tc.vals)])
				if err != nil || s == 0 {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkStringInt_strconvItoa(b *testing.B) {
	var s string
	for _, tc := range testsInts {
		tests := make([]int, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := strconv.Atoi(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
		}

		b.ResetTimer()
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s = strconv.Itoa(tests[n%len(tests)])
				if s == "" {
					b.Error("empty str")
				}
			}
		})
	}
}

func BenchmarkParseInt_strconvParseInt(b *testing.B) {
	var err error
	for _, tc := range testsInts {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("int32", func(b *testing.B) {
				if tc.name == "large" {
					b.Skip()
				}
				var s int64
				for n := 0; n < b.N; n++ {
					s, err = strconv.ParseInt(tc.vals[n%len(tc.vals)], 10, 32)
					if err != nil || s == 0 {
						b.Error(s, err)
					}
				}
			})

			b.Run("int64", func(b *testing.B) {
				var s int64
				for n := 0; n < b.N; n++ {
					s, err = strconv.ParseInt(tc.vals[n%len(tc.vals)], 10, 64)
					if err != nil || s == 0 {
						b.Error(s, err)
					}
				}
			})
		})
	}
}

func BenchmarkStringInt_strconvFormatInt(b *testing.B) {
	var s string
	for _, tc := range testsInts {
		b.Run(tc.name, func(b *testing.B) {
			if tc.name == "large" {
				b.Skip()
			}

			tests := make([]int64, 0, len(tc.vals))
			for _, q := range tc.vals {
				v, err := strconv.ParseInt(q, 10, 32)
				if err != nil {
					b.Error(err)
				}
				tests = append(tests, v)
			}

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				s = strconv.FormatInt(tests[n%len(tests)], 10)
				if s == "" {
					b.Error("empty str")
				}
			}
		})
	}
}

func BenchmarkParseFloat_strconvParseFloat(b *testing.B) {
	var s float64
	var err error
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("float32", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s, err = strconv.ParseFloat(tc.vals[n%len(tc.vals)], 32)
					if err != nil || s == 0 {
						b.Error(s, err)
					}
				}
			})

			b.Run("float64", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s, err = strconv.ParseFloat(tc.vals[n%len(tc.vals)], 64)
					if err != nil || s == 0 {
						b.Error(s, err)
					}
				}
			})
		})
	}
}

func BenchmarkStringFloat_strconvFormatFloat(b *testing.B) {
	var s string
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			tests := make([]float64, 0, len(tc.vals))
			for _, q := range tc.vals {
				v, err := strconv.ParseFloat(q, 32)
				if err != nil {
					b.Error(err)
				}
				tests = append(tests, v)
			}

			b.ResetTimer()
			b.Run("float32", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s = strconv.FormatFloat(tests[n%len(tests)], 'f', 3, 32)
					if s == "" {
						b.Error("empty str")
					}
				}
			})

			b.Run("float64", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					s = strconv.FormatFloat(tests[n%len(tests)], 'f', 3, 64)
					if s == "" {
						b.Error("empty str")
					}
				}
			})
		})
	}
}

func BenchmarkParseFloat_fmtSscanf(b *testing.B) {
	var s float32
	var err error
	for _, tc := range testsFloats {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err = fmt.Sscanf(tc.vals[n%len(tc.vals)], "%f", &s)
				if err != nil || s == 0 {
					b.Error(s, err)
				}
			}
		})
	}
}

func BenchmarkStringFloat_fmtSprintf(b *testing.B) {
	var s string
	var err error
	for _, tc := range testsFloats {
		tests := make([]float32, 0, len(tc.vals))
		for _, q := range tc.vals {
			var s float32
			_, err = fmt.Sscanf(q, "%f", &s)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, s)
		}

		b.ResetTimer()
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s = fmt.Sprintf("%.3f", tests[n%len(tests)])
				if s == "" {
					b.Error("empty str")
				}
			}
		})
	}
}
