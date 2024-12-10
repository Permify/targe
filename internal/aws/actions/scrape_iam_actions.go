package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://docs.aws.amazon.com/service-authorization/latest/reference/"
)

// Logger setup
var logger = log.New(os.Stdout, "IAM-CONNECTOR: ", log.LstdFlags)

// Fetches and parses a document from a given URL
func fetchDocument(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}

// Scrapes the main reference page to get all service links
func getServiceLinks() ([]string, error) {
	doc, err := fetchDocument(baseURL + "reference_policies_actions-resources-contextkeys.html")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch main page: %w", err)
	}

	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists && strings.HasPrefix(href, "./list_") {
			links = append(links, href[2:]) // Remove './' prefix
		}
	})

	return links, nil
}

// Parses service action details from a given service URL
func parseServiceActions(link string) (map[string]map[string]string, string, error) {
	doc, err := fetchDocument(baseURL + link)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch service page %s: %w", link, err)
	}

	actions := make(map[string]map[string]string)
	servicePrefix := doc.Find("code").First().Text()

	doc.Find("div.table-contents").Each(func(i int, table *goquery.Selection) {
		headers := table.Find("th").Map(func(i int, s *goquery.Selection) string {
			return strings.ToLower(s.Text())
		})

		if !contains(headers, "actions") || !contains(headers, "description") {
			return
		}

		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			cells := row.Find("td")
			// Ensure at least 3 cells exist in the row before proceeding
			if cells.Length() < 3 {
				return
			}

			// Extract the action, description, and access level safely
			actionText := cells.Eq(0).Text()
			if actionText == "" {
				return // Skip rows where the action is empty
			}

			actionFields := strings.Fields(actionText)
			if len(actionFields) == 0 {
				return // Skip rows where no valid action name is found
			}
			action := actionFields[0]

			description := cells.Eq(1).Text()
			accessLevel := cells.Eq(2).Text()

			// Initialize maps if necessary
			if _, exists := actions[action]; !exists {
				actions[action] = make(map[string]string)
			}

			// Assign values to the action map
			actions[action]["description"] = description
			actions[action]["access"] = accessLevel
		})
	})

	return actions, servicePrefix, nil
}

// Combines data for all services
func collectAllActions(links []string) map[string]map[string]map[string]string {
	data := make(map[string]map[string]map[string]string)

	for _, link := range links {
		logger.Printf("Processing link: %s\n", link)
		actions, servicePrefix, err := parseServiceActions(link)
		if err != nil {
			logger.Printf("Error processing %s: %v", link, err)
			continue
		}

		if _, exists := data[servicePrefix]; !exists {
			data[servicePrefix] = make(map[string]map[string]string)
		}

		for action, details := range actions {
			data[servicePrefix][action] = details
		}
	}

	return data
}

// Saves data to a JSON file
func saveToJSONFile(data map[string]map[string]map[string]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// Checks if a slice contains a specific item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func SaveActions() {
	logger.Println("Starting IAM actions scraping...")

	// Step 1: Get service links
	links, err := getServiceLinks()
	if err != nil {
		logger.Fatalf("Error fetching service links: %v", err)
	}

	// Step 2: Collect all actions
	data := collectAllActions(links)
	// Step 3: Save results to a JSON file
	outputFile := "actions.json"
	if err := saveToJSONFile(data, outputFile); err != nil {
		logger.Fatalf("Error saving JSON file: %v", err)
	}

	logger.Printf("Scraping completed! Data saved to %s\n", outputFile)
}

type ActionConfig map[string]map[string]ActionDetail

type ActionDetail struct {
	Access      string `json:"access"`
	Description string `json:"description"`
}

func LoadActions(filePath string) (ActionConfig, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open actions.json: %v", err)
	}
	defer file.Close()

	// Decode the JSON data into an ActionConfig (map of service -> actions)
	var actions ActionConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&actions); err != nil {
		return nil, fmt.Errorf("unable to decode actions.json: %v", err)
	}

	return actions, nil
}
