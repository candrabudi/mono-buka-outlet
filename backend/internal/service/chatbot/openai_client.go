package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ──────────────────────────────────────────────────────────────
// OpenAI API Client
// Supports both Chat Completions & Responses API (web search)
// ──────────────────────────────────────────────────────────────

type OpenAIClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// ──────────────────────────────────────────────────────────────
// Responses API types (for web search)
// ──────────────────────────────────────────────────────────────

type ResponsesAPIRequest struct {
	Model       string                   `json:"model"`
	Tools       []map[string]interface{} `json:"tools"`
	Input       []OpenAIMessage          `json:"input"`
	Temperature float64                  `json:"temperature,omitempty"`
}

type ResponsesAPIResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output []struct {
		ID      string `json:"id"`
		Type    string `json:"type"` // "message", "web_search_call"
		Role    string `json:"role,omitempty"`
		Content []struct {
			Type        string             `json:"type"` // "output_text", "refusal"
			Text        string             `json:"text,omitempty"`
			Annotations []SearchAnnotation `json:"annotations,omitempty"`
		} `json:"content,omitempty"`
		Status string `json:"status,omitempty"` // for web_search_call
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

type SearchAnnotation struct {
	Type  string `json:"type"` // "url_citation"
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
	Start int    `json:"start_index,omitempty"`
	End   int    `json:"end_index,omitempty"`
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if model == "" {
		model = "gpt-4o" // default fallback
	}
	return &OpenAIClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.openai.com/v1/chat/completions",
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

// ChatCompletion sends a chat completion request to OpenAI (no web search)
func (c *OpenAIClient) ChatCompletion(messages []OpenAIMessage, temperature float64, maxTokens int) (*OpenAIResponse, error) {
	reqBody := OpenAIRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		TopP:        0.95,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result OpenAIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("OpenAI error [%s]: %s", result.Error.Type, result.Error.Message)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OpenAI returned status %d: %s", resp.StatusCode, string(body))
	}

	return &result, nil
}

// ChatCompletionWithSearch uses the OpenAI Responses API with web_search_preview tool.
// This lets GPT search the web directly and include real-time information.
func (c *OpenAIClient) ChatCompletionWithSearch(messages []OpenAIMessage, temperature float64) (*OpenAIResponse, []SearchAnnotation, error) {
	reqBody := ResponsesAPIRequest{
		Model: c.model,
		Tools: []map[string]interface{}{
			{
				"type":                "web_search_preview",
				"search_context_size": "medium",
			},
		},
		Input:       messages,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result ResponsesAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to parse Responses API response: %w (body: %s)", err, string(body))
	}

	if result.Error != nil {
		return nil, nil, fmt.Errorf("OpenAI error [%s]: %s", result.Error.Type, result.Error.Message)
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("OpenAI returned status %d: %s", resp.StatusCode, string(body))
	}

	// Extract the reply text and annotations from the Responses API output
	replyText := ""
	var allAnnotations []SearchAnnotation

	for _, item := range result.Output {
		if item.Type == "message" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					replyText += content.Text
					allAnnotations = append(allAnnotations, content.Annotations...)
				}
			}
		}
	}

	// Convert to standard OpenAIResponse for compatibility
	standardResp := &OpenAIResponse{
		ID:    result.ID,
		Model: c.model,
		Choices: []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{
					Role:    "assistant",
					Content: replyText,
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     result.Usage.InputTokens,
			CompletionTokens: result.Usage.OutputTokens,
			TotalTokens:      result.Usage.TotalTokens,
		},
	}

	return standardResp, allAnnotations, nil
}

// GetReply extracts the reply text from the response
func (c *OpenAIClient) GetReply(resp *OpenAIResponse) string {
	if resp != nil && len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return ""
}
