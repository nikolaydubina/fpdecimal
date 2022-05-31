package fpdecimal_test

import (
	"testing"

	"github.com/nikolaydubina/fpdecimal"
)

func FuzzParseFixedPointDecimal(f *testing.F) {
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
		v, err := fpdecimal.ParseFixedPointDecimal(s, 3)
		if err != nil {
			if v != 0 {
				t.Errorf("has to be 0 on error")
			}
			if err.Error() == "" {
				t.Error(err, err.Error())
			}
			return
		}
	})
}
