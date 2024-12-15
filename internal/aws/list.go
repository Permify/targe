package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/aws/aws-sdk-go-v2/service/iam"

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

func ListGroupsForUser(ctx context.Context, cfg aws.Config, username string) ([]string, error) {
	var groups []string
	client := iam.NewFromConfig(cfg)

	input := &iam.ListGroupsForUserInput{
		UserName: aws.String(username),
	}

	resp, err := client.ListGroupsForUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to list groups for user %s: %v", username, err)
	}

	// Extract group names
	for _, group := range resp.Groups {
		groups = append(groups, *group.GroupName)
	}

	return groups, nil
}

func ListUsers(ctx context.Context, cfg aws.Config) (*iam.ListUsersOutput, error) {
	client := iam.NewFromConfig(cfg)
	input := &iam.ListUsersInput{}
	return client.ListUsers(ctx, input)
}

func ListPolicies(ctx context.Context, cfg aws.Config) (*iam.ListPoliciesOutput, error) {
	client := iam.NewFromConfig(cfg)

	// Initialize the request input
	input := &iam.ListPoliciesInput{
		Scope:    types.PolicyScopeTypeLocal,
		MaxItems: aws.Int32(500),
	}

	// List policies with the current input
	resp, err := client.ListPolicies(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ListGroups(ctx context.Context, cfg aws.Config) (*iam.ListGroupsOutput, error) {
	client := iam.NewFromConfig(cfg)
	input := &iam.ListGroupsInput{}
	return client.ListGroups(ctx, input)
}

func ListRoles(ctx context.Context, cfg aws.Config) (*iam.ListRolesOutput, error) {
	client := iam.NewFromConfig(cfg)
	input := &iam.ListRolesInput{}
	return client.ListRoles(ctx, input)
}

func ListAttachedUserPolicies(ctx context.Context, cfg aws.Config, username string) ([]string, error) {
	client := iam.NewFromConfig(cfg)
	var names []string

	// List the user's attached policies
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(username),
	}

	// Call AWS API to get the list of inline policies for the user
	resp, err := client.ListAttachedUserPolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to list inline policies for user %s: %v", username, err)
	}

	for _, p := range resp.AttachedPolicies {
		names = append(names, *p.PolicyName)
	}

	return names, nil
}

func ListUserInlinePolicies(ctx context.Context, cfg aws.Config, username string) ([]string, error) {
	client := iam.NewFromConfig(cfg)

	// List the user's attached policies
	input := &iam.ListUserPoliciesInput{
		UserName: aws.String(username),
	}

	// Call AWS API to get the list of inline policies for the user
	resp, err := client.ListUserPolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to list inline policies for user %s: %v", username, err)
	}

	return resp.PolicyNames, nil
}

func ListGroupInlinePolicies(ctx context.Context, cfg aws.Config, groupname string) ([]string, error) {
	client := iam.NewFromConfig(cfg)

	input := &iam.ListGroupPoliciesInput{
		GroupName: aws.String(groupname),
	}

	resp, err := client.ListGroupPolicies(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.PolicyNames, nil
}

func ListAttachedGroupPolicies(ctx context.Context, cfg aws.Config, groupname string) ([]string, error) {
	client := iam.NewFromConfig(cfg)

	var names []string

	// List the user's attached policies
	input := &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupname),
	}

	resp, err := client.ListAttachedGroupPolicies(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, p := range resp.AttachedPolicies {
		names = append(names, *p.PolicyName)
	}

	return names, nil
}

func ListRoleInlinePolicies(ctx context.Context, cfg aws.Config, rolename string) ([]string, error) {
	client := iam.NewFromConfig(cfg)

	input := &iam.ListRolePoliciesInput{
		RoleName: aws.String(rolename),
	}

	resp, err := client.ListRolePolicies(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return resp.PolicyNames, nil
}

func ListAttachedRolePolicies(ctx context.Context, cfg aws.Config, rolename string) ([]string, error) {
	client := iam.NewFromConfig(cfg)

	var names []string

	// List the user's attached policies
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(rolename),
	}

	resp, err := client.ListAttachedRolePolicies(context.TODO(), input)
	if err != nil {
		log.Printf("Error fetching policies for user %s: %v", rolename, err)
		return nil, nil
	}

	for _, p := range resp.AttachedPolicies {
		names = append(names, *p.PolicyName)
	}

	return names, nil
}
