package coco

import (
	"net/http"

	"github.com/leehainuo/coco/internal"
	"github.com/leehainuo/coco/internal/config"
)

type Option func(*config.Config)

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

func Spec(data []byte) Option {
	return func(c *config.Config) {
		c.Data = data
		c.Path = ""
	}
}

func SpecURL(url string) Option {
	return func(c *config.Config) {
		c.URL = url
		c.Path = ""
	}
}

func Title(title string) Option {
	return func(c *config.Config) {
		c.Title = title
	}
}

func Theme(theme string) Option {
	return func(c *config.Config) {
		c.Theme = theme
	}
}

func Lang(lang string) Option {
	return func(c *config.Config) {
		c.Lang = lang
	}
}

func EnableDebug(debug bool) Option {
	return func(c *config.Config) {
		c.Debug = debug
	}
}

func EnableExport(export bool) Option {
	return func(c *config.Config) {
		c.Export = export
	}
}

func EnableHistory(history bool) Option {
	return func(c *config.Config) {
		c.History = history
	}
}
