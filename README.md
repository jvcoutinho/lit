# Lit

Lit is an expressive and fast HTTP framework for Golang. It aims to enhance development by
providing great simplicity, extensibility and maintainability.

## Documentation

Check [Lit package](https://pkg.go.dev/github.com/jvcoutinho/lit) documentation.

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

	server := &http.Server{Addr: ":8080", Handler: r}
	log.Fatalln(server.ListenAndServe())
}

func HelloWorld(r *lit.Request) lit.Response {
	return render.OK("Hello, World!")
}
```

Then you can start adding new routes and middlewares.

## Features

- **Great speed**: It uses [httprouter](https://github.com/julienschmidt/httprouter), a very fast zero-allocation
  router. See benchmarks.
- **Maintainability**: Its constructs make room for declarative programming, creating code that is more readable,
  maintainable, extensible and testable.
- **Flexibility**: It allows one to easily extend its constructs. Creating new validations, middlewares and responses,
  for example, is very simple and intuitive.

## Benchmarks

