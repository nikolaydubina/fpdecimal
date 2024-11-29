# Fixed-Point Decimals

> [!CAUTION]
> DO NOT initialize package level constants if you are modifying `FractionDigits`. Variables may be initialized with different scaling factor depending on where `init()` call happens.

> To use in money, look at [github.com/nikolaydubina/fpmoney](https://github.com/nikolaydubina/fpmoney)

> _Be Precise. Using floats to represent currency is almost criminal. — Robert.C.Martin, "Clean Code" p.301_

[![codecov](https://codecov.io/gh/nikolaydubina/fpdecimal/branch/main/graph/badge.svg?token=0pf0P5qloX)](https://codecov.io/gh/nikolaydubina/fpdecimal)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/fpdecimal.svg)](https://pkg.go.dev/github.com/nikolaydubina/fpdecimal)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#financial)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/fpdecimal)](https://goreportcard.com/report/github.com/nikolaydubina/fpdecimal)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/fpdecimal/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/fpdecimal)

* `int64` inside
* does not use `float` neither in parsing nor printing
* as fast as `int64` in parsing, printing, arithmetics — 3x faser `float`, 20x faster [shopspring/decimal](https://github.com/shopspring/decimal), 30x faster `fmt`
* zero-overhead
* preventing error-prone fixed-point arithmetics
* Fuzz tests, Benchmarks
* JSON
* 200LOC

```go
import fp "github.com/nikolaydubina/fpdecimal"

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
$ go test -bench=BenchmarkParse -benchtime=5s -benchmem .
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkParse/fromString/small-10                             534307098            11.36 ns/op           0 B/op           0 allocs/op
BenchmarkParse/fromString/large-10                             254741558            23.42 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/small-10                          816873427             7.32 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/large-10                          272173255            22.16 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_Atoi/small-10                      1000000000             4.87 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_Atoi/large-10                       420536834            14.31 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/small/int32-10             561137575            10.67 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/small/int64-10             564200026            10.64 ns/op           0 B/op           0 allocs/op
BenchmarkParse_int_strconv_ParseInt/large/int64-10             219626983            27.17 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/small/float32-10       345666214            17.36 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/small/float64-10       339620222            17.68 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/large/float32-10       128824344            46.68 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_strconv_ParseFloat/large/float64-10       128140617            46.89 ns/op           0 B/op           0 allocs/op
BenchmarkParse_float_fmt_Sscanf/small-10                        21202892           281.6  ns/op          69 B/op           2 allocs/op
BenchmarkParse_float_fmt_Sscanf/large-10                        10074237           599.2  ns/op          88 B/op           3 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    116.249s
```

Print
```
$ go test -bench=BenchmarkPrint -benchtime=5s -benchmem .
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkPrint/small-10                                      191982066            31.24 ns/op           8 B/op           1 allocs/op
BenchmarkPrint/large-10                                      150874335            39.89 ns/op          24 B/op           1 allocs/op
BenchmarkPrint_int_strconv_Itoa/small-10                     446302868            13.39 ns/op           3 B/op           0 allocs/op
BenchmarkPrint_int_strconv_Itoa/large-10                     237484774            25.20 ns/op          18 B/op           1 allocs/op
BenchmarkPrint_int_strconv_FormatInt/small-10                444861666            13.70 ns/op           3 B/op           0 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/small/float32-10     55003357           104.2  ns/op          31 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/small/float64-10     43565430           137.4  ns/op          31 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/large/float32-10     64069650            92.07 ns/op          48 B/op           2 allocs/op
BenchmarkPrint_float_strconv_FormatFloat/large/float64-10     68441746            87.36 ns/op          48 B/op           2 allocs/op
BenchmarkPrint_float_fmt_Sprintf/small-10                     46503666           127.7  ns/op          16 B/op           2 allocs/op
BenchmarkPrint_float_fmt_Sprintf/large-10                     51764224           115.8  ns/op          28 B/op           2 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    79.192s
```

Arithmetics
```
$ go test -bench=BenchmarkArithmetic -benchtime=5s -benchmem .
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/fpdecimal
BenchmarkArithmetic/add-10                   1000000000             0.316 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic/div-10                   1000000000             0.950 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic/divmod-10                1000000000             1.890 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/add-10             1000000000             0.314 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/div-10             1000000000             0.316 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/divmod-10          1000000000             1.261 ns/op           0 B/op           0 allocs/op
BenchmarkArithmetic_int64/mod-10             1000000000             0.628 ns/op           0 B/op           0 allocs/op
PASS
ok      github.com/nikolaydubina/fpdecimal    6.721s
```

## References

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
var b fpdecimal.FromInt(1000)

// does not compile
a + b
```

However, untyped constants will be resolved to underlying type `int64` and will be allowed.  
```go
const a 10000
var b fpdecimal.FromInt(1000)

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

## Appendix E: Print into destination

To avoid mallocs, it is advantageous to print formatted value to pre-allocated destination.
Similarly, to `strconv.AppendInt`, we provide `AppendFixedPointDecimal`.
This is utilized in `github.com/nikolaydubina/fpmoney` package.

```
BenchmarkFixedPointDecimalToString/small-10     28522474         35.43 ns/op       24 B/op        1 allocs/op
BenchmarkFixedPointDecimalToString/large-10     36883687         32.32 ns/op       24 B/op        1 allocs/op
BenchmarkAppendFixedPointDecimal/small-10       38105520         30.51 ns/op      117 B/op        0 allocs/op
BenchmarkAppendFixedPointDecimal/large-10       55147478         29.52 ns/op      119 B/op        0 allocs/op
```

## Appendix F: DivMod notation

In early versions, `Div` and `Mul` operated on `int` and `Div` returned remainder.
As recommended by @vanodevium and more in line with other common libraries, notation is changed.
Bellow is survey as of 2023-05-18.

Go, https://pkg.go.dev/math/big
```go
func (z *Int) Div(x, y *Int) *Int
func (z *Int) DivMod(x, y, m *Int) (*Int, *Int)
func (z *Int) Mod(x, y *Int) *Int
```

Go, github.com/shopspring/decimal
```go
func (d Decimal) Div(d2 Decimal) Decimal
// X no DivMod
func (d Decimal) Mod(d2 Decimal) Decimal
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal
```

Python, https://docs.python.org/3/library/decimal.html
```python
divide(x, y) number
divide_int(x, y) number // truncates
divmod(x, y) number
remainder(x, y) number
```

Pytorch, https://pytorch.org/docs/stable/generated/torch.div.html
```python
torch.div(input, other, *, rounding_mode=None, out=None) → [Tensor] // discards remainder
torch.remainder(input, other, *, out=None) → [Tensor] // remainder
```

numpy, https://numpy.org/doc/stable/reference/generated/numpy.divmod.html
```python
np.divmod(x, y) (number, number) // is equivalent to (x // y, x % y
np.mod(x, y) number
np.remainder(x, y) number
np.divide(x, y) number
np.true_divide(x, y) number // same as divide
np.floor_divide(x, y) number // rounding down
```

## Appendix G: generics switch for decimal counting

Go does not support numerics in templates. However, defining multiple types each associated with specific number of decimals and passing them to functions and defining constraint as union of these types — is an attractive option.
This does not work well since Go does not support switch case (casting generic) back to integer well.

## Appendix H: `string` vs `[]byte` in interface

The typical usage of parsing number is through some JSON or other mechanism. Those APIs are dealing with `[]byte`.
Now, conversion from `[]byte` to `string` requires to copy data, since `string` is immutable.
To improve performance, we are using `[]byte` in signatures.

Using `string`
```
BenchmarkParse/fromString/small-10                 831217767             7.07 ns/op           0 B/op           0 allocs/op
BenchmarkParse/fromString/large-10                 275009497            21.79 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/small-10              553035127            10.98 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/large-10              248815030            24.14 ns/op           0 B/op           0 allocs/op
```

Using `[]byte`
```
BenchmarkParse/fromString/small-10                 523937236            11.32 ns/op           0 B/op           0 allocs/op
BenchmarkParse/fromString/large-10                 257542226            23.23 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/small-10              809793006             7.31 ns/op           0 B/op           0 allocs/op
BenchmarkParse/UnmarshalJSON/large-10              272087984            22.04 ns/op           0 B/op           0 allocs/op
```
