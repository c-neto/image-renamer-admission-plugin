package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/config"
	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	cfg                   config.Config
)

// Load configuration from file or environment variable
func LoadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	return nil
}

// Replace image based on rules
func replaceImage(image string, rules []config.Rule) string {
	for _, rule := range rules {
		if strings.HasPrefix(image, rule.Source) {
			return strings.Replace(image, rule.Source, rule.Target, 1)
		}
	}
	return image
}

// Patch container images
func PatchContainerImages(containers []corev1.Container) {
	for i := range containers {
		containers[i].Image = replaceImage(containers[i].Image, cfg.Rules)
	}
}

// Admission handler
func AdmissionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	var admissionReview admissionv1.AdmissionReview
	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReview); err != nil {
		http.Error(w, "could not decode admission review", http.StatusBadRequest)
		return
	}

	admissionResponse := admissionv1.AdmissionResponse{
		UID: admissionReview.Request.UID,
	}

	var pod corev1.Pod
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		admissionResponse.Result = &metav1.Status{
			Message: err.Error(),
		}
	} else {
		PatchContainerImages(pod.Spec.Containers)
		PatchContainerImages(pod.Spec.InitContainers)

		patchBytes, err := json.Marshal(pod)
		if err != nil {
			admissionResponse.Result = &metav1.Status{
				Message: err.Error(),
			}
		} else {
			admissionResponse.Patch = patchBytes
			patchType := admissionv1.PatchTypeJSONPatch
			admissionResponse.PatchType = &patchType
			admissionResponse.Allowed = true
		}
	}

	admissionReview.Response = &admissionResponse
	respBytes, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// Health handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Readiness handler
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
