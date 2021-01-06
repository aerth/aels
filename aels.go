package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	toml "github.com/pelletier/go-toml"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type LicenseServer struct {
	Title        string
	Port         string
	DatabasePath string
	MaxTries     int
	Debug        bool
	AllowCopies  bool
	PrivateKey   string
	handler      http.Handler
}

type License struct {
	Key []byte
}

func (l License) String() string {
	return string(l.Key)
}

var ErrFatal = fmt.Errorf("fatal")
var ErrConfig = fmt.Errorf("bad config")

func New(configpath ...string) (*LicenseServer, error) {
	l := &LicenseServer{}
	if configpath != nil && configpath[0] != "" {
		println("Loading", configpath[0])
		ttree, err := toml.LoadFile(configpath[0])
		if err != nil {
			return nil, err
		}
		err = ttree.Unmarshal(l)
		return l, err
	}
	if s := os.Getenv("SECRET"); s != "" {
		println("found environmental: secret", s)
		l.PrivateKey = s
	}

	if s := os.Getenv("ADDR"); s != "" {
		println("found environmental: port", s)
		l.Port = s
	}
	if l.PrivateKey == "" || l.Port == "" {
		// println("Loading config.toml")
		// ttree, err := toml.LoadFile("config.toml")
		// if err != nil {
		return nil, ErrConfig
		// }
		// err = ttree.Unmarshal(l)
		// return l, err
	}

	return l, nil
}

func (l *LicenseServer) ListenAndServe() error {
	if l.Port == "" {
		println("missing port")
		return ErrConfig
	}
	if l.handler == nil {
		l.handler = http.HandlerFunc(l.handlerFunc)
	}
	return http.ListenAndServe(l.Port, l.handler)
}

func (l *LicenseServer) handlerFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("[visit]", r.UserAgent(), r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("x-forwarded-for"))
	switch r.Method {
	case http.MethodGet:
		if s := r.FormValue("cmd"); s == "gen" {
			for i := 0; i < 100; i++ {
				l.generateLicense()
			}
			return
		}
		if r.URL.Path == "/" {
			w.Write([]byte("Welcome to the AELS System\b"))
		}
		http.NotFound(w, r)
		return
	case http.MethodPost:
		key := r.PostFormValue("license")
		if !l.checkLicense([]byte(key)) {
			log.Println("key bad:", key)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		log.Println("key good:", key)
		w.WriteHeader(http.StatusOK)
	default:
		http.NotFound(w, r)
		return
	}
}

func (l *LicenseServer) generateLicense() License {
	b, err := bcrypt.GenerateFromPassword([]byte(l.PrivateKey), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return License{}
	}
	b = []byte(hex.EncodeToString(b))
	fmt.Printf("%s,", b)
	return License{Key: b}
}
func (l *LicenseServer) checkLicense(key []byte) bool {
	key, err := hex.DecodeString(string(key))
	if err != nil {
		log.Println(err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword(key, []byte(l.PrivateKey)); err != nil {
		log.Println(err)
		return false
	}
	return true
}