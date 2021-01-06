package licensed

// as plugin

import (
	"flag"
	"log"
	"os"
	//"github.com/aerth/aels/lib/aels"
)

func main() {
	logger := log.New(os.Stderr, "[ÆLicenseServer] ", log.LstdFlags)
	flag.Parse()
	l, err := New(logger, "")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Fatal(l.ListenAndServe())
}
