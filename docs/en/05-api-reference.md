# API Reference

Complete API reference documentation for the Coco library.

## Core Functions

### New

Create a new documentation handler.

```go
func New(path string, opts ...Option) http.Handler
```

**Parameters**:
- `path` - Path to OpenAPI specification file (relative or absolute)
- `opts` - Optional configuration options

**Returns**:
- `http.Handler` - Standard HTTP handler interface implementation

**Examples**:
```go
// Load from file
handler := coco.New("./openapi.json")

// With configuration options
handler := coco.New("./openapi.json",
    coco.Title("My API"),
    coco.Theme("dark"),
)
```

---

## Configuration Options

All configuration options are `Option` type functions that can be passed to `New()`.

### Spec

Load OpenAPI specification from byte array.

```go
func Spec(data []byte) Option
```

**Parameters**:
- `data` - OpenAPI specification byte array (JSON format)

**Example**:
```go
spec := []byte(`{"openapi": "3.0.0", ...}`)
handler := coco.New("", coco.Spec(spec))
```

**Note**: When using this option, pass empty string as first parameter.

---

### SpecURL

Load OpenAPI specification from remote URL.

```go
func SpecURL(url string) Option
```

**Parameters**:
- `url` - Remote URL of OpenAPI specification

**Example**:
```go
handler := coco.New("", 
    coco.SpecURL("https://api.example.com/openapi.json"),
)
```

**Note**: When using this option, pass empty string as first parameter.

---

### I18n

Register a custom language from a local translation file.

```go
func I18n(path string) Option
```

**Parameters**:
- `path` - Path to the translation JSON file (relative or absolute)

**Example**:
```go
handler := coco.New("./openapi.json",
    coco.I18n("./i18n.fr.json"),
    coco.Lang("custom"),
)
```

**File format**:
```json
{
  "name": "Français",
  "messages": {
    "search": "Recherche",
    "loading": "Chargement..."
  }
}
```

**Note**: `name` is shown in the language switcher; `messages` maps translation keys to localized strings. Missing keys fall back to English.

---

### I18nData

Register a custom language from an in-memory byte array.

```go
func I18nData(data []byte) Option
```

**Parameters**:
- `data` - Translation JSON as a byte array

**Example**:
```go
//go:embed i18n.fr.json
var frI18n []byte

handler := coco.New("./openapi.json",
    coco.I18nData(frI18n),
    coco.Lang("custom"),
)
```

**Note**: The source is arbitrary — embedded files, an inline Go string (`[]byte("...")`), or JSON built at runtime from a database. See the file format under [I18n](#i18n).

---

### I18nURL

Register a custom language by fetching a translation file from a remote URL.

```go
func I18nURL(url string) Option
```

**Parameters**:
- `url` - Remote URL of the translation file

**Example**:
```go
handler := coco.New("./openapi.json",
    coco.I18nURL("https://cdn.example.com/i18n.fr.json"),
    coco.Lang("custom"),
)
```

**Note**: See the file format under [I18n](#i18n). Missing keys fall back to English.

---

### Title

Set the documentation page title.

```go
func Title(title string) Option
```

**Parameters**:
- `title` - Documentation title string

**Default**: `"Coco API Docs"`

**Example**:
```go
handler := coco.New("./openapi.json",
    coco.Title("My API Documentation"),
)
```

**Effect**:
- Displayed in browser tab title
- Displayed at top of documentation page

---

### Theme

Set the default theme.

```go
func Theme(theme string) Option
```

**Parameters**:
- `theme` - Theme name

**Available values**:
- `"light"` - Light theme
- `"dark"` - Dark theme
- `"auto"` - Follow system (default)

**Default**: `"auto"`

**Example**:
```go
handler := coco.New("./openapi.json",
    coco.Theme("dark"),
)
```

**Note**: Users can switch themes anytime in the top-right corner, and the choice is saved in browser local storage.

---

### Lang

Set the interface language.

```go
func Lang(lang string) Option
```

**Parameters**:
- `lang` - Language code

**Available values**:
- `"en"` - English
- `"zh"` - Chinese
- `"custom"` - Custom language (requires `I18n` / `I18nData` / `I18nURL`)

**Default**: `"en"`

**Example**:
```go
handler := coco.New("./openapi.json",
    coco.Lang("en"),
)
```

**Note**: Users can switch languages anytime in the top-right corner.

---

### EnableDebug

Enable or disable debug panel.

```go
func EnableDebug(debug bool) Option
```

**Parameters**:
- `debug` - `true` to enable, `false` to disable

**Default**: `true`

**Example**:
```go
// Enable debug panel (default)
handler := coco.New("./openapi.json",
    coco.EnableDebug(true),
)

// Disable debug panel
handler := coco.New("./openapi.json",
    coco.EnableDebug(false),
)
```

**Features**:
- Test APIs in documentation interface
- Fill parameters and request body
- View response results
- Support various HTTP methods

---

### EnableExport

Enable or disable export feature.

```go
func EnableExport(export bool) Option
```

**Parameters**:
- `export` - `true` to enable, `false` to disable

**Default**: `true`

**Example**:
```go
// Enable export feature (default)
handler := coco.New("./openapi.json",
    coco.EnableExport(true),
)

// Disable export feature
handler := coco.New("./openapi.json",
    coco.EnableExport(false),
)
```

**Features**:
- Allow users to download OpenAPI specification files (JSON format)

---

### EnableHistory

Enable or disable history feature.

```go
func EnableHistory(history bool) Option
```

**Parameters**:
- `history` - `true` to enable, `false` to disable

**Default**: `true`

**Example**:
```go
// Enable history (default)
handler := coco.New("./openapi.json",
    coco.EnableHistory(true),
)

// Disable history
handler := coco.New("./openapi.json",
    coco.EnableHistory(false),
)
```

**Features**:
- Save debug panel request history
- Stored in browser local storage
- Convenient for repeated testing

---

## Type Definitions

### Option

Configuration option function type.

```go
type Option func(*config.Config)
```

This is a function type used to modify configuration. All configuration functions return this type.

---

## Configuration Structures

While these structures are used internally, understanding them helps understand how configuration works.

### Config

```go
type Config struct {
    Spec
    UI
    Feature
}
```

### Spec

Specification configuration.

```go
type Spec struct {
    Path string  // File path
    Data []byte  // Byte array
    URL  string  // Remote URL
}
```

### I18n

Custom translation source (resolved: bytes, then URL, then file path).

```go
type I18n struct {
    I18nPath string // File path
    I18nData []byte // Byte array
    I18nURL  string // Remote URL
}
```

### UI

UI configuration.

```go
type UI struct {
    Title string // Document title
    Theme string // Theme: "light", "dark", "auto"
    Lang  string // Language: "en", "zh", "custom"
    I18n         // Custom translation source
}
```

### Feature

Feature toggles.

```go
type Feature struct {
    Debug   bool // Debug panel
    Export  bool // Export feature
    History bool // History
}
```

---

## Usage Examples

### Minimal Configuration

```go
handler := coco.New("./openapi.json")
http.Handle("/docs/", handler)
```

### Complete Configuration

```go
handler := coco.New("./openapi.json",
    coco.Title("My API Documentation"),
    coco.Theme("dark"),
    coco.Lang("en"),
    coco.EnableDebug(true),
    coco.EnableExport(true),
    coco.EnableHistory(true),
)
http.Handle("/docs/", handler)
```

### Using Byte Array

```go
import _ "embed"

//go:embed openapi.json
var spec []byte

handler := coco.New("",
    coco.Spec(spec),
    coco.Title("Embedded API"),
)
```

### Using Remote URL

```go
handler := coco.New("",
    coco.SpecURL("https://api.example.com/openapi.json"),
    coco.Title("Remote API"),
)
```

### Using a Custom Language

```go
import _ "embed"

//go:embed i18n.fr.json
var frI18n []byte

handler := coco.New("./openapi.json",
    coco.I18nData(frI18n),
    coco.Lang("custom"),
)
```

### Environment-Based Configuration

```go
import "os"

func createHandler() http.Handler {
    isProd := os.Getenv("ENV") == "production"
    
    opts := []coco.Option{
        coco.Title("My API"),
    }
    
    if !isProd {
        opts = append(opts,
            coco.EnableDebug(true),
            coco.Lang("en"),
        )
    } else {
        opts = append(opts,
            coco.EnableDebug(false),
            coco.EnableExport(false),
        )
    }
    
    return coco.New("./openapi.json", opts...)
}
```

---

## Error Handling

Coco returns error responses in the following situations:

1. **File not found** - Returns 404 Not Found
2. **Invalid file format** - Returns 500 Internal Server Error
3. **URL inaccessible** - Returns 500 Internal Server Error

Check log output during development to diagnose issues.

---

## Performance Considerations

1. **File Caching** - OpenAPI specification files are cached in memory
2. **Static Resources** - Frontend resources are embedded in binary
3. **Concurrency Safe** - All operations are concurrency-safe

---

## Compatibility

- **Go Version**: >= 1.21
- **OpenAPI Version**: 2.0, 3.0.x, 3.1.x
- **Browsers**: Modern browsers (Chrome, Firefox, Safari, Edge)

---

## Related Documentation

- [Quick Start](./01-getting-started.md)
- [Configuration Guide](./02-configuration.md)
- [Framework Integration](./03-framework-integration.md)
- [OpenAPI Generation](./04-openapi-generation.md)
