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
	"crypto/tls"
	"fmt"
	batchv1 "github.com/a8uhnf/opensearch-ism-crd/api/v1"
	"github.com/a8uhnf/opensearch-ism-crd/internal/pkg/opensearch"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// OSIndexPolicyReconciler reconciles a OSIndexPolicy object
type OSIndexPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch.a8uhnf.com,resources=osindexpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.a8uhnf.com,resources=osindexpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch.a8uhnf.com,resources=osindexpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OSIndexPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *OSIndexPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := logf.FromContext(ctx)

	logr.Info("Reconciling OSIndexPolicy", "name", req.Name, "namespace", req.Namespace)

	// Fetch the OSIndexPolicy instance
	osIndexPolicy := &batchv1.OSIndexPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OSIndexPolicy",
			APIVersion: batchv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Get(ctx, req.NamespacedName, osIndexPolicy); err != nil {
		logr.Error(err, "Failed to get OSIndexPolicy")
		// If the resource is not found, it may have been deleted after the request was queued.
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		// Resource not found, return and don't requeue
		return ctrl.Result{}, nil
	}

	opensearchClient, err := opensearch.NewOpenSearchClient(ctx, opensearch.OpenSearchConfig{
		URL:      osIndexPolicy.Spec.OpensearhConnection.URL,      // Use the URL from the request spec",
		Username: osIndexPolicy.Spec.OpensearhConnection.Username, // Use the username from the request spec
		Password: osIndexPolicy.Spec.OpensearhConnection.Password, // Use the password from the request spec
		TLSConfig: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Set to true for testing purposes, should be false in production
			},
		},
	})
	if err != nil {
		logr.Error(err, "Failed to create OpenSearch client")
		// If the OpenSearch client cannot be created, return an error to requeue the request.
		return ctrl.Result{
			RequeueAfter: 30 * 1000000000, // Requeue after 30 seconds
		}, err
	}

	policy, err := opensearchClient.GetIndexPolicy(ctx, osIndexPolicy.Spec.PolicyID)

	if errors.IsNotFound(err) {
		logr.Error(err, "Index policy not found in OpenSearch, creating new policy", "policyName", osIndexPolicy.Name)

		if err := opensearchClient.CreateIndexPolicy(ctx, osIndexPolicy.Spec.PolicyID, &osIndexPolicy.Spec.Policy); err != nil {
			logr.Error(err, "Failed to create index policy in OpenSearch", "policyName", osIndexPolicy.Name)
			// If the index policy cannot be created, return an error to requeue the request.
			return ctrl.Result{
				RequeueAfter: 30 * 1000000000, // Requeue after 30 seconds
			}, err
		}
		logr.Info("Index policy created successfully in OpenSearch", "policyName", osIndexPolicy.Name)

		return ctrl.Result{
			RequeueAfter: 30 * 1000000000, // Requeue after 30 seconds
		}, nil
	}

	if err != nil {
		logr.Error(err, "Failed to retrieve index policy from OpenSearch")
		return ctrl.Result{
			RequeueAfter: 30 * 1000000000, // Requeue after 30 seconds
		}, err
	}

	fmt.Println(policy)

	// Here you would add your logic to handle the OSIndexPolicy.
	logr.Info("Successfully reconciled OSIndexPolicy", "name", osIndexPolicy.Name, "namespace", osIndexPolicy.Namespace)
	// For example, you might want to check the policy's status and update it accordingly.
	// This is a placeholder for your reconciliation logic.
	// You can update the status of the OSIndexPolicy if needed.
	// osIndexPolicy.Status.SomeField = "some value"
	if err := r.Status().Update(ctx, osIndexPolicy); err != nil {
		logr.Error(err, "Failed to update OSIndexPolicy status")
		return ctrl.Result{}, err
	}
	// If you want to requeue the request after some time, you can return a Result with RequeueAfter set.
	// If you want to requeue the request immediately, you can return ctrl.Result{Requeue: true}.
	// If you want to stop requeuing, return ctrl.Result{}.
	logr.Info("OSIndexPolicy reconciled successfully", "name", osIndexPolicy.Name)

	// Returning a result to requeue the request after 30 seconds.
	logr.Info("Requeuing OSIndexPolicy reconciliation after 30 seconds", "name", osIndexPolicy.Name)

	return ctrl.Result{
		RequeueAfter: 30 * 1000000000, // Requeue after 30 seconds
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OSIndexPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.OSIndexPolicy{}).
		Named("osindexpolicy").
		Complete(r)
}
