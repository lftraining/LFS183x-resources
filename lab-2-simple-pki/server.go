package main

import (
    "net/http"
    "log"
)

func ZeroTrustServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("Zero Trust is awesome!\n"))
}

func main() {
    http.HandleFunc("/zero", ZeroTrustServer)
    err := http.ListenAndServeTLS(":8443", "certs/control-plane.example.crt", "certs/control-plane.example.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}