# 🛫 Small Fixed-Point Decimals (fp3.Decimal)

> _When you have small and simple float-like numbers. Precise and Fast. Perfect for money._

[![codecov](https://codecov.io/gh/nikolaydubina/fpdecimal/branch/main/graph/badge.svg?token=0pf0P5qloX)](https://codecov.io/gh/nikolaydubina/fpdecimal)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/fpdecimal.svg)](https://pkg.go.dev/github.com/nikolaydubina/fpdecimal)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#financial)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/fpdecimal)](https://goreportcard.com/report/github.com/nikolaydubina/fpdecimal)

* ~100LOC
* stores internally as `int64`
* does not use `float64` in parsing nor printing
* Fuzz tests (parsing, printing, arithmetics)
* 100% coverage
* JSON encoding/decoding
* As fast as integers (parsing, printing, arithmetics)
* 3x faster parsing than float
* 2x faster printing than float
* 30x faster parsing than `fmt`
* 20x faster parsing than [shopspring/decimal](https://github.com/shopspring/decimal)
* no overhead for arithmetic operations
* no overhead for memory (same size as `int64`)
* blocking error-prone fixed-point arithmetics

```go
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
if v.SP500.HigherThan(BuySP500Price) {
    amountToBuy = amountToBuy.Add(v.SP500.Mul(2))
}

fmt.Println(amountToBuy)
// Output: 18000.046
```

### Implementation

Parsing and Printing is expensive operation and requires a lot of code.
However, if you know that your numbers are always small and simple and you do not care or do not permit lots of fractions like `-1234.567`, then parsing and printing can be greatly simplified.
Code is heavily influenced by hot-path from Go core `strconv` package.

It is wrapped into struct to prevent bugs:
- block multiplication by `fpdecimal` type, which leads to increase in decimal fractions and loose of precision
- block additions of untyped constants, which leads to errors if you forget to scale by factor

### Benchmarks

Parse
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkParse_FP3Decimal/small-10                             845515756             7.04 ns/op           0 B/op           0 allocs/op
BenchmarkParse_FP3Decimal/large-10                             278560885            21.43 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_Atoi/small-10                      1000000000             4.74 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_Atoi/large-10                       424242687            14.17 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/small/int32-10             566976321            10.65 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/small/int64-10             552894133            10.85 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/large/int64-10             219031276            27.56 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/small/float32-10       344793511            17.43 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/small/float64-10       335880535            17.82 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/large/float32-10       129427171            46.40 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/large/float64-10       128508513            46.75 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_fmt_Sscanf/small-10                        20424795           295.6  ns/op          69 B/op           2 allocs/op
BenchmarkParse_float_fmt_Sscanf/large-10                         9479828           633.9  ns/op          88 B/op           3 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    194.558s
```

Print
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkPrint_FP3Decimal/small-10                            164560935            36.44 ns/op          10 B/op           1 allocs/op
BenchmarkPrint_FP3Decimal/large-10                             95584335            63.17 ns/op          48 B/op           2 allocs/op
BenchmarkPrint_int_strconv_Itoa/small-10                      731089902             8.19 ns/op           1 B/op           0 allocs/op
BenchmarkPrint_int_strconv_Itoa/large-10                      234441306            25.77 ns/op          16 B/op           1 allocs/op
BenchmarkPrint_int_strconv_FormatInt/small-10                 728307549             8.26 ns/op           1 B/op           0 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/small/float32-10      49801364           117.8  ns/op          31 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/small/float64-10      40938864           148.3  ns/op          31 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/large/float32-10      58160480            99.12 ns/op          48 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/large/float64-10      61878582            97.22 ns/op          48 B/op           2 allocs/op
BenchmarkPrint_float_fmt_Sprintf/small-10                      43542469           138.8  ns/op          16 B/op           2 allocs/op
BenchmarkPrint_float_fmt_Sprintf/large-10                      47824404           125.7  ns/op          28 B/op           2 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    194.558s
```

Arithmetics
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkArithmetic_FP3Decimal/add_x1-10           1000000000             0.31 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_FP3Decimal/add_x100-10          181966545            32.75 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/add_x1-10                1000000000             0.31 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/add_x100-10               182298925            32.99 ns/op           0 B/op           0 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    194.558s
```

### References

- [Fixed-Point Arithmetic Wiki](https://en.wikipedia.org/wiki/Fixed-point_arithmetic)
- [shopspring/decimal](https://github.com/shopspring/decimal)

## Appendix A: Comparison to other libraries

- https://github.com/shopspring/decimal solves arbitrary precision, fpdecimal solves only simple small decimals
- https://github.com/Rhymond/go-money solves typed number (currency), decodes through `interface{}` and float64, no precision in decoding, expects encoding to be in cents

## Appendix B: Benchmarking [shopspring/decimal](https://github.com/shopspring/decimal)

`2022-05-28`
```
$ go test -bench=. -benchtime=5s -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/shopspring/decimal
BenchmarkNewFromFloatWithExponent-10                        59701516          97.7 ns/op         106 B/op           4 allocs/op
BenchmarkNewFromFloat-10                                    14771503         410.3 ns/op          67 B/op           2 allocs/op
BenchmarkNewFromStringFloat-10                              16246342         375.2 ns/op         175 B/op           5 allocs/op
Benchmark_FloorFast-10                                    1000000000           2.1 ns/op           0 B/op           0 allocs/op
Benchmark_FloorRegular-10                                   53857244         106.3 ns/op         112 B/op           6 allocs/op
Benchmark_DivideOriginal-10                                        7   715322768   ns/op   737406446 B/op    30652495 allocs/op
Benchmark_DivideNew-10                                            22   262893689   ns/op   308046721 B/op    12054905 allocs/op
BenchmarkDecimal_RoundCash_Five-10                           9311530         636.5 ns/op         616 B/op          28 allocs/op
Benchmark_Cmp-10                                                  44   133191579   ns/op          24 B/op           1 allocs/op
Benchmark_decimal_Decimal_Add_different_precision-10        31561636         176.6 ns/op         280 B/op           9 allocs/op
Benchmark_decimal_Decimal_Sub_different_precision-10        36892767         164.4 ns/op         240 B/op           9 allocs/op
Benchmark_decimal_Decimal_Add_same_precision-10            134831919          44.9 ns/op          80 B/op           2 allocs/op
Benchmark_decimal_Decimal_Sub_same_precision-10            134902627          43.1 ns/op          80 B/op           2 allocs/op
BenchmarkDecimal_IsInteger-10                               92543083          66.1 ns/op           8 B/op           1 allocs/op
BenchmarkDecimal_NewFromString-10                             827455        7382   ns/op        3525 B/op         216 allocs/op
BenchmarkDecimal_NewFromString_large_number-10                212538       28836   ns/op       16820 B/op         360 allocs/op
BenchmarkDecimal_ExpHullAbraham-10                             10000      572091   ns/op      486628 B/op         568 allocs/op
BenchmarkDecimal_ExpTaylor-10                                  26343      222915   ns/op      431226 B/op        3172 allocs/op
PASS
ok      github.com/shopspring/decimal    123.541sa
```

## Appendix C: Why this is good fit for money?

There are only ~200 currencies in the world.
All currencies have at most 3 decimal digits, thus it is sufficient to handle 3 decimal fractions.
Next, currencies without decimal digits are typically 1000x larger than dollar, but even then maximum number that fits into `int64` (without 3 decimal fractions) is `9 223 372 036 854 775.807` which is ~9 quadrillion. This should be enough for most operations with money.

## Appendix D: Is it safe to use arithmetic operators in Go?

Sort of... 

In one of iterations, I did Type Alias, but it required some effort to use it carefully.

Operations with defined types (variables) will fail.
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

Both of this can be addressed at compile time by providing linter.
This can be also addressed by wrapping into a struct and defining methods.
Formed is hard to achieve in Go, due to lack of operator overload and lots of work required to write AST parser.
Later has been implemented in this pacakge, and, as benchmarks show, without any extra memory or calls overhead as compared to `int64`.
