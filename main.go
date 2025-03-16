package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	config                Config
)

func loadConfig() error {
	file, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}
	return nil
}

func replaceImage(image string, rules []Rule) string {
	for _, rule := range rules {
		if strings.HasPrefix(image, rule.Source) {
			return strings.Replace(image, rule.Source, rule.Target, 1)
		}
	}
	return image
}

func patchContainerImages(containers []corev1.Container) {
	for i := range containers {
		containers[i].Image = replaceImage(containers[i].Image, config.Rules)
	}
}

func admissionHandler(w http.ResponseWriter, r *http.Request) {
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
		patchContainerImages(pod.Spec.Containers)
		patchContainerImages(pod.Spec.InitContainers)

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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func main() {
	if err := loadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	http.HandleFunc("/mutate", admissionHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/readyz", readinessHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
