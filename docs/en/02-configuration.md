# Configuration Guide

Coco provides rich configuration options to customize the appearance and behavior of your documentation.

## Configuration Method

All configurations are passed to `coco.New()` via `Option` functions:

```go
handler := coco.New("./openapi.json",
    coco.Title("My API"),
    coco.Theme("dark"),
    coco.Lang("en"),
    // ... more options
)
```

## UI Configuration

### Document Title

Set the title of the documentation page (displayed in browser tab and page header):

```go
coco.Title("My API Documentation")
```

**Default**: `"Coco API Docs"`

### Theme Settings

Set the default theme:

```go
coco.Theme("dark")  // Dark theme
coco.Theme("light") // Light theme
coco.Theme("auto")  // Follow system (default)
```

**Available values**:
- `"light"` - Light theme
- `"dark"` - Dark theme
- `"auto"` - Follow system settings

**Default**: `"auto"`

**Note**: Users can switch themes anytime in the top-right corner, and the choice is saved in browser local storage.

### Language Settings

Set the interface language:

```go
coco.Lang("en")     // English
coco.Lang("zh")     // Chinese
coco.Lang("custom") // Custom language (see below)
```

**Available values**:
- `"en"` - English
- `"zh"` - Chinese
- `"custom"` - A custom language loaded via `coco.I18n` / `coco.I18nData` / `coco.I18nURL`

**Default**: `"en"`

**Note**: Users can switch languages anytime in the top-right corner.

### Custom Language (i18n)

Besides the built-in English and Chinese, you can add a third language by
providing your own translation file. Register it with one of the following
options, then select it via `coco.Lang("custom")`:

```go
coco.I18n("./i18n.fr.json")                     // from a local file
coco.I18nData(frI18n)                           // from in-memory bytes
coco.I18nURL("https://cdn.example.com/i18n.fr.json") // from a remote URL
```

The translation file is a JSON object of the following shape:

```json
{
  "name": "Français",
  "messages": {
    "search": "Recherche",
    "loading": "Chargement...",
    "response": "Réponse"
  }
}
```

- `name` - The label shown in the language switcher.
- `messages` - A flat map of translation keys to localized strings.

**Notes**:
- Missing keys automatically fall back to English.
- When no custom language is configured, the switcher only shows English and Chinese.
- Set `coco.Lang("custom")` to make the custom language the initial one; otherwise users can pick it from the switcher.

## Specification Configuration

### Load from File

Most common method, load OpenAPI spec from local file:

```go
coco.New("./openapi.json")
coco.New("./swagger.json")
coco.New("./docs/api-spec.json")
```

Supports relative and absolute paths, JSON format only.

### Load from Byte Array

Suitable for dynamically generated or embedded specs:

```go
spec := []byte(`{
    "openapi": "3.0.0",
    "info": {
        "title": "My API",
        "version": "1.0.0"
    }
}`)

handler := coco.New("", coco.Spec(spec))
```

**Use cases**:
- Integration with code generation tools like Huma
- Load specs from database or config center
- Embed spec files using Go embed

**Example - Using embed**:
```go
import _ "embed"

//go:embed openapi.json
var spec []byte

func main() {
    handler := coco.New("", coco.Spec(spec))
    // ...
}
```

### Load from Remote URL

Load spec from remote server:

```go
handler := coco.New("", 
    coco.SpecURL("https://api.example.com/openapi.json"),
)
```

**Note**: 
- Ensure URL is accessible
- Consider network latency and availability
- Use local files or embedded method in production

## Translation Configuration

A custom language is loaded from a translation file. Like the specification,
it can come from three sources, resolved in this priority: in-memory bytes,
then remote URL, then local file.

### Load from File

Load the translation file from the local filesystem:

```go
handler := coco.New("./openapi.json",
    coco.I18n("./i18n.fr.json"),
    coco.Lang("custom"),
)
```

**Use cases**:
- Local development, where you edit the translation file without recompiling
- The translation file ships next to the binary

### Load from Byte Array

Pass translation bytes that already live in memory:

```go
handler := coco.New("./openapi.json",
    coco.I18nData(frI18n),
    coco.Lang("custom"),
)
```

`I18nData` accepts any `[]byte`, so the source is up to you:

**Example - Using embed** (single self-contained binary):
```go
import _ "embed"

//go:embed i18n.fr.json
var frI18n []byte

func main() {
    handler := coco.New("./openapi.json",
        coco.I18nData(frI18n),
        coco.Lang("custom"),
    )
    // ...
}
```

**Example - Using an inline Go string** (no external file at all):
```go
const frI18n = `{
    "name": "Français",
    "messages": {
        "search": "Recherche",
        "loading": "Chargement..."
    }
}`

handler := coco.New("./openapi.json",
    coco.I18nData([]byte(frI18n)),
    coco.Lang("custom"),
)
```

**Example - Built at runtime** (e.g. loaded from a database):
```go
payload, _ := json.Marshal(map[string]any{
    "name":     "Français",
    "messages": messagesFromDB, // map[string]string
})

handler := coco.New("./openapi.json",
    coco.I18nData(payload),
    coco.Lang("custom"),
)
```

**Use cases**:
- Embed the translation file with Go `embed` for zero external dependencies
- Define the JSON directly as a Go string constant
- Generate translations dynamically or load them from a database / config center

### Load from Remote URL

Fetch the translation file from a remote server:

```go
handler := coco.New("./openapi.json",
    coco.I18nURL("https://cdn.example.com/i18n.fr.json"),
    coco.Lang("custom"),
)
```

**Use cases**:
- Translations hosted on a CDN, object storage, or config center
- Multiple services sharing and updating one translation file centrally

## Feature Toggles

### Debug Panel

Enable or disable API debug panel:

```go
coco.EnableDebug(true)  // Enable (default)
coco.EnableDebug(false) // Disable
```

**Default**: `true`

**Features**:
- Test APIs directly in documentation
- Fill parameters and request body
- View response results
- Support various HTTP methods

**Recommendation**:
- Development: Enable
- Production: Decide based on needs (can enable for internal docs)

### Export Feature

Enable or disable OpenAPI spec export:

```go
coco.EnableExport(true)  // Enable (default)
coco.EnableExport(false) // Disable
```

**Default**: `true`

**Features**:
- Allow users to download OpenAPI spec files (JSON format)
- Easy integration with other tools

### History

Enable or disable request history:

```go
coco.EnableHistory(true)  // Enable (default)
coco.EnableHistory(false) // Disable
```

**Default**: `true`

**Features**:
- Save debug panel request history
- Stored in browser local storage
- Convenient for repeated testing

## Complete Configuration Examples

### Development Environment

```go
handler := coco.New("./openapi.json",
    coco.Title("Development - API Documentation"),
    coco.Theme("dark"),
    coco.Lang("en"),
    coco.EnableDebug(true),
    coco.EnableExport(true),
    coco.EnableHistory(true),
)
```

### Production Environment

```go
handler := coco.New("./openapi.json",
    coco.Title("Production - API Documentation"),
    coco.Theme("auto"),
    coco.Lang("en"),
    coco.EnableDebug(false),  // Disable debug
    coco.EnableExport(false), // Disable export
    coco.EnableHistory(false),
)
```

### Internal Documentation

```go
handler := coco.New("./openapi.json",
    coco.Title("Internal API Documentation"),
    coco.Theme("auto"),
    coco.Lang("en"),
    coco.EnableDebug(true),   // Internal can debug
    coco.EnableExport(true),  // Allow export
    coco.EnableHistory(true),
)
```

### Using Embedded Spec

```go
import _ "embed"

//go:embed openapi.json
var apiSpec []byte

func main() {
    handler := coco.New("",
        coco.Spec(apiSpec),
        coco.Title("Embedded API Documentation"),
        coco.Lang("en"),
    )
    // ...
}
```

## Advanced Usage

### Dynamic Configuration

Configure dynamically based on environment variables:

```go
import "os"

func createHandler() http.Handler {
    isDev := os.Getenv("ENV") == "development"
    
    return coco.New("./openapi.json",
        coco.Title(getTitle()),
        coco.EnableDebug(isDev),
        coco.EnableExport(isDev),
    )
}

func getTitle() string {
    if os.Getenv("ENV") == "production" {
        return "Production API"
    }
    return "Development API"
}
```

### Multiple Documentation Instances

Create different documentation for different API versions:

```go
// API v1 documentation
v1Handler := coco.New("./openapi-v1.json",
    coco.Title("API v1 Documentation"),
)
http.Handle("/docs/v1/", v1Handler)

// API v2 documentation
v2Handler := coco.New("./openapi-v2.json",
    coco.Title("API v2 Documentation"),
)
http.Handle("/docs/v2/", v2Handler)
```

## Best Practices

1. **Use Environment Variables** - Configure feature toggles dynamically based on environment
2. **Embed Spec Files** - Use `embed` in production to avoid file path issues
3. **Reasonable Defaults** - Enable all features in development, choose carefully in production
4. **Clear Titles** - Use descriptive titles for easy identification
5. **Consider User Habits** - Use appropriate language for your audience

## Related Documentation

- [Quick Start](./01-getting-started.md)
- [Framework Integration](./03-framework-integration.md)
- [API Reference](./05-api-reference.md)
