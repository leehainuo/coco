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

	mux.Handle("/docs/", coco.New("",
		coco.SpecURL("https://petstore3.swagger.io/api/v3/openapi.json"),
		coco.Title("Coco Demo - Remote Spec (Petstore)"),
	))

	log.Println("Server starting on http://localhost:8000")
	log.Println("API docs available at http://localhost:8000/docs/")
	log.Println("Using remote spec from Petstore API")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
