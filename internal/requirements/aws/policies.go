package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type Policies struct{}

func (p Policies) GetName() string {
	return "aws::policies"
}

func (p Policies) GetFileName() string {
	return "policies.json"
}

func (p Policies) Install() error {
	policies, err := p.getPolicies()
	if err != nil {
		return err
	}
	return writeServicesToJSONFile(Folder, p.GetFileName(), policies)
}

type Policy struct {
	Name string `json:"name"`
	Arn  string `json:"arn"`
}

// getPolicies .
func (p Policies) getPolicies() ([]Policy, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := iam.NewFromConfig(cfg)

	// Initialize the request input
	input := &iam.ListPoliciesInput{
		Scope:    types.PolicyScopeTypeLocal,
		MaxItems: aws.Int32(500),
	}

	// List policies with the current input
	resp, err := client.ListPolicies(context.TODO(), input)
	if err != nil {
		log.Fatalf("unable to list policies, %v", err)
	}

	var policies []Policy
	for _, policy := range resp.Policies {
		p := Policy{
			Name: aws.ToString(policy.PolicyName),
			Arn:  aws.ToString(policy.Arn),
		}
		policies = append(policies, p)
	}

	return policies, nil
}

// GetPolicies reads the services from a JSON file in a specified folder
func (p Policies) GetPolicies() ([]Policy, error) {
	// Ensure the folder exists
	if _, err := os.Stat(Folder); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder does not exist: %s", Folder)
	}

	// Construct the full file path
	filePath := Folder + "/" + p.GetFileName()

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the policies
	var policies []Policy
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&policies); err != nil {
		return nil, err
	}

	return policies, nil
}
