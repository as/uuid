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

- What about the gofrs package?
> That package just makes you check for an error and still reads from /dev/random. This package actually solves the problem at the root.

- How exactly does it work?
> It initializes an AES block cipher with a random key and random initialization vector. It uses the output of CryptBlocks to build the uuid and then runs the operation again to generate the next random UUID. The output will never repeat, even if the process is restarted.

- I don't trust AES. What if the hard-core predicate of AES doesn't hold?
> You have bigger things to worry about. Also, most UUID implementations do not require a true random number generator. Especially if your database has a unique constraint, as such an error is much better than a panic in a production service.

- Has this been used in production services?
> This package has been in use in at least 10 production services recieving a moderate volume of traffic. It has never emitted an error on the database side.

- License
>BSD


