package licensed

import (
	"fmt"
	"os"
	"strconv"

	toml "github.com/pelletier/go-toml"
)

var ErrFatal = fmt.Errorf("fatal")
var ErrConfig = fmt.Errorf("bad config")

func New(logger interface{ Printf(string, ...interface{}) }, optionalConfigPath ...string) (*LicenseServer, error) {
	if len(optionalConfigPath) > 1 {
		return nil, fmt.Errorf("need one config")
	}
	var configpath = ""
	if len(optionalConfigPath) != 0 {
		configpath = optionalConfigPath[0]
	}
	l := &LicenseServer{
		log:        logger,
		configpath: configpath,
	}
	if configpath != "" {
		l.log.Printf("Loading config: %q", configpath)
		ttree, err := toml.LoadFile(configpath)
		if err != nil {
			return nil, err
		}
		err = ttree.Unmarshal(l)
		if err != nil {
			return nil, err
		}
	}
	if s := os.Getenv("SECRET"); s != "" {
		l.log.Printf("found environmental: secret")
		l.PrivateKey = s
	}

	if s := os.Getenv("ADDR"); s != "" {
		println("found environmental: addr", s)
		l.Addr = s
	}
	if s := os.Getenv("PORT"); s != "" {
		println("found environmental: port", s)
		var err error
		l.Port, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	if l.PrivateKey == "" {
		println("Missing PrivateKey")
		return nil, ErrConfig
	}
	l.log = logger
	return l, nil
}
