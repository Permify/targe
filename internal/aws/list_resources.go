package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Resource represents a single AWS resource.
type Resource struct {
	Name string
	Arn  string
}

// ListResources fetches resources of a given type from the AWS Resource Explorer API.
func ListResources(resourceType string) ([]Resource, error) {
	const service = "resource-explorer"
	const endpoint = "https://resource-explorer.%s.amazonaws.com/resources-list"

	// Validate inputs
	if resourceType == "" {
		return nil, fmt.Errorf("resourceType cannot be empty")
	}

	// Load AWS configuration (including credentials from ~/.aws/credentials)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Format the endpoint with the region
	url := fmt.Sprintf(endpoint, cfg.Region)

	// Retrieve actual credentials
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AWS credentials: %w", err)
	}

	// Prepare the request payload
	payload := fmt.Sprintf(`{"ResourceType":"%s"}`, resourceType)

	// Hash the payload
	hash := sha256.Sum256([]byte(payload))
	payloadHash := hex.EncodeToString(hash[:])

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add required headers
	addHeaders(req, payloadHash)

	// Sign the request using AWS Signature Version 4
	err = signRequest(req, payloadHash, creds, service, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %s", resp.Status)
	}

	// Parse the response body
	var listResourcesResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&listResourcesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Extract resources
	resources, err := extractResources(listResourcesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to extract resources: %w", err)
	}

	return resources, nil
}

// addHeaders adds required headers to the HTTP request.
func addHeaders(req *http.Request, payloadHash string) {
	req.Header.Add("accept", "*/*")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-amz-content-sha256", payloadHash)
	req.Header.Add("x-amz-date", time.Now().UTC().Format("20060102T150405Z"))
	req.Header.Add("user-agent", "Custom-Client/1.0")
}

// signRequest signs the HTTP request using AWS Signature Version 4 and the loaded configuration.
func signRequest(req *http.Request, payloadHash string, creds aws.Credentials, service, region string) error {
	signer := v4.NewSigner()
	return signer.SignHTTP(context.TODO(), creds, req, payloadHash, service, region, time.Now())
}

// extractResources extracts resources from the API response.
func extractResources(response map[string]interface{}) ([]Resource, error) {
	resourceArns, ok := response["ResourceArns"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: 'ResourceArns' missing or invalid")
	}

	var resources []Resource
	for _, arnInterface := range resourceArns {
		arn, ok := arnInterface.(string)
		if !ok {
			return nil, fmt.Errorf("invalid ARN format in response")
		}

		resources = append(resources, Resource{
			Name: extractResourceName(arn),
			Arn:  arn,
		})
	}
	return resources, nil
}

// extractResourceName extracts the last component of an ARN.
func extractResourceName(arn string) string {
	re := regexp.MustCompile(`([^/:]+)$`)
	matches := re.FindStringSubmatch(arn)
	if len(matches) > 1 {
		return matches[1]
	}
	return arn // Fallback to the full ARN if no match
}
