# 🐣 Fixed-Point Floats (3FPFloat)

> _When you have small and simple float-like numbers. Precise and Fast. Perfect for money._

* 100LOC
* Fuzz tests
* 100% coverage
* JSON encoding/decoding
* Arithmetic operators
* As fast as integers (parsing, formatting, operations)
* 40x faster parsing than float
* 2x faster printing than float
* 20x faster than [shopspring/decimal](https://github.com/shopspring/decimal)

```go
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
```

### Implementation

Parsing and Printing is expensive operation and requires a lot of code.
However, if you know that your numbers are always small and simple and you do not care or do not permit lots of fractions like `-1234.567`, then parsing and printing can be greatly simplified.
Code is heavily influenced by hot-path from Go core `strconv` package.

### Benchmarks

Parse
```
nikolaydubina@Nikolays-MacBook-Pro fpfloat % go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpfloat
BenchmarkFP3FloatFromString/small-10          	            855162544	         6.933 ns/op	   0 B/op	       0 allocs/op
BenchmarkFP3FloatFromString/large-10          	            275428832	        21.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt_strconvAtoi/small-10        	            1000000000	         4.847 ns/op	   0 B/op	       0 allocs/op
BenchmarkParseInt_strconvAtoi/large-10        	            425393912	        14.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt_strconvParseInt/small/int32-10         	558908575	        10.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt_strconvParseInt/small/int64-10         	565205856	        10.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt_strconvParseInt/large/int64-10         	218402264	        27.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat_strconvParseFloat/small/float32-10   	343771328	        17.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat_strconvParseFloat/small/float64-10   	335683317	        17.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat_strconvParseFloat/large/float32-10   	128418057	        46.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat_strconvParseFloat/large/float64-10   	127516434	        46.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat_fmtSscanf/small-10                   	20711220	       289.1 ns/op	      69 B/op	       2 allocs/op
BenchmarkParseFloat_fmtSscanf/large-10                   	 9629078	       623.5 ns/op	      88 B/op	       3 allocs/op
PASS
ok  	github.com/nikolaydubina/fpfloat	175.917s
```

Format
```
nikolaydubina@Nikolays-MacBook-Pro fpfloat % go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpfloat
BenchmarkFP3Float_String/small-10             	            152158830	        39.09 ns/op	      10 B/op	       1 allocs/op
BenchmarkFP3Float_String/large-10             	            95981250	        63.11 ns/op	      48 B/op	       2 allocs/op
BenchmarkStringInt_strconvItoa/small-10       	            709914073	         8.556 ns/op	   1 B/op	       0 allocs/op
BenchmarkStringInt_strconvItoa/large-10       	            235245307	        25.64 ns/op	      16 B/op	       1 allocs/op
BenchmarkStringInt_strconvFormatInt/small-10             	705413114	         8.522 ns/op	   1 B/op	       0 allocs/op
BenchmarkStringFloat_strconvFormatFloat/small/float32-10 	51060498	       117.7 ns/op	      31 B/op	       2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/small/float64-10 	40995589	       146.3 ns/op	      31 B/op	       2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/large/float32-10 	60731203	        99.07 ns/op	      48 B/op	       2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/large/float64-10 	61169353	        96.83 ns/op	      48 B/op	       2 allocs/op
BenchmarkStringFloat_fmtSprintf/small-10                 	42804757	       138.5 ns/op	      16 B/op	       2 allocs/op
BenchmarkStringFloat_fmtSprintf/large-10                 	47333761	       126.2 ns/op	      28 B/op	       2 allocs/op
PASS
ok  	github.com/nikolaydubina/fpfloat	175.917s
```

### Future Work

- Adding wrapper into a struct to block arithmetic operations (+benchmarks encoding, method calls overhead, long chains of calls)
- Separate repo to benchmark other open source versions (decimals) with same benchmark tests as in this repo (do not include other modules here just for benchmarking)
- Linter to warn about using constants in expressions containing `fpfloat` type.

### References

- [Fixed-Point Arithmetic Wiki](https://en.wikipedia.org/wiki/Fixed-point_arithmetic)
- [shopspring/decimal](https://github.com/shopspring/decimal)

### Appendix A: Comparison to other libraries

- https://github.com/shopspring/decimal solves arbitrary precision, fpfloat solves only simple small floats

### Appendix B: Benchmarking https://github.com/shopspring/decimal (2022-05-28)
```
nikolaydubina@Nikolays-MacBook-Pro decimal % go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/shopspring/decimal
BenchmarkNewFromFloatWithExponent-10                    	59701516	        97.07 ns/op	     106 B/op	       4 allocs/op
BenchmarkNewFromFloat-10                                	14771503	       410.3 ns/op	      67 B/op	       2 allocs/op
BenchmarkNewFromStringFloat-10                          	16246342	       375.2 ns/op	     175 B/op	       5 allocs/op
Benchmark_FloorFast-10                                  	1000000000	         2.143 ns/op	       0 B/op	       0 allocs/op
Benchmark_FloorRegular-10                               	53857244	       106.3 ns/op	     112 B/op	       6 allocs/op
Benchmark_DivideOriginal-10                             	       7	 715322768 ns/op	737406446 B/op	30652495 allocs/op
Benchmark_DivideNew-10                                  	      22	 262893689 ns/op	308046721 B/op	12054905 allocs/op
BenchmarkDecimal_RoundCash_Five-10                      	 9311530	       636.5 ns/op	     616 B/op	      28 allocs/op
Benchmark_Cmp-10                                        	      44	 133191579 ns/op	      24 B/op	       1 allocs/op
Benchmark_decimal_Decimal_Add_different_precision-10    	31561636	       176.6 ns/op	     280 B/op	       9 allocs/op
Benchmark_decimal_Decimal_Sub_different_precision-10    	36892767	       164.4 ns/op	     240 B/op	       9 allocs/op
Benchmark_decimal_Decimal_Add_same_precision-10         	134831919	        44.96 ns/op	      80 B/op	       2 allocs/op
Benchmark_decimal_Decimal_Sub_same_precision-10         	134902627	        43.61 ns/op	      80 B/op	       2 allocs/op
BenchmarkDecimal_IsInteger-10                           	92543083	        66.51 ns/op	       8 B/op	       1 allocs/op
BenchmarkDecimal_NewFromString-10                       	  827455	      7382 ns/op	    3525 B/op	     216 allocs/op
BenchmarkDecimal_NewFromString_large_number-10          	  212538	     28836 ns/op	   16820 B/op	     360 allocs/op
BenchmarkDecimal_ExpHullAbraham-10                      	   10000	    572091 ns/op	  486628 B/op	     568 allocs/op
BenchmarkDecimal_ExpTaylor-10                           	   26343	    222915 ns/op	  431226 B/op	    3172 allocs/op
PASS
ok  	github.com/shopspring/decimal	123.541s
```

### Appendix C: Why this is good fit for money?

There are only ~200 currencies in the world.
All currencies have at most 3 decimal digits, thus it is sufficient to handle 3 decimal fractions.
Next, currencies without decimal digits are typically 1000x larger than dollar, but even then maximum number that fits into `int64` (without 3 decimal fractions) is `9 223 372 036 854 775.807` which is ~9 quadrillion. This should be enough for most operations with money.

### Appendix D: Is it safe to use arithmetic operators in Go?

Sort of. Operations with defined types (variables) will fail.
```go
var a int64
var b fpfloat.FP3FloatFromInt(1000)

// does not compile
a + b
```

However, untyped constants will be resolved to underlying type `int64` and will be allowed.  
```go
const a 10000
var b fpfloat.FP3FloatFromInt(1000)

// compiles
a + b

// also compiles
b - 42

// this one too
b *= 23
```

Is this a problem? For multiplication and division - no. For addition substraction - yes, it can be. You have to be careful and remind yourself that constants would be reduced 1000x. This may be addressed at compile time by providing linter. This can be also addressed by wrapping into struct and defining methods.