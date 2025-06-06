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

package v1

import (
	"context"
	"fmt"

	"crypto/tls"
	"github.com/a8uhnf/opensearch-ism-crd/internal/pkg/opensearch"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	batchv1 "github.com/a8uhnf/opensearch-ism-crd/api/v1"
)

// nolint:unused
// log is for logging in this package.
var osindexpolicylog = logf.Log.WithName("osindexpolicy-resource")

// SetupOSIndexPolicyWebhookWithManager registers the webhook for OSIndexPolicy in the manager.
func SetupOSIndexPolicyWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&batchv1.OSIndexPolicy{}).
		WithValidator(&OSIndexPolicyCustomValidator{}).
		WithDefaulter(&OSIndexPolicyCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-batch-a8uhnf-com-v1-osindexpolicy,mutating=true,failurePolicy=fail,sideEffects=None,groups=batch.a8uhnf.com,resources=osindexpolicies,verbs=create;update,versions=v1,name=mosindexpolicy-v1.kb.io,admissionReviewVersions=v1

// OSIndexPolicyCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OSIndexPolicy when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OSIndexPolicyCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting

}

var _ webhook.CustomDefaulter = &OSIndexPolicyCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind OSIndexPolicy.
func (d *OSIndexPolicyCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	osindexpolicy, ok := obj.(*batchv1.OSIndexPolicy)

	if !ok {
		return fmt.Errorf("expected an OSIndexPolicy object but got %T", obj)
	}
	osindexpolicylog.Info("Defaulting for OSIndexPolicy", "name", osindexpolicy.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-batch-a8uhnf-com-v1-osindexpolicy,mutating=false,failurePolicy=fail,sideEffects=None,groups=batch.a8uhnf.com,resources=osindexpolicies,verbs=create;update,versions=v1,name=vosindexpolicy-v1.kb.io,admissionReviewVersions=v1

// OSIndexPolicyCustomValidator struct is responsible for validating the OSIndexPolicy resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OSIndexPolicyCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &OSIndexPolicyCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type OSIndexPolicy.
func (v *OSIndexPolicyCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	osindexpolicy, ok := obj.(*batchv1.OSIndexPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a OSIndexPolicy object but got %T", obj)
	}
	osindexpolicylog.Info("Validation for OSIndexPolicy upon creation", "name", osindexpolicy.GetName())

	if osindexpolicy.Spec.PolicyID == "" {
		return nil, fmt.Errorf("policy_id must be specified in the OSIndexPolicy spec")
	}
	if osindexpolicy.Spec.OpensearhConnection.URL == "" {
		return nil, fmt.Errorf("opensearch_connection.url must be specified in the OSIndexPolicy spec")
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type OSIndexPolicy.
func (v *OSIndexPolicyCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	osindexpolicy, ok := newObj.(*batchv1.OSIndexPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a OSIndexPolicy object for the newObj but got %T", newObj)
	}
	osindexpolicylog.Info("Validation for OSIndexPolicy upon update", "name", osindexpolicy.GetName())

	// TODO(user): fill in your validation logic upon object update.

	if osindexpolicy.Spec.PolicyID == "" {
		return nil, fmt.Errorf("policy_id must be specified in the OSIndexPolicy spec")
	}
	if osindexpolicy.Spec.OpensearhConnection.URL == "" {
		return nil, fmt.Errorf("opensearch_connection.url must be specified in the OSIndexPolicy spec")
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type OSIndexPolicy.
func (v *OSIndexPolicyCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	osindexpolicy, ok := obj.(*batchv1.OSIndexPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a OSIndexPolicy object but got %T", obj)
	}
	osindexpolicylog.Info("Validation for OSIndexPolicy upon deletion", "name", osindexpolicy.GetName())

	opensearchClient, err := opensearch.NewOpenSearchClient(ctx, opensearch.OpenSearchConfig{
		URL:      osindexpolicy.Spec.OpensearhConnection.URL,      // Use the URL from the request spec
		Username: osindexpolicy.Spec.OpensearhConnection.Username, // Use the username from the request spec
		Password: osindexpolicy.Spec.OpensearhConnection.Password, // Use the password from the request spec
		TLSConfig: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Set to true for testing purposes, should be false in production
			},
		},
	})

	opensearchClient.DeleteIndexPolicy(ctx, osindexpolicy.Spec.PolicyID)
	if err != nil {
		osindexpolicylog.Error(err, "Failed to delete index policy in OpenSearch", "policyName", osindexpolicy.Spec.PolicyID)
		// If the index policy cannot be deleted, return an error to requeue the request.
		return nil, fmt.Errorf("failed to delete index policy in OpenSearch: %w", err)
	}
	osindexpolicylog.Info("Index policy deleted successfully in OpenSearch", "policyName", osindexpolicy.Spec.PolicyID)

	return nil, nil
}
