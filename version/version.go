package version

import (
	_ "embed"
)

var Name = "x-realy"

//go:embed version
var Version string
