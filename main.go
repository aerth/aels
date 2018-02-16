package main

import (
	"flag"
	"log"

	"github.com/aerth/aels/lib/aels"
)

func main() {
	log.SetFlags(log.Lshortfile)
	configpath := flag.String("conf", "", "path to toml config file. use $PORT and $SECRET environment to skip config.")
	flag.Parse()
	l, err := aels.New(*configpath)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(l.ListenAndServe())
}
