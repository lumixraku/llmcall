package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

// 模拟工具函数
func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getWeather(location string) string {
	// 模拟天气数据
	return fmt.Sprintf(`{"location": "%s", "temperature": "22°C", "weather": "晴", "humidity": "45%%"}`, location)
}

func generateImage(prompt string) string {
	// 模拟图片生成
	return fmt.Sprintf(`{"image_url": "https://example.com/generated-image.png", "prompt": "%s"}`, prompt)
}

// 执行工具调用
func executeTool(name string, arguments string) string {
	var args map[string]any
	json.Unmarshal([]byte(arguments), &args)

	switch name {
	case "getCurrentTime":
		return getCurrentTime()
	case "getWeather":
		location, _ := args["location"].(string)
		return getWeather(location)
	case "generateImage":
		prompt, _ := args["prompt"].(string)
		return generateImage(prompt)
	default:
		return fmt.Sprintf(`{"error": "unknown tool: %s"}`, name)
	}
}

func main() {

	// https://platform.openai.com/account/api-keys
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	apiKey := os.Getenv("MOONSHOT_API_KEY")
	log.Printf("Using API Key: %s...%s", apiKey[:10], apiKey[len(apiKey)-4:])

	// 使用 Kimi (月之暗面) OpenAI 兼容接口
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.moonshot.cn/v1"),
	)

	ctx := context.Background()

	// 定义工具
	tools := []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "getCurrentTime",
				Description: param.Opt[string]{Value: "获取当前时间"},
				Parameters: openai.FunctionParameters{
					"type": "object",
				},
			},
		},
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "getWeather",
				Description: param.Opt[string]{Value: "获取某地天气"},
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{"type": "string"},
					},
					"required": []string{"location"},
				},
			},
		},
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "generateImage",
				Description: param.Opt[string]{Value: "根据描述生成图片"},
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"prompt": map[string]any{"type": "string"},
					},
					"required": []string{"prompt"},
				},
			},
		},
	}

	// 初始消息
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("请获取深圳的天气和当前时间，然后根据天气情况生成一张符合天气风格的图片，最后做成图文简报。"),
	}

	// Agentic Loop: 循环直到模型不再调用工具
	for {
		log.Printf("发送请求，消息数: %d", len(messages))
		msgsJSON, _ := json.MarshalIndent(messages, "", "  ")
		log.Printf("messages: %s", string(msgsJSON))

		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    "moonshot-v1-8k",
			Messages: messages,
			Tools:    tools,
		})
		if err != nil {
			panic(err)
		}

		choice := resp.Choices[0]
		log.Printf("收到响应，finish_reason: %s", choice.FinishReason)

		// 如果没有 tool_calls，结束循环
		if choice.FinishReason != "tool_calls" || len(choice.Message.ToolCalls) == 0 {
			fmt.Println("\n===== 最终回复 =====")
			fmt.Println(choice.Message.Content)
			break
		}

		// 打印完整响应用于调试
		choiceJSON, _ := json.MarshalIndent(choice, "", "  ")
		log.Printf("choice: %s", string(choiceJSON))

		// 处理所有 tool calls
		for _, toolCall := range choice.Message.ToolCalls {
			log.Printf("调用工具: %s (id=%s), 参数: %s", toolCall.Function.Name, toolCall.ID, toolCall.Function.Arguments)

			result := executeTool(toolCall.Function.Name, toolCall.Function.Arguments)
			log.Printf("工具返回: %s", result)

			// 将 assistant 的 tool_call 和 tool 结果一起加入
			messages = append(messages,
				openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: []openai.ChatCompletionMessageToolCallParam{
							toolCall.ToParam(),
						},
					},
				},
				openai.ToolMessage(result, toolCall.ID),
			)
		}
	}
}