package main

import (
	"github.com/urlshortener/shortener"
	"github.com/gorilla/mux"
	"net/http"
	"flag"
	"fmt"
	"os"
)

func main() {

	path := flag.String("c", "", "Path to the configuration file")
	flag.Parse()

	config, err := shortener.ReadConfig(*path)
	if err != nil {
		fmt.Printf("Could not read configuration file '%v'. %v\n", *path, err)
		os.Exit(1)
	}

	shortener.ConfigGl = *config

	err = shortener.CreateDBConnection(config)
	if err != nil {
		fmt.Printf("Could not create database connection. %v\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", shortener.IndexHandler)
	r.HandleFunc("/create", shortener.CreateHandler)
	r.HandleFunc("/url/{id}", shortener.IdHandler)

	port := ":" + config.ServerPort
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Printf("Could not start http server. %v\n", err)
	}
}
