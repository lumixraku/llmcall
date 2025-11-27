package gemini

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func ImageEdit() {
	// Load .env file from current directory
	_ = godotenv.Load()

	apiKey := os.Getenv("SONETTO_API_KEY")
	if apiKey == "" {
		fmt.Println("SONETTO_API_KEY is not set in the environment")
		os.Exit(1)
	}

	// Create OpenAI client with custom base URL
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://vip.sonetto.top/v1"),
	)

	// Load image as base64
	imageData, err := loadImageAsBase64("assets/two.png")
	if err != nil {
		fmt.Printf("Failed to load image: %v\n", err)
		os.Exit(1)
	}

	// Create chat completion request with image
	ctx := context.Background()
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "gemini-3-pro-image-preview-16",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				{
					OfImageURL: &openai.ChatCompletionContentPartImageParam{
						ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
							URL: fmt.Sprintf("data:image/png;base64,%s", imageData),
						},
					},
				},
				{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: "Japanese Animate Style",
					},
				},
			}),
		},
	})
	if err != nil {
		fmt.Printf("API request failed: %v\n", err)
		os.Exit(1)
	}

	// Process response
	if len(resp.Choices) == 0 {
		fmt.Println("No response choices returned")
		os.Exit(1)
	}

	content := resp.Choices[0].Message.Content

	// Try to extract base64 image from markdown format: ![image](data:image/png;base64,...)
	pattern := regexp.MustCompile(`!\[.*?\]\(data:image/(\w+);base64,([A-Za-z0-9+/=]+)\)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	if len(matches) > 0 {
		for idx, match := range matches {
			imgFormat := match[1]
			b64Data := match[2]

			imageBytes, err := base64.StdEncoding.DecodeString(b64Data)
			if err != nil {
				fmt.Printf("Failed to decode base64 image %d: %v\n", idx, err)
				continue
			}

			outputPath := filepath.Join("assets", "output", fmt.Sprintf("edited_%d.%s", idx, imgFormat))
			// Ensure output directory exists
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				fmt.Printf("Failed to create output directory: %v\n", err)
				continue
			}

			if err := os.WriteFile(outputPath, imageBytes, 0644); err != nil {
				fmt.Printf("Failed to save image %d: %v\n", idx, err)
				continue
			}
			fmt.Printf("编辑后的图片已保存到: %s\n", outputPath)
		}
	} else {
		// Print raw content if no image found
		if len(content) > 200 {
			fmt.Printf("未能解析图片数据: %s...\n", content[:200])
		} else {
			fmt.Printf("未能解析图片数据: %s\n", content)
		}
	}
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
