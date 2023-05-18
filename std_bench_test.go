package fpdecimal_test

import (
	"fmt"
	"strconv"
	"testing"
)

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

func BenchmarkParse_int_strconv_Atoi(b *testing.B) {
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

func BenchmarkPrint_int_strconv_Itoa(b *testing.B) {
	var s string
	for _, tc := range testsInts {
		tests := make([]int, 0, len(tc.vals))
		for _, q := range tc.vals {
			v, err := strconv.Atoi(q)
			if err != nil {
				b.Error(err)
			}
			tests = append(tests, v)
			tests = append(tests, -v)
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

func BenchmarkParse_int_strconv_ParseInt(b *testing.B) {
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

func BenchmarkPrint_int_strconv_FormatInt(b *testing.B) {
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
				tests = append(tests, -v)
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

func BenchmarkParse_float_strconv_ParseFloat(b *testing.B) {
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

func BenchmarkPrint_float_strconv_FormatFloat(b *testing.B) {
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

func BenchmarkParse_float_fmt_Sscanf(b *testing.B) {
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

func BenchmarkPrint_float_fmt_Sprintf(b *testing.B) {
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

func BenchmarkArithmetic_int64(b *testing.B) {
	var x int64 = 251231
	var y int64 = 21231001

	var s int64

	b.Run("add", func(b *testing.B) {
		s = 0
		for n := 0; n < b.N; n++ {
			s = x + y
		}
	})

	b.Run("div", func(b *testing.B) {
		s = 0
		for n := 0; n < b.N; n++ {
			s = x + y
		}
	})

	b.Run("divmod", func(b *testing.B) {
		s = 0
		for n := 0; n < b.N; n++ {
			s += x / y
			s += x % y
		}
	})

	b.Run("mod", func(b *testing.B) {
		s = 0
		for n := 0; n < b.N; n++ {
			s = x % y
		}
	})

	if s == 0 {
		b.Error()
	}
}
