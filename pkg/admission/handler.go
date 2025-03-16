package admission

import (
	"net/http"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/config"
	"github.com/c-neto/image-renamer-admission-plugin/pkg/server"
	corev1 "k8s.io/api/core/v1"
)

var cfg config.Config

// Set configuration
func SetConfig(c config.Config) {
	cfg = c
}

// Patch container images
func patchContainerImages(containers []corev1.Container) {
	server.PatchContainerImages(containers)
}

// Admission handler
func AdmissionHandler(w http.ResponseWriter, r *http.Request) {
	server.AdmissionHandler(w, r)
}

// Health handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	server.HealthHandler(w, r)
}

// Readiness handler
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	server.ReadinessHandler(w, r)
}
