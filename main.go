package main

import (
	"rpg/MDRegeln/cmd"
)

var (
	// holds clis version
	Version = "1.0.0"
)

func main() {
	cmd.Version = Version
	cmd.Execute()
}
