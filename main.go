package main

import (
	"laplace/core"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) <= 2 || os.Args[1] == "" || os.Args[2] == "" {
		log.Fatalln("Usage: laplace [ip] [port]")
	}
	ip, port := os.Args[1], os.Args[2]
	server := core.GetHttp(ip + ":" + port)
	if err := http.ListenAndServeTLS("0.0.0.0:" + port, "files/server.crt", "files/server.key", server); err != nil {
		log.Fatalln(err)
	}
}
