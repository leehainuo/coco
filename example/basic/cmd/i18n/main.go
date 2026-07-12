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

	// Load a user-provided i18n.json to add a custom language (French here)
	// alongside the built-in English and Chinese. Lang("custom") selects it
	// as the initial language; the language switcher lets users pick any of
	// the three at runtime.
	mux.Handle("/docs/", coco.New("openapi.json",
		coco.Title("Coco Demo - Custom i18n"),
		coco.I18n("i18n.json"),
		coco.Lang("custom"),
	))

	log.Println("Server starting on http://localhost:8000")
	log.Println("API docs available at http://localhost:8000/docs/")
	log.Println("Using a custom i18n.json (French) as an extra language")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
