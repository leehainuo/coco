package main

import (
	_ "embed"
	"log"
	"net/http"

	"example/basic/pkg"

	"github.com/leehainuo/coco"
)

//go:embed openapi.json
var spec []byte

func main() {
	mux := http.NewServeMux()

	pkg.RegisterHTTPAPI(mux)

	mux.Handle("/docs/", coco.New("",
		coco.Spec(spec),
		coco.Title("Coco Demo - Embedded Spec"),
	))

	log.Println("Server starting on http://localhost:8000")
	log.Println("API docs available at http://localhost:8000/docs/")
	log.Println("Using embedded OpenAPI spec")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
