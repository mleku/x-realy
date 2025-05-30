package main

import (
	"os"

	"x.realy.lol/bech32encoding"
	"x.realy.lol/chk"
	"x.realy.lol/config"
	"x.realy.lol/hex"
	"x.realy.lol/log"
	"x.realy.lol/p256k"
	"x.realy.lol/version"
)

func main() {
	cfg := config.New()
	if cfg.Superuser == "" {
		log.F.F("SUPERUSER is not set")
		os.Exit(1)
	}
	log.I.F("starting %s version %s", version.Name, version.Version)
	a := cfg.Superuser
	var err error
	var dst []byte
	if dst, err = bech32encoding.NpubToBytes([]byte(a)); chk.E(err) {
		if _, err = hex.DecBytes(dst, []byte(a)); chk.E(err) {
			log.F.F("SUPERUSER is invalid: %s", a)
			os.Exit(1)
		}
	}
	super := &p256k.Signer{}
	if err = super.InitPub(dst); chk.E(err) {
		return
	}

}
