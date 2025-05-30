package main

import (
	"os"

	"x.realy.lol/config"
	"x.realy.lol/log"
	"x.realy.lol/version"
)

func main() {
	cfg := config.New()
	if cfg.Superuser == "" {
		log.F.F("SUPERUSER is not set")
		os.Exit(1)
	}
	log.I.F("starting %s version %s", version.Name, version.Version)
}
