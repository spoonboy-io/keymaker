package main

import (
	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/reprise"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger

func init(){
	logger = &koan.Logger{}
}

func main() {
	// write a console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Keymaker",
		Description:  "TODO",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/keymaker",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})
}