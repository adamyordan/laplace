package main

import (
	"flag"
	"laplace/core"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:443", "Listen address")
	noTls := flag.Bool("no-tls", false, "Don't use TLS")
	certFile := flag.String("certFile", "files/server.crt", "TLS cert file")
	keyFile := flag.String("keyFile", "files/server.key", "TLS key file")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	server := core.GetHttp()

	if !*noTls {
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
