// Package coco serves an interactive, self-contained API documentation UI
// for an OpenAPI/Swagger specification.
//
// Create an http.Handler with New and mount it on any router:
//
//	http.Handle("/docs/", http.StripPrefix("/docs", coco.New("openapi.json")))
//
// Behaviour is customised through functional Options, for example:
//
//	handler := coco.New("openapi.json",
//		coco.Title("My API"),
//		coco.Theme("dark"),
//		coco.Lang("zh"),
//		coco.I18n("i18n.fr.json"),
//	)
package coco

import (
	"net/http"

	"github.com/leehainuo/coco/internal"
	"github.com/leehainuo/coco/internal/config"
)

// Option configures the documentation handler. Options are applied in
// order, so later Options override earlier ones.
type Option func(*config.Config)

// New builds an http.Handler that serves the documentation UI along with
// the OpenAPI spec loaded from path (a local JSON file).
//
// The spec source can be overridden with Spec or SpecURL. Additional
// Options tweak the UI (Title, Theme, Lang, I18n) and toggle features
// (EnableDebug, EnableExport, EnableHistory). Sensible defaults are used
// when no Options are supplied.
func New(path string, opts ...Option) http.Handler {
	c := config.Config{
		Spec: config.Spec{
			Path: path,
		},
		UI: config.UI{
			Title: "Coco API Docs",
			Theme: "auto",
			Lang:  "en",
		},
		Feature: config.Feature{
			Debug:   true,
			Export:  true,
			History: true,
		},
	}
	for _, opt := range opts {
		opt(&c)
	}

	return internal.New(c)
}

// Spec serves the OpenAPI spec from the given in-memory bytes instead of a
// file path. It clears any previously configured file path.
func Spec(data []byte) Option {
	return func(c *config.Config) {
		c.Data = data
		c.Path = ""
	}
}

// SpecURL serves the OpenAPI spec by fetching it from the given URL instead
// of a file path. It clears any previously configured file path.
func SpecURL(url string) Option {
	return func(c *config.Config) {
		c.URL = url
		c.Path = ""
	}
}

// Title sets the document title shown in the browser tab and sidebar.
func Title(title string) Option {
	return func(c *config.Config) {
		c.Title = title
	}
}

// Theme sets the initial color theme. Accepted values are "light", "dark"
// and "auto" (follow the system preference).
func Theme(theme string) Option {
	return func(c *config.Config) {
		c.Theme = theme
	}
}

// Lang sets the initial UI language. The built-in languages are "en" and
// "zh"; use "custom" together with I18n to select a user-provided language.
func Lang(lang string) Option {
	return func(c *config.Config) {
		c.Lang = lang
	}
}

// I18n registers a custom translation file loaded from the given local
// path, adding a user-provided language alongside the built-in en/zh.
//
// The file is a JSON object of the form:
//
//	{"name": "Français", "messages": {"search": "Recherche", ...}}
//
// where name is the label shown in the language switcher and messages maps
// translation keys to localized strings. Missing keys fall back to English.
func I18n(path string) Option {
	return func(c *config.Config) {
		c.I18nPath = path
	}
}

// I18nData registers a custom translation file from the given in-memory
// bytes. See I18n for the expected JSON shape.
func I18nData(data []byte) Option {
	return func(c *config.Config) {
		c.I18nData = data
	}
}

// I18nURL registers a custom translation file fetched from the given URL.
// See I18n for the expected JSON shape.
func I18nURL(url string) Option {
	return func(c *config.Config) {
		c.I18nURL = url
	}
}

// EnableDebug toggles the built-in "try it out" request debugger. It is
// enabled by default.
func EnableDebug(debug bool) Option {
	return func(c *config.Config) {
		c.Debug = debug
	}
}

// EnableExport toggles the button that exports the OpenAPI spec. It is
// enabled by default.
func EnableExport(export bool) Option {
	return func(c *config.Config) {
		c.Export = export
	}
}

// EnableHistory toggles persistence of the debugger's request history. It
// is enabled by default.
func EnableHistory(history bool) Option {
	return func(c *config.Config) {
		c.History = history
	}
}
