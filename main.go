package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Perplexity API supported models
const (
	// Search models
	ModelSonarPro = "sonar-pro" // Advanced search offering with grounding
	ModelSonar    = "sonar"     // Lightweight, cost-effective search model

	// Research models
	ModelDeepResearch = "sonar-deep-research" // Expert-level research model for comprehensive reports

	// Reasoning models
	ModelReasoningPro = "sonar-reasoning-pro" // Premier reasoning model with Chain of Thought
	ModelReasoning    = "sonar-reasoning"     // Fast, real-time reasoning model

	// Offline models
	ModelR1 = "r1-1776" // A version of DeepSeek R1, post-trained for uncensored, unbiased and factual information
)

// Perplexity API endpoint
const (
	perplexityAPIURL = "https://api.perplexity.ai/chat/completions"
)

// PerplexitySearchResult defines the structure for search results
type PerplexitySearchResult struct {
	Query            string   `json:"query"`
	Text             string   `json:"text"`
	Citations        []string `json:"citations,omitempty"`
	RelatedQuestions []string `json:"related_questions,omitempty"`
	Images           []string `json:"images,omitempty"`
}

// PerplexityChatMessage represents a message in the chat
type PerplexityChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PerplexityChatRequest defines the request structure for the Perplexity Chat API
type PerplexityChatRequest struct {
	Model                  string                  `json:"model"`
	Messages               []PerplexityChatMessage `json:"messages"`
	MaxTokens              int                     `json:"max_tokens,omitempty"`
	Temperature            float64                 `json:"temperature,omitempty"`
	TopP                   float64                 `json:"top_p,omitempty"`
	SearchDomainFilter     []string                `json:"search_domain_filter,omitempty"`
	ReturnImages           bool                    `json:"return_images,omitempty"`
	ReturnRelatedQuestions bool                    `json:"return_related_questions,omitempty"`
	SearchRecencyFilter    string                  `json:"search_recency_filter,omitempty"`
	TopK                   int                     `json:"top_k,omitempty"`
	Stream                 bool                    `json:"stream,omitempty"`
	PresencePenalty        float64                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty       float64                 `json:"frequency_penalty,omitempty"`
	ResponseFormat         *struct{}               `json:"response_format,omitempty"`
	WebSearchOptions       *struct {
		SearchContextSize string `json:"search_context_size,omitempty"`
	} `json:"web_search_options,omitempty"`
}

// PerplexityChatResponse represents the API response
type PerplexityChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
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
}

// buildSearchTool constructs the MCP tool definition for Perplexity search
func buildSearchTool() mcp.Tool {
	return mcp.NewTool("perplexity_search",
		mcp.WithDescription("Perform web search using Perplexity API and return results"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query string"),
		),
		mcp.WithString("model",
			mcp.Description("Model to use for the search. Options include:\n"+
				"- sonar-pro: Advanced search offering with grounding (default)\n"+
				"- sonar: Lightweight, cost-effective search model\n"+
				"- sonar-deep-research: Expert-level research model for comprehensive reports\n"+
				"- sonar-reasoning-pro: Premier reasoning model with Chain of Thought\n"+
				"- sonar-reasoning: Fast, real-time reasoning model\n"+
				"- r1-1776: Offline chat model (no search capability)"),
			mcp.Enum(ModelSonarPro, ModelSonar, ModelDeepResearch, ModelReasoningPro, ModelReasoning, ModelR1),
		),
		mcp.WithString("search_recency_filter",
			mcp.Description("Filter search results by recency (options: month, week, day, hour)"),
			mcp.Enum("month", "week", "day", "hour"),
		),
		mcp.WithNumber("max_tokens",
			mcp.Description("Maximum number of tokens returned by the API (max 8k for sonar-pro)"),
		),
		mcp.WithNumber("temperature",
			mcp.Description("Amount of randomness in the response, valued between 0 and 2"),
			mcp.DefaultNumber(0.2),
		),
		mcp.WithNumber("top_p",
			mcp.Description("Nucleus sampling threshold, valued between 0 and 1"),
			mcp.DefaultNumber(0.9),
		),
		mcp.WithArray("search_domain_filter",
			mcp.Description("List of domains to limit search results to"),
		),
		mcp.WithBoolean("return_images",
			mcp.Description("Whether search results should include images"),
		),
		mcp.WithBoolean("return_related_questions",
			mcp.Description("Whether related questions should be returned"),
		),
		mcp.WithNumber("top_k",
			mcp.Description("Number of tokens to keep for top-k filtering"),
			mcp.DefaultNumber(0),
		),
		mcp.WithBoolean("stream",
			mcp.Description("Whether to stream the response incrementally"),
		),
		mcp.WithNumber("presence_penalty",
			mcp.Description("Positive values increase the likelihood of discussing new topics"),
			mcp.DefaultNumber(0),
		),
		mcp.WithNumber("frequency_penalty",
			mcp.Description("Decreases likelihood of repetition based on prior frequency"),
			mcp.DefaultNumber(1),
		),
		mcp.WithObject("web_search_options",
			mcp.Description("Configuration for using web search in model responses. The 'search_context_size' property can be set to 'low', 'medium', or 'high' to control how much search context is retrieved (default: medium)"),
		),
	)
}

// handleSearchTool processes search tool requests and returns results
func handleSearchTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY environment variable not set")
	}

	// Extract required parameters
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("invalid or empty query parameter")
	}

	// Prepare default values
	defaultMaxTokens := 2000             // Default max tokens
	defaultModel := ModelSonarPro        // Default model
	defaultReturnImages := false         // Default for return images
	defaultReturnRelatedQs := false      // Default for return related questions
	defaultTemperature := 0.2            // Default temperature
	defaultTopP := 0.9                   // Default TopP
	defaultTopK := 0                     // Default TopK
	defaultPresencePenalty := 0.0        // Default presence penalty
	defaultFrequencyPenalty := 1.0       // Default frequency penalty
	defaultStream := false               // Default stream setting
	defaultSearchContextSize := "medium" // Default search context size

	// Check environment variables for default overrides
	if envMaxTokens := os.Getenv("DEFAULT_MAX_TOKENS"); envMaxTokens != "" {
		if _, err := fmt.Sscanf(envMaxTokens, "%d", &defaultMaxTokens); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Invalid DEFAULT_MAX_TOKENS value: %s\n", envMaxTokens)
		}
	}

	if envModel := os.Getenv("PERPLEXITY_MODEL"); envModel != "" {
		defaultModel = envModel
	}

	if envContextSize := os.Getenv("DEFAULT_SEARCH_CONTEXT_SIZE"); envContextSize != "" {
		defaultSearchContextSize = envContextSize
	}

	// Prepare the chat request
	chatReq := PerplexityChatRequest{
		Model: defaultModel,
		Messages: []PerplexityChatMessage{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens:              defaultMaxTokens,
		Temperature:            defaultTemperature,
		TopP:                   defaultTopP,
		TopK:                   defaultTopK,
		ReturnImages:           defaultReturnImages,
		ReturnRelatedQuestions: defaultReturnRelatedQs,
		Stream:                 defaultStream,
		PresencePenalty:        defaultPresencePenalty,
		FrequencyPenalty:       defaultFrequencyPenalty,
		WebSearchOptions: &struct {
			SearchContextSize string `json:"search_context_size,omitempty"`
		}{
			SearchContextSize: defaultSearchContextSize,
		},
	}

	// Validate model selection
	if model, ok := request.Params.Arguments["model"].(string); ok && model != "" {
		// Verify the model is valid
		validModel := false
		validModels := []string{
			ModelSonarPro, ModelSonar, ModelDeepResearch,
			ModelReasoningPro, ModelReasoning, ModelR1,
		}
		for _, m := range validModels {
			if model == m {
				validModel = true
				break
			}
		}

		if validModel {
			chatReq.Model = model
		} else {
			return nil, fmt.Errorf("invalid model: %s", model)
		}
	}

	// Set optional parameters (override defaults)
	if maxTokens, ok := request.Params.Arguments["max_tokens"].(float64); ok {
		chatReq.MaxTokens = int(maxTokens)
	}

	if temperature, ok := request.Params.Arguments["temperature"].(float64); ok {
		chatReq.Temperature = temperature
	}

	if topP, ok := request.Params.Arguments["top_p"].(float64); ok {
		chatReq.TopP = topP
	}

	// Process search_domain_filter parameter
	if domainFilter, ok := request.Params.Arguments["search_domain_filter"].([]interface{}); ok && len(domainFilter) > 0 {
		chatReq.SearchDomainFilter = make([]string, len(domainFilter))
		for i, domain := range domainFilter {
			if domain, ok := domain.(string); ok {
				chatReq.SearchDomainFilter[i] = domain
			}
		}
	}

	// Process return_images parameter
	if returnImages, ok := request.Params.Arguments["return_images"].(bool); ok {
		chatReq.ReturnImages = returnImages
	}

	// Process return_related_questions parameter
	if returnRelatedQuestions, ok := request.Params.Arguments["return_related_questions"].(bool); ok {
		chatReq.ReturnRelatedQuestions = returnRelatedQuestions
	}

	// Process search_recency_filter parameter
	if recencyFilter, ok := request.Params.Arguments["search_recency_filter"].(string); ok && recencyFilter != "" {
		chatReq.SearchRecencyFilter = recencyFilter
	}

	if topK, ok := request.Params.Arguments["top_k"].(float64); ok {
		chatReq.TopK = int(topK)
	}

	if stream, ok := request.Params.Arguments["stream"].(bool); ok {
		chatReq.Stream = stream
	}

	if presencePenalty, ok := request.Params.Arguments["presence_penalty"].(float64); ok {
		chatReq.PresencePenalty = presencePenalty
	}

	if frequencyPenalty, ok := request.Params.Arguments["frequency_penalty"].(float64); ok {
		chatReq.FrequencyPenalty = frequencyPenalty
	}

	// Process response_format parameter
	if _, ok := request.Params.Arguments["response_format"].(map[string]interface{}); ok {
		chatReq.ResponseFormat = &struct{}{}
	}

	// Process web_search_options parameter
	if webSearchOptions, ok := request.Params.Arguments["web_search_options"].(map[string]interface{}); ok {
		if searchContextSize, ok := webSearchOptions["search_context_size"].(string); ok && searchContextSize != "" {
			// Validate search_context_size value
			if searchContextSize == "low" || searchContextSize == "medium" || searchContextSize == "high" {
				chatReq.WebSearchOptions.SearchContextSize = searchContextSize
			}
		}
	}

	// For offline model r1-1776, remove search-related parameters
	if chatReq.Model == ModelR1 {
		// Clear search-related parameters
		chatReq.SearchDomainFilter = nil
		chatReq.ReturnImages = false
		chatReq.ReturnRelatedQuestions = false
		chatReq.SearchRecencyFilter = ""
		chatReq.WebSearchOptions = nil
	}

	// Execute the chat request
	responseText, err := performPerplexityChat(apiKey, chatReq)
	if err != nil {
		return nil, fmt.Errorf("search execution failed: %v", err)
	}

	return mcp.NewToolResultText(responseText), nil
}

// performPerplexityChat executes a chat request to the Perplexity API
func performPerplexityChat(apiKey string, req PerplexityChatRequest) (string, error) {
	// Serialize the request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", perplexityAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set request headers
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned error status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var chatResp PerplexityChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract content from the response
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"perplexity-search",
		"1.0.0",
		server.WithLogging(),
	)

	// Add search tool
	s.AddTool(buildSearchTool(), handleSearchTool)

	// Start the server (using stdin/stdout as communication channel)
	fmt.Fprintf(os.Stderr, "Perplexity Search MCP Server running...\n")
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
