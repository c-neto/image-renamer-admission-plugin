package main

import (
	"fmt"
	"net/http"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/server"
)

func main() {
	if err := server.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	http.HandleFunc("/mutate", server.AdmissionHandler)
	http.HandleFunc("/healthz", server.HealthHandler)
	http.HandleFunc("/readyz", server.ReadinessHandler)
	fmt.Println("Starting server at port 8443")
	if err := http.ListenAndServeTLS(":8443", "/tls/tls.crt", "/tls/tls.key", nil); err != nil {
		fmt.Println(err)
	}
}
