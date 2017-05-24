package main

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("Will1.com"), //your domain here
		Cache:      autocert.DirCache("certs"),          //folder for storing certificates
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt
}
