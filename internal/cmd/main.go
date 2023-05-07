package main

import (
	"log"
	"net/http"

	"github.com/jvcoutinho/lit/render"

	"github.com/jvcoutinho/lit"
)

func main() {
	router := lit.NewRouter()
	router.Handle("/app", http.MethodGet, WriteHelloWorld)
	log.Fatalln(router.ListenAndServe(":8080"))
}

func WriteHelloWorld(_ *lit.Context) lit.Result {
	return render.Ok("Hello World!")
}
