package aws

import (
	`encoding/json`
	`fmt`
	`io/ioutil`
	`net/http`
	`os`
)

type ManagedPolicies struct{}

func (p ManagedPolicies) GetName() string {
	return "aws::managed_policies"
}

func (p ManagedPolicies) GetFileName() string {
	return "managed_policies.json"
}

func (p ManagedPolicies) Install() error {
	policies, err := p.getManagedPolicies()
	if err != nil {

	}
	return writeServicesToJSONFile(Folder, p.GetFileName(), policies)
}

type ManagedPolicy struct {
	Name string `json:"name"`
	Arn  string `json:"arn"`
}

// getManagedPolicies .
func (p ManagedPolicies) getManagedPolicies() ([]ManagedPolicy, error) {
	// URL of the JSON file
	url := "https://aws-managed-policies-list.s3.eu-central-1.amazonaws.com/policies.json"

	// Fetch the JSON file
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var policies []ManagedPolicy
	err = json.Unmarshal(body, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

// GetPolicies reads the services from a JSON file in a specified folder
func (p ManagedPolicies) GetPolicies() ([]ManagedPolicy, error) {
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
	var policies []ManagedPolicy
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&policies); err != nil {
		return nil, err
	}

	return policies, nil
}
