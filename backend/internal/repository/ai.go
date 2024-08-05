package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type AIRepository interface {
	GenerateAbstractAndKeyword(ctx context.Context, content string) (string, string, error)
}

func NewAIRepository(
	r *Repository,
) AIRepository {
	return &aiRepository{
		Repository: r,
	}
}

type aiRepository struct {
	*Repository
}

// Define the structure of the response based on the provided JSON format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

// Function to extract abstract and keywords
func extractAbstractAndKeywords(content string) (string, string) {
	// Define regex patterns for abstract and keywords
	abstractPattern := `摘要：([\s\S]*?)\n`
	keywordsPattern := `关键字：([\s\S]*)`

	// Compile the regex patterns
	abstractRegex := regexp.MustCompile(abstractPattern)
	keywordsRegex := regexp.MustCompile(keywordsPattern)

	// Find abstract
	abstractMatch := abstractRegex.FindStringSubmatch(content)
	var abstract string
	if len(abstractMatch) > 1 {
		abstract = abstractMatch[1]
	}

	// Find keywords
	keywordsMatch := keywordsRegex.FindStringSubmatch(content)
	var keywords string
	if len(keywordsMatch) > 1 {
		keywords = keywordsMatch[1]
	}

	// Clean up keywords: split by both English and Chinese commas, remove spaces, and join with English commas
	keywordList := strings.FieldsFunc(keywords, func(r rune) bool {
		return r == ',' || r == '，'
	})
	for i, keyword := range keywordList {
		keywordList[i] = strings.TrimSpace(keyword)
	}
	formattedKeywords := strings.Join(keywordList, ",")

	return abstract, formattedKeywords
}

func (r *aiRepository) GenerateAbstractAndKeyword(ctx context.Context, content string) (string, string, error) {
	aiGenerate := os.Getenv("ai_generate")
	if aiGenerate == "" {
		aiGenerate = r.conf.GetString("ai.generate")
	}
	if aiGenerate == "close" {
		return "", "", nil
	}

	// Check if content is empty or starts with "Error processing link"
	if strings.TrimSpace(content) == "" || strings.HasPrefix(content, "Error processing link") {
		return "", "", nil
	}

	// Define the request payload
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "system", "content": "你是一位优秀的文章总结高手。请根据我提供的文字，总结出一个简洁的摘要，要求在100字以内，并提取出至少3个关键字。请将结果以以下格式返回：\n\n摘要：xxx\n关键字：xxx,xxx,xxx"},
			{"role": "user", "content": content},
		},
		"stream": false,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %v", err)
	}

	apiKey := os.Getenv("api_key")
	if apiKey == "" {
		apiKey = r.conf.GetString("open.api_key")
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response JSON
	var apiResponse Response
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extract the content where role is "assistant"
	for _, choice := range apiResponse.Choices {
		if choice.Message.Role == "assistant" {
			abstract, keyword := extractAbstractAndKeywords(choice.Message.Content)
			return abstract, keyword, nil
		}
	}

	return "", "", fmt.Errorf("no assistant role found in the response")
}
