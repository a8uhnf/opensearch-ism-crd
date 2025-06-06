package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	apiv1 "github.com/a8uhnf/opensearch-ism-crd/api/v1"
	"github.com/opensearch-project/opensearch-go"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// OpensearchIndexPolicy represents an index policy in OpenSearch.

type openSearchClient struct {
	client *opensearch.Client
	url    string
}

func (c *openSearchClient) CreateIndexPolicy(ctx context.Context, policyName string, policy *apiv1.OpensearchIndexPolicy) error {
	logr := logf.FromContext(ctx)
	logr.Info("Creating index policy", "policyName", policyName)
	if policyName == "" {
		logr.Error(nil, "Policy name cannot be empty")
		return errors.NewBadRequest("policyName cannot be empty")
	}
	p := make(map[string]interface{})

	p = map[string]interface{}{
		"policy": policy,
	}

	body, err := json.Marshal(p)
	bBody := bytes.NewBuffer(body)
	fmt.Println("Creating index policy with body:", string(body))
	fmt.Println("Creating index policy with body:", fmt.Sprintf("%s/_plugins/_ism/policies/%s", c.url, policyName))
	// Create a new HTTP request to create the index policy
	// Note: The OpenSearch client does not directly support creating index policies,
	// so we need to use the HTTP API directly.
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/_plugins/_ism/policies/%s", c.url, policyName), bBody)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logr.Error(err, "Failed to create HTTP request for index policy")
		return errors.NewInternalError(err)
	}
	resp, err := c.client.Transport.Perform(req)
	if err != nil {
		logr.Error(err, "Failed to create index policy")

		return errors.NewInternalError(err)
	}
	defer req.Body.Close()
	if resp.StatusCode >= 300 {
		iout, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(iout))
		logr.Error(err, "Failed to create index policy")
		return errors.NewInternalError(fmt.Errorf("failed to create policy: %d", resp.StatusCode))
	}

	logr.Info("Index policy created successfully", "policyName", policyName)
	// If the policy creation is successful, we can return nil
	// to indicate that the operation was successful.
	return nil
}
func (c *openSearchClient) GetIndexPolicy(ctx context.Context, policyName string) (*apiv1.OpensearchIndexPolicy, error) {
	// Implementation for retrieving an index policy from OpenSearch
	logr := logf.FromContext(ctx)
	logr.Info("Retrieving index policy", "policyName", policyName)
	if policyName == "" {
		return nil, errors.NewBadRequest("policyName cannot be empty")
	}
	// Create a new HTTP request to get the index policy
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/_plugins/_ism/policies/%s", c.url, policyName), nil)
	if err != nil {
		logr.Error(err, "Failed to create HTTP request for index policy")
		return nil, errors.NewInternalError(err)
	}
	logr.Info("Performing HTTP request to retrieve index policy", "policyName", policyName, "url", req.URL.String())

	resp, err := c.client.Transport.Perform(req)

	if resp.StatusCode == 404 {
		logr.Info("Failed to retrieve index policy")

		schemaResource := schema.GroupResource{
			Group:    apiv1.GroupVersion.Group,
			Resource: apiv1.GroupVersion.WithResource("opensearchindexpolicies").Resource,
		}
		err = errors.NewNotFound(
			schemaResource, fmt.Sprintf("index policy %s not found", policyName))

		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		logr.Error(err, "not found")
		return nil, errors.NewInternalError(fmt.Errorf("failed to retrieve policy: %d", resp.StatusCode))
	}
	// Here we would typically parse the response body to extract the index policy details.
	// For simplicity, we return an empty OpensearchIndexPolicy.
	logr.Info("Index policy retrieved successfully", "policyName", policyName)
	// If the policy retrieval is successful, we can return an empty OpensearchIndexPolicy
	// to indicate that the operation was successful.

	policy := &apiv1.OpensearchIndexPolicy{}
	jsonDecoder := json.NewDecoder(resp.Body)
	logr.Info("Decoding index policy response")
	// Decode the response body into the OpensearchIndexPolicy struct
	// Assuming the response body is in JSON format and contains the policy details.
	if err := jsonDecoder.Decode(policy); err != nil {
		logr.Error(err, "Failed to decode index policy response")
		return nil, errors.NewInternalError(err)
	}
	// You would typically unmarshal the response body into an OpensearchIndexPolicy struct.
	// For now, we return an empty policy as a placeholder.
	// Assuming the response body contains the policy details, we would unmarshal it here.
	// For example:
	return policy, nil
}
func (c *openSearchClient) DeleteIndexPolicy(ctx context.Context, policyName string) error {
	// Implementation for deleting an index policy from OpenSearch
	logr := logf.FromContext(ctx)
	logr.Info("Deleting index policy", "policyName", policyName)
	if policyName == "" {
		return errors.NewBadRequest("policyName cannot be empty")
	}
	// Create a new HTTP request to delete the index policy
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/_plugins/_ism/policies/%s", c.url, policyName), nil)
	if err != nil {
		logr.Error(err, "Failed to create HTTP request for deleting index policy")
		return errors.NewInternalError(err)
	}
	resp, err := c.client.Transport.Perform(req)
	if err != nil {
		logr.Error(err, "Failed to delete index policy")
		return errors.NewInternalError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		logr.Error(err, "Failed to delete index policy", "statusCode", resp.StatusCode)
		// If the response status code indicates an error, we return an internal error.
		return errors.NewInternalError(fmt.Errorf("failed to delete policy: %s", err.Error()))
	}
	logr.Info("Index policy deleted successfully", "policyName", policyName)
	return nil
}

func (c *openSearchClient) GetClusterHealth(ctx context.Context) (string, error) {
	// Implementation for retrieving the cluster health from OpenSearch
	logr := logf.FromContext(ctx)
	logr.Info("Retrieving cluster health")
	// Create a new HTTP request to get the cluster health
	req, err := http.NewRequest("GET", "/_cluster/health", nil)
	if err != nil {
		logr.Error(err, "Failed to create HTTP request for cluster health")
		return "", errors.NewInternalError(err)
	}
	resp, err := c.client.Transport.Perform(req)
	if err != nil {
		logr.Error(err, "Failed to retrieve cluster health")
		return "", errors.NewInternalError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", errors.NewInternalError(fmt.Errorf("failed to retrieve cluster health: %s", err.Error()))
	}
	// Here we would typically parse the response body to extract the cluster health details.
	// For simplicity, we return an empty string.
	logr.Info("Cluster health retrieved successfully")
	return "Cluster is healthy", nil
}

// OpenSearchConfig holds the configuration for connecting to an OpenSearch cluster.

type OpenSearchConfig struct {
	// URL is the OpenSearch cluster URL.
	URL string `json:"url"`
	// Username is the username for OpenSearch authentication.
	Username string `json:"username"`
	// Password is the password for OpenSearch authentication.
	Password string `json:"password"`
	// TLSConfig contains TLS configuration for secure connections.
	TLSConfig *http.Transport `json:"tls_config,omitempty"`
}
