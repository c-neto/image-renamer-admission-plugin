package admission

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/c-neto/image-renamer-admission-plugin/pkg/config"
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

func SetConfig(c config.Config) {
	cfg = c
}

func replaceImage(image string, rules []config.Rule) string {
	for _, rule := range rules {
		if strings.HasPrefix(image, rule.Source) {
			return strings.Replace(image, rule.Source, rule.Target, 1)
		}
	}
	return image
}

func patchContainerImages(containers []corev1.Container) {
	for i := range containers {
		containers[i].Image = replaceImage(containers[i].Image, cfg.Rules)
	}
}

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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
