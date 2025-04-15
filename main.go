package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/aerth/aels/licensed"
)

func main() {
	var (
		defaultconfig = "aels.toml"
		logger        = log.New(os.Stderr, "[ÆLicenseServer] ", log.LstdFlags)
		configpath    = flag.String("conf", defaultconfig, "path to toml config file. use $PORT and $SECRET environment to skip config.")
		isDebug       = flag.Bool("debug", false, "debug")
	)

	flag.Parse()

	if *isDebug {
		logger.SetFlags(log.Lshortfile | log.LstdFlags)
	}

	l, err := licensed.New(logger, *configpath)
	if err != nil {
		logger.Fatal("startup", err)
	}

	go func() {
		<-time.After(time.Second)
		logger.Println("Serving")
	}()
	logger.Fatal(l.ListenAndServe())
}
