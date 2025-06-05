package opensearch

import (
	"context"

	apiv1 "github.com/a8uhnf/opensearch-ism-crd/api/v1"
	"github.com/opensearch-project/opensearch-go"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type OpenSearch interface {
	// CreateIndexPolicy creates an index policy in OpenSearch.
	CreateIndexPolicy(ctx context.Context, policyName string, policy *apiv1.OpensearchIndexPolicy) error
	// GetIndexPolicy retrieves an index policy from OpenSearch.
	GetIndexPolicy(ctx context.Context, policyName string) (*apiv1.OpensearchIndexPolicy, error)
	// DeleteIndexPolicy deletes an index policy from OpenSearch.
	DeleteIndexPolicy(ctx context.Context, policyName string) error
	// // GetIndexPolicies retrieves all index policies from OpenSearch.
	// GetIndexPolicies(ctx context.Context) ([]OpensearchIndexPolicy, error)
	GetClusterHealth(ctx context.Context) (string, error)
}

func NewOpenSearchClient(ctx context.Context, config OpenSearchConfig) (OpenSearch, error) {

	logr := logf.FromContext(ctx)
	logr.Info("Creating OpenSearch client", "url", config.URL)
	// Implementation of OpenSearch client creation
	// This would typically involve setting up a connection to the OpenSearch cluster
	// using the provided configuration.
	oCli, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{config.URL},
		Username:  config.Username,
		Password:  config.Password,
	})
	if err != nil {
		return nil, err
	}
	return &openSearchClient{
		client: oCli,
	}, nil
}
