package locale

import "embed"

//go:embed active.*.toml
var LocaleFS embed.FS
