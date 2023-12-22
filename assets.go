package switcher

import "embed"

// Empty used when reflected package name.
type Empty struct{}

//go:embed static
var StaticDir embed.FS
