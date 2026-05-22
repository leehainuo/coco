package assets

import "embed"

// DistFS holds the compiled frontend assets built from frontend/.
// The dist/ directory is populated by running: cd frontend && npm run build
//
//go:embed all:dist
var DistFS embed.FS
