package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	cnf        *Config
	configPath string
)

func initConfig() error {
	flag.StringVar(&configPath, "c", "config.yaml", "Configuration File")
	flag.Parse()

	c, err := NewCfg(configPath)
	if err != nil {
		return err
	}
	cnf = c

	return err
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cnf.HttpCfg().Dir)))

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", cnf.HttpCfg().Host, cnf.HttpCfg().Port),
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	fmt.Println(fmt.Sprintf("Your server are running on %s:%d", cnf.HttpCfg().Host, cnf.HttpCfg().Port))
	server.ListenAndServeTLS(cnf.HttpCfg().Certificate, cnf.HttpCfg().TLSKey)
}
