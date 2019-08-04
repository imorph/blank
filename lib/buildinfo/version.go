package buildinfo

import (
	"fmt"
)

// Version must be set via -ldflags '-X'
var Version string

// PrintVersion prints version to stout
func PrintVersion() {
	fmt.Println("app version:", Version)
}
