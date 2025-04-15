package licensed

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type LicenseServer struct {
	Title        string
	Port         int
	Addr         string
	DatabasePath string
	MaxTries     int
	Debug        bool
	AllowCopies  bool
	PrivateKey   string
	Handler      http.Handler
	log          interface {
		Printf(string, ...interface{})
	}
	configpath string
}

type License struct {
	Key []byte
}

func (l License) String() string {
	return string(l.Key)
}

func (l *LicenseServer) ListenAndServe() error {
	if l.log == nil {
		return fmt.Errorf("nil log")
	}
	if l.Port != 0 && l.Addr != "" {
		println("cant have both port and addr, pick one")
		return ErrConfig
	}
	if l.Port == 0 && l.Addr == "" {
		println("missing port and addr, pick one.")
		return ErrConfig
	}
	if l.Addr == "" {
		l.Addr = fmt.Sprintf("127.0.0.1:%d", l.Port)
	}
	if l.Handler == nil {
		l.Handler = l
	}
	return http.ListenAndServe(l.Addr, l.Handler)
}

func (l *LicenseServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l.log != nil {
		l.log.Printf("[visit] %s %s %s %s %s", r.UserAgent(), r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("x-forwarded-for"))
	}
	switch r.Method {
	default:
		http.NotFound(w, r)
		return
	case http.MethodGet:
		if s := r.FormValue("cmd"); s == "gen" {
			if os.Getenv("GEN") != "1" {
				fmt.Fprintf(w, "need GEN=1\n")
				return
			}
			for i := 0; i < 100; i++ {
				l.generateLicense()
			}
			return
		}
		if r.URL.Path == "/" {
			w.Write([]byte("Welcome to the AELS System\b\n"))
			return
		}
		http.NotFound(w, r)
		return
	case http.MethodPost:
		if ct := r.Header.Get("Content-Type"); ct == "application/json" {
			t := struct {
				License string `json:"license"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				r.Body.Close()
				return
			}
			if !l.checkLicense([]byte(t.License)) {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			l.log.Printf("Access granted: %q", t.License[:6])
			w.WriteHeader(http.StatusOK)
			return
		}
		key := r.PostFormValue("license")
		if len(key) == 0 {
			l.log.Printf("key bad: license=%q", key)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if !l.checkLicense([]byte(key)) {
			l.log.Printf("key bad: %q", key)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		l.log.Printf("key good: %q", key)
		w.WriteHeader(http.StatusOK)
	}
}

func (l *LicenseServer) GenerateLicense() License {
	return l.generateLicense()
}
func (l *LicenseServer) generateLicense() License {
	b, err := bcrypt.GenerateFromPassword([]byte(l.PrivateKey), bcrypt.DefaultCost)
	if err != nil {
		l.log.Printf("error generating: %v", err)
		return License{}
	}
	b = []byte(hex.EncodeToString(b))
	fmt.Printf("%s,", b)
	return License{Key: b}
}
func (l *LicenseServer) checkLicense(key []byte) bool {
	key, err := hex.DecodeString(string(key))
	if err != nil {
		l.log.Printf("error decoding license key: %v", err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword(key, []byte(l.PrivateKey)); err != nil {
		l.log.Printf("error checking license validity: %v", err)
		return false
	}
	return true
}
