package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GPTResponse struct {
	Action                string            `json:"action"`
	Principal             map[string]string `json:"principal"`
	RequestedResourceType string            `json:"requested_resource_type"`
	RequestedResource     string            `json:"requested_resource"`
	IsManagedPolicy       bool              `json:"isManagedPolicy"`
	Policy                string            `json:"policy"`
	Error                 bool              `json:"error"`
	Confidence            int               `json:"confidence"`
}

func generatePrompt(userInput string) string {
	return userInput
}

func callGPTWithSchema(apiKey, model, prompt string, temperature float64) (GPTResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"
	schema := map[string]interface{}{
		"name": "iam_request",
		"schema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        []string{"string", "null"},
					"description": "The action type. If undetermined, return null.",
					"enum": []string{
						"attach_policy_to_user",
						"attach_policy_to_group",
						"attach_policy_to_role",
						"create_policy",
						"add_user_to_group",
						"need_custom_policy",
					},
				},
				"error": map[string]interface{}{
					"type":        "boolean",
					"description": "Indicates if the command cannot be managed by the specified actions.",
				},
				"principal": map[string]interface{}{
					"type":        "object",
					"description": "The target entity to which the policy will be attached.",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": []string{"string", "null"},
							"enum": []string{"user", "group", "role"},
						},
						"name": map[string]interface{}{
							"type":        []string{"string", "null"},
							"description": "The name of the target entity.",
						},
					},
					"required":             []string{"type", "name"},
					"additionalProperties": false,
				},
				"requested_resource_type": map[string]interface{}{
					"type":        []string{"string", "null"},
					"description": "The type of aws resource user wants access for. S3, RDS, EC2 etc",
				},
				"requested_resource": map[string]interface{}{
					"type":        []string{"string", "null"},
					"description": "The name of the resource user wants access for.",
				},
				"isManagedPolicy": map[string]interface{}{
					"type":        "boolean",
					"description": "Indicates if the provided policy is an exact AWS managed policy.",
				},
				"policy": map[string]interface{}{
					"type":        []string{"string", "null"},
					"description": "The name of the policy. If it's too vague, return null. If the user input does not provide any meaningful context, the model must not guess a policy",
				},
				"confidence": map[string]interface{}{
					"type":        "integer",
					"description": "Confidence level from 1 to 10 about the policy name.",
				},
			},
			"required": []string{
				"error",
			},
			"additionalProperties": false,
		},
	}
	payload := map[string]interface{}{
		"model":       model,
		"temperature": temperature,
		"messages": []map[string]string{
			{"role": "system", "content": `
				You are an assistant designed to interpret IAM-related requests and convert them into structured JSON objects.
	
				Your task is to:
				1. Analyze the user's input.
				2. Only provide a field if you are very certain (confidence >= 8). If you are not sure, return null for that field.
				3. Identify the requested action, target entity, and resource details. 
				4. If input is vauge return null for specified section.
				5. If target is a specific resource use custom policies. "Example: For 'Allow access to bucket production-data', identify the required resource-specific policy and set isManagedPolicy = false.
				6. Identify policy. If target is a service try getting aws managed policies first. Be certain about aws managed policy names if its wrong correct it.
				7. If the identified policy name has low confidence, set confidence < 5."
				8. If the user input does not provide any meaningful context, the model must not guess a policy.
			`},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]interface{}{
			"type":        "json_schema",
			"json_schema": schema,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return GPTResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return GPTResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GPTResponse{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return GPTResponse{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response into an intermediate structure
	var intermediateResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &intermediateResponse)
	if err != nil {
		return GPTResponse{}, fmt.Errorf("failed to parse intermediate response: %w", err)
	}

	if len(intermediateResponse.Choices) == 0 || intermediateResponse.Choices[0].Message.Content == "" {
		return GPTResponse{}, fmt.Errorf("no content found in response")
	}

	// Extract the actual structured content
	content := intermediateResponse.Choices[0].Message.Content
	content = strings.TrimSpace(content)

	// Unmarshal the structured JSON content into GPTResponse
	var gptResponse GPTResponse
	err = json.Unmarshal([]byte(content), &gptResponse)
	if err != nil {
		return GPTResponse{}, fmt.Errorf("failed to parse structured content: %w", err)
	}

	return gptResponse, nil
}

func GenerateCLICommand(response GPTResponse) string {
	placeholder := "?"

	principalType := placeholder
	principalName := placeholder
	if response.Principal != nil {
		if val, ok := response.Principal["type"]; ok {
			if val != "" {
				principalType = val
			}
		}
		if val, ok := response.Principal["name"]; ok {
			if val != "" {
				principalName = val
			}
		}
	}

	policyName := placeholder
	if response.Policy != "" {
		policyName = response.Policy
	}

	action := response.Action
	if action == "" {
		action = placeholder
	}

	return fmt.Sprintf("aws %s %s %s %s", principalType, principalName, action, policyName)
}

func main() {
	apiKey := "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"
	model := "gpt-4o"
	temperature := 0.1
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your request (or 'exit' to quit): ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if strings.ToLower(userInput) == "exit" {
			break
		}

		prompt := generatePrompt(userInput)
		response, err := callGPTWithSchema(apiKey, model, prompt, temperature)
		if err != nil {
			fmt.Println("Error calling GPT:", err)
			continue
		}

		cliCommand := GenerateCLICommand(response)
		fmt.Println("CLI Command:", cliCommand)
		fmt.Printf("Parsed Response: %+v\n", response)
	}
}
