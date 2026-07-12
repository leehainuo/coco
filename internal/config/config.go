package config

type Config struct {
	Spec
	UI
	Feature
}

type Spec struct {
	Path string
	Data []byte
	URL  string
}

type UI struct {
	Title string
	Theme string
	Lang  string
	I18n
}

// I18n holds the source of a custom translation file. The resolution
// priority mirrors Spec: inline Data, then URL, then file Path.
type I18n struct {
	I18nPath string
	I18nData []byte
	I18nURL  string
}

type Feature struct {
	Debug   bool
	Export  bool
	History bool
}
