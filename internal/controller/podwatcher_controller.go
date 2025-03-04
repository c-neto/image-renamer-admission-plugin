/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	// corev1 "github.com/c-neto/image-renamer-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/c-neto/image-renamer-operator/api/v1"
)

// PodWatcherReconciler reconciles a PodWatcher object
type PodWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.example.com,resources=podwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.example.com,resources=podwatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.example.com,resources=podwatchers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodWatcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PodWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("AAAA")
	logger.Info("carlos neto debbuging")

	// Fetch the PodWatcher resource
	var podWatcher appv1.PodWatcher
	if err := r.Get(ctx, req.NamespacedName, &podWatcher); err != nil {
		logger.Error(err, "unable to fetch PodWatcher")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Watch for Pods matching the label selector
	labelSelector := podWatcher.Spec.LabelSelector
	if labelSelector == nil {
		logger.Info("No label selector specified. Skipping reconciliation.")
		return ctrl.Result{}, nil
	}

	podList := &corev1.PodList{}
	listOpts := &client.ListOptions{
		Namespace: req.Namespace,
	}
	if err := r.List(ctx, podList, listOpts); err != nil {
		logger.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	logger.Info("carlos neto debbuging")

	for _, pod := range podList.Items {
		fmt.Sprintln(pod.Name)
	}

	for _, pod := range podList.Items {
		// Check if the pod matches the label selector
		matches := true
		for key, value := range labelSelector {
			if pod.Labels[key] != value {
				matches = false
				break
			}
		}
		if matches {
			for _, status := range pod.Status.ContainerStatuses {
				if status.RestartCount >= 0 {
					message := fmt.Sprintf("Pod '%s' in namespace '%s' has restarted %d times!",
						pod.Name, pod.Namespace, status.RestartCount)
					fmt.Println(message)

					// Update PodWatcher status
					podWatcher.Status.LastPodRestartTime = time.Now().String()
					if err := r.Status().Update(ctx, &podWatcher); err != nil {
						logger.Error(err, "failed to update PodWatcher status")
					}
				}

			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.PodWatcher{}).
		Complete(r)
}
