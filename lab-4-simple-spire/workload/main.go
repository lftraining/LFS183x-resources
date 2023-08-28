package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	apiv2 "github.com/spiffe/go-spiffe/v2/workloadapi"
)

func main() {
	ctx := context.Background()

	x509Source, err := apiv2.NewX509Source(ctx)
	if err != nil {
		log.Fatalf("unable to create x509 source: %v\n", err)
	}
	defer x509Source.Close()

	tlsConf := tlsconfig.TLSServerConfig(x509Source)

	srv := &http.Server{
		Addr:      ":8443",
		Handler:   &handler{},
		TLSConfig: tlsConf,
	}

	if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		log.Fatalf("unable to launch server: %v\n", err)
	}
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "TADA!")
}
