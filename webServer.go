package main

import (
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

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", cnf.HttpCfg().Host, cnf.HttpCfg().Port),
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	fmt.Println(fmt.Sprintf("Your server are running on %s:%d", cnf.HttpCfg().Host, cnf.HttpCfg().Port))
	server.ListenAndServe()
}
