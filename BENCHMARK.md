# Router

#### Depth order
```
a: AddPattern{ Topic: "order" }
b: AddPattern{ Topic: "order", Cmd: "create" }
c: AddPattern{ Topic: "order", Cmd: "create", Type: 3 }

ActPattern{ Topic: "order", Cmd: "create" } // b Matched
ActPattern{ Topic: "order" } // a Matched
ActPattern{ Topic: "order", Type: 3 } // c Matched
```

#### Insertion order
```
a: AddPattern{ Topic: "order" }
b: AddPattern{ Topic: "order", Cmd: "create" }
c: AddPattern{ Topic: "order", Cmd: "create", Type: 3 }

ActPattern{ Topic: "order", Cmd: "create" } // a Matched
ActPattern{ Topic: "order" } // a Matched
ActPattern{ Topic: "order", Type: 3 } // a Matched
```

## Benchmark
- `Lookup` on 10000 Pattern
- `List` on 10000 Pattern
- `Add` with struct of depth 4
```
BenchmarkLookupWeightDepth7-4             200000              7236 ns/op
BenchmarkLookupWeightDepth6-4              10000            139158 ns/op
BenchmarkLookupWeightDepth5-4               5000            281219 ns/op
BenchmarkLookupWeightDepth4-4               2000            705551 ns/op
BenchmarkLookupWeightDepth3-4               2000            557297 ns/op
BenchmarkLookupWeightDepth2-4               2000            690949 ns/op
BenchmarkLookupWeightDepth1-4               2000            682166 ns/op
BenchmarkListDepth100000-4                   500           2504608 ns/op
BenchmarkAddDepth-4                        10000            128326 ns/op
BenchmarkLookupWeightInsertion7-4         200000              7424 ns/op
BenchmarkLookupWeightInsertion6-4         200000              7020 ns/op
BenchmarkLookupWeightInsertion5-4         200000              6845 ns/op
BenchmarkLookupWeightInsertion4-4         200000              6480 ns/op
BenchmarkLookupWeightInsertion3-4         200000              6355 ns/op
BenchmarkLookupWeightInsertion2-4         200000              5895 ns/op
BenchmarkLookupWeightInsertion1-4           3000            468402 ns/op
BenchmarkListInsertion10000-4                500            2627245 ns/op
BenchmarkAddInsertion-4                    10000            734603 ns/op
PASS
```