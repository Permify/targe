package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

var Folder = "requirements"

type Types struct{}

func (t Types) GetName() string {
	return "aws::types"
}

func (t Types) GetFileName() string {
	return "types.json"
}

func (t Types) Install() error {
	services := t.getServices()
	return writeServicesToJSONFile(Folder, t.GetFileName(), services)
}

type ListTypesResponse struct {
	TypeSummaries []types.TypeSummary `json:"TypeSummaries"`
	NextToken     *string             `json:"NextToken,omitempty"`
}

type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetServices retrieves all CloudFormation resource types and binds them to a slice of Service structs
func (t Types) getServices() []Service {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// Create a CloudFormation client
	client := cloudformation.NewFromConfig(cfg)

	var services []Service
	var nextToken *string

	for {
		// Create the input for the ListTypes API call
		input := &cloudformation.ListTypesInput{
			Type:             types.RegistryTypeResource,
			Visibility:       types.VisibilityPublic,
			ProvisioningType: types.ProvisioningTypeFullyMutable,
			MaxResults:       aws.Int32(100),
			NextToken:        nextToken,
		}

		// Call the ListTypes API
		resp, err := client.ListTypes(context.TODO(), input)
		if err != nil {
			log.Fatalf("failed to list types, %v", err)
		}

		// Append results to the services slice
		for _, t := range resp.TypeSummaries {
			service := Service{
				Name:        aws.ToString(t.TypeName),
				Description: aws.ToString(t.Description),
			}
			services = append(services, service)
		}

		// Check if there is a next page
		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return services
}

// WriteServicesToJSONFile writes the services slice to a JSON file
func writeServicesToJSONFile(folder, filename string, services []Service) error {
	// Ensure the folder exists
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	// Combine folder and filename to get full file path
	filePath := folder + "/" + filename

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write JSON data to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ") // Pretty print with indentation
	if err := encoder.Encode(services); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// GetServices reads the services from a JSON file in a specified folder
func (t Types) GetServices() ([]Service, error) {
	// Ensure the folder exists
	if _, err := os.Stat(Folder); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder does not exist: %s", Folder)
	}

	// Construct the full file path
	filePath := Folder + "/" + t.GetFileName()

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the services
	var services []Service
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&services); err != nil {
		return nil, err
	}

	return services, nil
}
