package main

import (
	"log"
	"net/http"

	"example/basic/pkg"

	"github.com/leehainuo/coco"
)

func main() {
	mux := http.NewServeMux()

	pkg.RegisterHTTPAPI(mux)

	mux.Handle("/docs/", coco.New("openapi.json",
		coco.Title("Coco Demo - Local File Spec"),
	))

	log.Println("Server starting on http://localhost:8000")
	log.Println("API docs available at http://localhost:8000/docs/")
	log.Println("Using local OpenAPI spec file")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
