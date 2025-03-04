package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodWatcherSpec defines the desired state of PodWatcher
type PodWatcherSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LabelSelector map[string]string `json:"labelSelector,omitempty"`
}

// PodWatcherStatus defines the observed state of PodWatcher
type PodWatcherStatus struct {
	LastPodRestartTime string `json:"lastPodRestartTime,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodWatcher is the Schema for the podwatchers API
type PodWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodWatcherSpec   `json:"spec,omitempty"`
	Status PodWatcherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodWatcherList contains a list of PodWatcher
type PodWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodWatcher{}, &PodWatcherList{})
}
