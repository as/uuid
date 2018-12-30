# uuid
simple uuid v4 generator without vacuous error returns

# usage

```
package main

import "github.com/as/uuid"

func main() {
	fmt.Println(uuid.V4())
}
```

# install

`go get github.com/as/uuid`


# questions

- Do I have to check for errors?
> no

- Does this panic when the entropy pool "depletes"?
> no

- How does it work? Does it use `math/random`?
> it uses the AES family of ciphers as a CSPRNG 

- Is it faster than uuid package `X`?
> see benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/as/uuid
BenchmarkV4-4           	10000000	       112 ns/op
BenchmarkV4Parallel/X2/A-4       	20000000	       109 ns/op
BenchmarkV4Parallel/X2/B-4       	20000000	       109 ns/op
BenchmarkV4Parallel/X5/A-4       	20000000	       113 ns/op
BenchmarkV4Parallel/X5/B-4       	20000000	       109 ns/op
BenchmarkV4Parallel/X5/C-4       	20000000	       109 ns/op
BenchmarkV4Parallel/X5/D-4       	20000000	       110 ns/op
BenchmarkV4Parallel/X5/E-4       	20000000	       109 ns/op
PASS
ok  	github.com/as/uuid	22.638s
```

- Will this generate a bunch of zeroes like some other `uuid` packages?
> read the tests

```
go test -list .
TestV4
TestRace
TestProbabilityDistribution
```

- Does it read from a file descriptor?
> It reads from your system's entropy source once, at initialization time.

- Are there any conditions under which this will panic
> If the initial `crypto/rand` `Reader` can't read `<100` random bytes. This happens once at init time.

- License
>BSD


