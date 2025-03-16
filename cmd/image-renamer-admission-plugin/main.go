package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/server"
)

func main() {
	// Parse command-line flags
	useHTTP := flag.Bool("http", false, "Run server in HTTP mode for local development")
	flag.Parse()

	// Load server configuration
	if err := server.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Register HTTP handlers
	http.HandleFunc("/mutate", server.AdmissionHandler)
	http.HandleFunc("/healthz", server.HealthHandler)
	http.HandleFunc("/readyz", server.ReadinessHandler)

	// Start the server in either HTTP or HTTPS mode based on the flag
	if *useHTTP {
		fmt.Println("Starting server at port 8080 in HTTP mode")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Starting server at port 8443 in HTTPS mode")
		if err := http.ListenAndServeTLS(":8443", "/tls/tls.crt", "/tls/tls.key", nil); err != nil {
			fmt.Println(err)
		}
	}
}
