package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Zero Trust is awesome!" to the response body
	_, _ = io.WriteString(w, "Zero Trust is awesome!\n")
}

func main() {
	// Set up a /zero resource handler
	http.HandleFunc("/zero", helloHandler)

	// Create a CA certificate pool and add the client cert chain pem file to it
	caCert, err := os.ReadFile("cert-chain-2.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("certs/control-plane.example.crt", "certs/control-plane.example.key"))
}
