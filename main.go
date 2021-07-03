package main

import (
	"flag"
	"laplace/config"
	"laplace/core"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:443", "Listen address")
	tls := flag.Bool("tls", false, "Use TLS")
	setconfig := flag.Bool("setconfig", false, "Generates a config file")
	certFile := flag.String("certFile", "files/server.crt", "TLS cert file")
	keyFile := flag.String("keyFile", "files/server.key", "TLS key file")
	flag.Parse()

	// Action performed when the config file is called
	if *setconfig {
		config.SetDefaults()
		return
	}

	rand.Seed(time.Now().UnixNano())
	server := core.GetHttp()

	if *tls {
		log.Println("Listening on TLS:", *addr)
		if err := http.ListenAndServeTLS(*addr, *certFile, *keyFile, server); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("Listening:", *addr)
		if err := http.ListenAndServe(*addr, server); err != nil {
			log.Fatalln(err)
		}
	}
}
