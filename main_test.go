package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c-neto/image-renamer-admission-controller/pkg/config"
	"github.com/c-neto/image-renamer-admission-controller/pkg/server"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestPatchContainerImages tests the PatchContainerImages function
func TestPatchContainerImages(t *testing.T) {
	cfg := config.Config{
		Rules: []config.Rule{
			{Source: "nginx", Target: "my-registry/repo/nginx"},
			{Source: "my-registry/repo/busybox", Target: "my-registry/repo/busybox"},
		},
	}
	server.SetConfig(cfg)

	containers := []corev1.Container{
		{Image: "nginx"},
		{Image: "my-registry/repo/busybox"},
	}

	server.PatchContainerImages(containers)

	if containers[0].Image != "my-registry/repo/nginx" {
		t.Errorf("expected %s, got %s", "my-registry/repo/nginx", containers[0].Image)
	}
	if containers[1].Image != "my-registry/repo/busybox" {
		t.Errorf("expected %s, got %s", "my-registry/repo/busybox", containers[1].Image)
	}
}

// TestAdmissionHandler tests the AdmissionHandler function
func TestAdmissionHandler(t *testing.T) {
	cfg := config.Config{
		Rules: []config.Rule{
			{Source: "nginx", Target: "my-registry/repo/nginx"},
		},
	}
	server.SetConfig(cfg)

	pod := corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Image: "nginx"},
			},
		},
	}
	rawPod, _ := json.Marshal(pod)

	admissionReview := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			UID:    "1234",
			Object: runtime.RawExtension{Raw: rawPod},
		},
	}
	rawAdmissionReview, _ := json.Marshal(admissionReview)

	req, err := http.NewRequest("POST", "/mutate", bytes.NewReader(rawAdmissionReview))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.AdmissionHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response admissionv1.AdmissionReview
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.Response.UID != "1234" {
		t.Errorf("expected UID %s, got %s", "1234", response.Response.UID)
	}

	if !response.Response.Allowed {
		t.Errorf("expected response to be allowed")
	}
}

// TestHealthHandler tests the HealthHandler function
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HealthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestReadinessHandler tests the ReadinessHandler function
func TestReadinessHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/readyz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ReadinessHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "ready"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
