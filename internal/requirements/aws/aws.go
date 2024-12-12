package aws

import (
	"encoding/json"
	"fmt"
	"os"
)

var Folder = "requirements"

// WriteServicesToJSONFile writes the services slice to a JSON file
func writeServicesToJSONFile(folder, filename string, services interface{}) error {
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
