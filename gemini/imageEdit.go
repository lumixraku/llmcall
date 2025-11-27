package gemini

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// ImageEditResult represents the result of an image edit operation
type ImageEditResult struct {
	Format string // Image format (e.g., "png", "jpeg")
	Data   []byte // Decoded image data
}

// ImageEdit takes a base64 encoded image string and a prompt, returns edited images
func ImageEdit(imageBase64 string, prompt string) ([]ImageEditResult, error) {
	// Load .env file from current directory
	_ = godotenv.Load()

	apiKey := os.Getenv("SONETTO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SONETTO_API_KEY is not set in the environment")
	}

	// Create OpenAI client with custom base URL
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://vip.sonetto.top/v1"),
	)

	// Create chat completion request with image
	ctx := context.Background()
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "gemini-3-pro-image-preview-16",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				{
					OfImageURL: &openai.ChatCompletionContentPartImageParam{
						ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
							URL: fmt.Sprintf("data:image/png;base64,%s", imageBase64),
						},
					},
				},
				{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: prompt,
					},
				},
			}),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	// Process response
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	content := resp.Choices[0].Message.Content

	// Try to extract base64 image from markdown format: ![image](data:image/png;base64,...)
	pattern := regexp.MustCompile(`!\[.*?\]\(data:image/(\w+);base64,([A-Za-z0-9+/=]+)\)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no image found in response: %s", truncateString(content, 200))
	}

	var results []ImageEditResult
	for idx, match := range matches {
		imgFormat := match[1]
		b64Data := match[2]

		imageBytes, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image %d: %w", idx, err)
		}

		results = append(results, ImageEditResult{
			Format: imgFormat,
			Data:   imageBytes,
		})
	}

	return results, nil
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// loadImageAsBase64 loads an image file and returns its base64 encoded string
func loadImageAsBase64(imagePath string) (string, error) {
	// Get absolute path relative to current directory
	absPath := imagePath

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
