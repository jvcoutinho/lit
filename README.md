# Lit 🔥

Lit is an expressive and fast HTTP framework for Golang. It aims to enhance development by
providing great simplicity, extensibility and maintainability.

## Documentation

Check [Lit](https://pkg.go.dev/github.com/jvcoutinho/lit#section-documentation) documentation.

## Getting started

Create a new Go project and import Lit with the command:

```
go get github.com/jvcoutinho/lit
```

Write this to your `main.go` file:

```go
package main

import (
	"log"
	"net/http"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
)

func main() {
	r := lit.NewRouter()
	r.Use(lit.Log)
	r.Use(lit.Recover)

	r.GET("/", HelloWorld)

	server := http.Server{Addr: ":8080", Handler: r}
	log.Fatalln(server.ListenAndServe())
}

func HelloWorld(r *lit.Request) lit.Response {
	return render.OK("Hello, World!")
}
```

Then you can start adding new routes and middlewares.

## Features

- **Speed and efficiency**: It uses [httprouter](https://github.com/julienschmidt/httprouter), a very fast
  zero-allocation
  router. See benchmarks.
- **Expressiveness**: Its constructs make room for declarative programming, creating code that is more readable,
  maintainable, extensible and testable.
- **Flexibility**: It allows one to easily extend its constructs. Creating new validations, middlewares and responses,
  for example, is very simple and intuitive.
- **Idiomatic**: It uses latest Go features, such as generics, in order to build a framework that is elegant to code.

## Benchmarks

The instances below were tested with the specifications:

```
goos: windows
goarch: amd64
pkg: github.com/julienschmidt/go-http-routing-benchmark
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
```

You can check the methodology or try yourself with
[go-http-routing-benchmark](https://github.com/jvcoutinho/go-http-routing-benchmark).

### GitHub API

| Router      | Memory for handler registration | Number repetitions | Latency per repetition | Heap memory per repetition | Allocations per repetition |
|-------------|---------------------------------|--------------------|------------------------|----------------------------|----------------------------|
| Chi         | 94888 B                         | 10000              | 105680 ns/op           | 61713 B/op                 | 406 allocs/op              |
| Gin         | 94888 B                         | 84152              | 14056 ns/op            | 0 B/op                     | 0 allocs/op                |
| **Lit**  🔥 | **42088 B**                     | **65989**          | **17919 ns/op**        | **0 B/op**                 | **0 allocs/op**            |
| Gorilla     | 1319632 B                       | 612                | 1886488 ns/op          | 199684 B/op                | 1588 allocs/op             |
| HttpRouter  | 37136 B                         | 114933             | 10511 ns/op            | 0 B/op                     | 0 allocs/op                |
| Martini     | 485032 B                        | 596                | 2070638 ns/op          | 231420 B/op                | 2731 allocs/op             |

### Parse API

| Router      | Memory for handler registration | Number repetitions | Latency per repetition | Heap memory per repetition | Allocations per repetition |
|-------------|---------------------------------|--------------------|------------------------|----------------------------|----------------------------|
| Chi         | 9656 B                          | 315666             | 11549 ns/op            | 7904 B/op                  | 52 allocs/op               |
| Gin         | 7864 B                          | 3118021            | 1144 ns/op             | 0 B/op                     | 0 allocs/op                |
| **Lit**  🔥 | **5776 B**                      | **2670331**        | **1342 ns/op**         | **0 B/op**                 | **0 allocs/op**            |
| Gorilla     | 105448 B                        | 75290              | 47507 ns/op            | 23632 B/op                 | 198 allocs/op              |
| HttpRouter  | 5072 B                          | 4421617            | 816 ns/op              | 0 B/op                     | 0 allocs/op                |
| Martini     | 45808 B                         | 45362              | 79979 ns/op            | 25696 B/op                 | 305 allocs/op              |

### Fake API (only static routes)

| Router      | Memory for handler registration | Number repetitions | Latency per repetition | Heap memory per repetition | Allocations per repetition |
|-------------|---------------------------------|--------------------|------------------------|----------------------------|----------------------------|
| Chi         | 83160 B                         | 51685              | 69608 ns/op            | 47728 B/op                 | 314 allocs/op              |
| Gin         | 34344 B                         | 385328             | 9302 ns/op             | 0 B/op                     | 0 allocs/op                |
| **Lit**  🔥 | **25560 B**                     | **450072**         | **7944 ns/op**         | **0 B/op**                 | **0 allocs/op**            |
| Gorilla     | 582536 B                        | 7695               | 500949 ns/op           | 113042 B/op                | 1099 allocs/op             |
| HttpRouter  | 21712 B                         | 638422             | 5716 ns/op             | 0 B/op                     | 0 allocs/op                |
| Martini     | 309880 B                        | 3915               | 916608 ns/op           | 129211 B/op                | 2031 allocs/op             |

As seen, Lit consumes less memory to register a batch of routes and has performance comparable to the top performers,
such as Gin Gonic and HttpRouter.

## Contributing

Feel free to open issues or pull requests!

---

Copyright (c) 2023-2024

João Victor de Sá Ferraz Coutinho <joao.coutinho9@gmail.com>