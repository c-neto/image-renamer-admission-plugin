package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/admission"
	"github.com/c-neto/image-renamer-admission-plugin/pkg/config"
)

func main() {
	configPath := "config.yaml"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	admission.SetConfig(cfg)

	http.HandleFunc("/mutate", admission.AdmissionHandler)
	http.HandleFunc("/healthz", admission.HealthHandler)
	http.HandleFunc("/readyz", admission.ReadinessHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
