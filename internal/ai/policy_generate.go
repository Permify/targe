package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type IAMPolicy struct {
	Version   string         `json:"Version"`
	Id        string         `json:"Id,omitempty"`
	Statement []IAMStatement `json:"Statement"`
}

type IAMStatement struct {
	Sid          string                 `json:"Sid,omitempty"`
	Effect       string                 `json:"Effect"`
	Principal    *IAMPrincipal          `json:"Principal,omitempty"`
	NotPrincipal *IAMPrincipal          `json:"NotPrincipal,omitempty"`
	Action       *IAMActionResource     `json:"Action,omitempty"`
	NotAction    *IAMActionResource     `json:"NotAction,omitempty"`
	Resource     *IAMActionResource     `json:"Resource,omitempty"`
	NotResource  *IAMActionResource     `json:"NotResource,omitempty"`
	Condition    map[string]interface{} `json:"Condition,omitempty"`
}

type IAMPrincipal struct {
	Star          bool
	AWS           []string
	Federated     []string
	Service       []string
	CanonicalUser []string
}

type IAMActionResource struct {
	IsStar    bool
	Resources []string
}

var IAMPolicySchema = map[string]interface{}{
	"name": "IAMPolicy",
	"schema": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"Version": map[string]interface{}{
				"type": "string",
				"enum": []string{"2012-10-17"},
			},
			"Id": map[string]interface{}{
				"type": "string",
			},
			"Statement": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"Sid": map[string]interface{}{
							"type": "string",
						},
						"Effect": map[string]interface{}{
							"type": "string",
							"enum": []string{"Allow", "Deny"},
						},
						"Principal": map[string]interface{}{
							"oneOf": []interface{}{
								map[string]interface{}{
									"type": "string",
									"enum": []string{"*"},
								},
								map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"AWS": map[string]interface{}{
											"oneOf": []interface{}{
												map[string]interface{}{
													"type": "string",
												},
												map[string]interface{}{
													"type":  "array",
													"items": map[string]interface{}{"type": "string"},
												},
											},
										},
									},
									"additionalProperties": false,
								},
							},
						},
						"Action": map[string]interface{}{
							"oneOf": []interface{}{
								map[string]interface{}{
									"type": "string",
									"enum": []string{"*"},
								},
								map[string]interface{}{
									"type":  "array",
									"items": map[string]interface{}{"type": "string"},
								},
							},
						},
						"Resource": map[string]interface{}{
							"oneOf": []interface{}{
								map[string]interface{}{
									"type": "string",
									"enum": []string{"*"},
								},
								map[string]interface{}{
									"type":  "array",
									"items": map[string]interface{}{"type": "string"},
								},
							},
						},
					},
					"required": []string{"Effect"},
				},
			},
		},
		"required": []string{"Version", "Statement"},
	},
}

func (p *IAMPrincipal) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "*" {
			p.Star = true
			return nil
		}
		p.AWS = []string{str}
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "AWS":
			p.AWS = toStringSlice(v)
		case "Federated":
			p.Federated = toStringSlice(v)
		case "Service":
			p.Service = toStringSlice(v)
		case "CanonicalUser":
			p.CanonicalUser = toStringSlice(v)
		}
	}
	return nil
}

func (ar *IAMActionResource) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "*" {
			ar.IsStar = true
			return nil
		}
		ar.Resources = []string{str}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		ar.Resources = arr
		return nil
	}

	return fmt.Errorf("invalid action/resource format")
}

func (p IAMPrincipal) MarshalJSON() ([]byte, error) {
	if p.Star {
		return json.Marshal("*")
	}

	obj := map[string]interface{}{}

	if len(p.AWS) == 1 {
		obj["AWS"] = p.AWS[0]
	} else if len(p.AWS) > 1 {
		obj["AWS"] = p.AWS
	}

	if len(p.Federated) == 1 {
		obj["Federated"] = p.Federated[0]
	} else if len(p.Federated) > 1 {
		obj["Federated"] = p.Federated
	}

	if len(p.Service) == 1 {
		obj["Service"] = p.Service[0]
	} else if len(p.Service) > 1 {
		obj["Service"] = p.Service
	}

	if len(p.CanonicalUser) == 1 {
		obj["CanonicalUser"] = p.CanonicalUser[0]
	} else if len(p.CanonicalUser) > 1 {
		obj["CanonicalUser"] = p.CanonicalUser
	}

	if len(obj) == 0 {
		return json.Marshal("*")
	}

	return json.Marshal(obj)
}

func (ar IAMActionResource) MarshalJSON() ([]byte, error) {
	if ar.IsStar {
		return json.Marshal("*")
	}
	if len(ar.Resources) == 1 {
		return json.Marshal(ar.Resources[0])
	}
	return json.Marshal(ar.Resources)
}

func toStringSlice(v interface{}) []string {
	switch val := v.(type) {
	case string:
		return []string{val}
	case []interface{}:
		var out []string
		for _, item := range val {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func GeneratePolicy(apiKey, prompt string) (IAMPolicy, error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model":       "gpt-4o",
		"temperature": 0.1,
		"messages": []map[string]string{
			{"role": "system", "content": "You are an assistant that produces IAM policies as JSON."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]interface{}{
			"type":        "json_schema",
			"json_schema": IAMPolicySchema,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return IAMPolicy{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return IAMPolicy{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return IAMPolicy{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return IAMPolicy{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

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
		return IAMPolicy{}, fmt.Errorf("failed to parse intermediate response: %w", err)
	}

	if len(intermediateResponse.Choices) == 0 || intermediateResponse.Choices[0].Message.Content == "" {
		return IAMPolicy{}, fmt.Errorf("no content found in response")
	}

	content := strings.TrimSpace(intermediateResponse.Choices[0].Message.Content)

	var policy IAMPolicy
	err = json.Unmarshal([]byte(content), &policy)
	if err != nil {
		return IAMPolicy{}, fmt.Errorf("failed to parse structured content: %w", err)
	}

	return policy, nil
}
