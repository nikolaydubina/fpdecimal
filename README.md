# 🐣 Small Fixed-Point Decimals (FP3Decimal)

> _When you have small and simple float-like numbers. Precise and Fast. Perfect for money._

[![codecov](https://codecov.io/gh/nikolaydubina/fpdecimal/branch/main/graph/badge.svg?token=0pf0P5qloX)](https://codecov.io/gh/nikolaydubina/fpdecimal)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/fpdecimal.svg)](https://pkg.go.dev/github.com/nikolaydubina/fpdecimal)

* 100LOC
* Fuzz tests
* 100% coverage
* JSON encoding/decoding
* Arithmetic operators
* As fast as integers (parsing, formatting, operations)
* 90x faster parsing than float
* 2x faster printing than float
* 20x faster than [shopspring/decimal](https://github.com/shopspring/decimal)

```go
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
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkFP3DecimalFromString/small-10                      827744476            7.138 ns/op           0 B/op           0 allocs/op
BenchmarkFP3DecimalFromString/large-10                      276668296            21.79 ns/op           0 B/op           0 allocs/op
BenchmarkParseInt_strconvAtoi/small-10                     1000000000            4.791 ns/op           0 B/op           0 allocs/op
BenchmarkParseInt_strconvAtoi/large-10                      416969704            14.18 ns/op           0 B/op           0 allocs/op
BenchmarkParseInt_strconvParseInt/small/int32-10            567803484            10.56 ns/op           0 B/op           0 allocs/op
BenchmarkParseInt_strconvParseInt/small/int64-10            567515059            10.56 ns/op           0 B/op           0 allocs/op
BenchmarkParseInt_strconvParseInt/large/int64-10            221833478            27.14 ns/op           0 B/op           0 allocs/op
BenchmarkParseFloat_strconvParseFloat/small/float32-10      349272979            17.34 ns/op           0 B/op           0 allocs/op
BenchmarkParseFloat_strconvParseFloat/small/float64-10      333610484            17.82 ns/op           0 B/op           0 allocs/op
BenchmarkParseFloat_strconvParseFloat/large/float32-10      129024007            46.45 ns/op           0 B/op           0 allocs/op
BenchmarkParseFloat_strconvParseFloat/large/float64-10      128212430            46.79 ns/op           0 B/op           0 allocs/op
BenchmarkParseFloat_fmtSscanf/small-10                       20381784            293.4 ns/op          69 B/op           2 allocs/op
BenchmarkParseFloat_fmtSscanf/large-10                        9484489            629.3 ns/op          88 B/op           3 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    175.518s
```

Format
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkFP3Decimal_String/small-10                         161572239            37.17 ns/op          10 B/op           1 allocs/op
BenchmarkFP3Decimal_String/large-10                          92706448            63.58 ns/op          48 B/op           2 allocs/op
BenchmarkStringInt_strconvItoa/small-10                     729627100            8.327 ns/op           1 B/op           0 allocs/op
BenchmarkStringInt_strconvItoa/large-10                     233921521            25.61 ns/op          16 B/op           1 allocs/op
BenchmarkStringInt_strconvFormatInt/small-10                736678662            8.141 ns/op           1 B/op           0 allocs/op
BenchmarkStringFloat_strconvFormatFloat/small/float32-10     50491785            117.8 ns/op          31 B/op           2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/small/float64-10     40790115            147.4 ns/op          31 B/op           2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/large/float32-10     60102750            99.38 ns/op          48 B/op           2 allocs/op
BenchmarkStringFloat_strconvFormatFloat/large/float64-10     61115224            97.45 ns/op          48 B/op           2 allocs/op
BenchmarkStringFloat_fmtSprintf/small-10                     43199199            138.2 ns/op          16 B/op           2 allocs/op
BenchmarkStringFloat_fmtSprintf/large-10                     47292736            126.2 ns/op          28 B/op           2 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    175.518s
```

### Future Work

- Adding wrapper into a struct to block arithmetic operations (+benchmarks encoding, method calls overhead, long chains of calls)
- Separate repo to benchmark other open source versions (decimals) with same benchmark tests as in this repo (do not include other modules here just for benchmarking)
- Linter to warn about using constants in expressions containing `fpdecimal` type.

### References

- [Fixed-Point Arithmetic Wiki](https://en.wikipedia.org/wiki/Fixed-point_arithmetic)
- [shopspring/decimal](https://github.com/shopspring/decimal)

### Appendix A: Comparison to other libraries

- https://github.com/shopspring/decimal solves arbitrary precision, fpdecimal solves only simple small decimals

### Appendix B: Benchmarking https://github.com/shopspring/decimal (2022-05-28)
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/shopspring/decimal
BenchmarkNewFromFloatWithExponent-10                        59701516         97.07 ns/op         106 B/op           4 allocs/op
BenchmarkNewFromFloat-10                                    14771503         410.3 ns/op          67 B/op           2 allocs/op
BenchmarkNewFromStringFloat-10                              16246342         375.2 ns/op         175 B/op           5 allocs/op
Benchmark_FloorFast-10                                    1000000000         2.143 ns/op           0 B/op           0 allocs/op
Benchmark_FloorRegular-10                                   53857244         106.3 ns/op         112 B/op           6 allocs/op
Benchmark_DivideOriginal-10                                        7     715322768 ns/op   737406446 B/op    30652495 allocs/op
Benchmark_DivideNew-10                                            22     262893689 ns/op   308046721 B/op    12054905 allocs/op
BenchmarkDecimal_RoundCash_Five-10                           9311530         636.5 ns/op         616 B/op          28 allocs/op
Benchmark_Cmp-10                                                  44     133191579 ns/op          24 B/op           1 allocs/op
Benchmark_decimal_Decimal_Add_different_precision-10        31561636         176.6 ns/op         280 B/op           9 allocs/op
Benchmark_decimal_Decimal_Sub_different_precision-10        36892767         164.4 ns/op         240 B/op           9 allocs/op
Benchmark_decimal_Decimal_Add_same_precision-10            134831919         44.96 ns/op          80 B/op           2 allocs/op
Benchmark_decimal_Decimal_Sub_same_precision-10            134902627         43.61 ns/op          80 B/op           2 allocs/op
BenchmarkDecimal_IsInteger-10                               92543083         66.51 ns/op           8 B/op           1 allocs/op
BenchmarkDecimal_NewFromString-10                             827455          7382 ns/op        3525 B/op         216 allocs/op
BenchmarkDecimal_NewFromString_large_number-10                212538         28836 ns/op       16820 B/op         360 allocs/op
BenchmarkDecimal_ExpHullAbraham-10                             10000        572091 ns/op      486628 B/op         568 allocs/op
BenchmarkDecimal_ExpTaylor-10                                  26343        222915 ns/op      431226 B/op        3172 allocs/op
PASS
ok      github.com/shopspring/decimal    123.541sa
```

### Appendix C: Why this is good fit for money?

There are only ~200 currencies in the world.
All currencies have at most 3 decimal digits, thus it is sufficient to handle 3 decimal fractions.
Next, currencies without decimal digits are typically 1000x larger than dollar, but even then maximum number that fits into `int64` (without 3 decimal fractions) is `9 223 372 036 854 775.807` which is ~9 quadrillion. This should be enough for most operations with money.

### Appendix D: Is it safe to use arithmetic operators in Go?

Sort of. Operations with defined types (variables) will fail.
```go
var a int64
var b fpdecimal.FP3DecimalFromInt(1000)

// does not compile
a + b
```

However, untyped constants will be resolved to underlying type `int64` and will be allowed.  
```go
const a 10000
var b fpdecimal.FP3DecimalFromInt(1000)

// compiles
a + b

// also compiles
b - 42

// this one too
b *= 23
```

Is this a problem? 
* For multiplication and division - yes, it can be. You have to be careful not to multiply two `fpdecimal` numbers, since scaling factor will quadruple. Multiplying by constants is ok tho.
* For addition substraction - yes, it can be. You have to be careful and remind yourself that constants would be reduced 1000x.

Both of this can be addressed at compile time by providing linter. This can be also addressed by wrapping into a struct and defining methods.
