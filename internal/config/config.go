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
}

type Feature struct {
	Debug   bool
	Export  bool
	History bool
}
