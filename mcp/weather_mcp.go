package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

// MCPToolManager 管理多个 MCP 服务器连接
type MCPToolManager struct {
	sessions map[string]*mcp.ClientSession
	tools    map[string]*mcp.Tool // tool name -> tool definition
	toolMap  map[string]string    // tool name -> session name
}

func NewMCPToolManager() *MCPToolManager {
	return &MCPToolManager{
		sessions: make(map[string]*mcp.ClientSession),
		tools:    make(map[string]*mcp.Tool),
		toolMap:  make(map[string]string),
	}
}

// ConnectToLocalServer 连接到本地 MCP 服务器 (通过命令启动)
func (m *MCPToolManager) ConnectToLocalServer(ctx context.Context, name string, command string, args ...string) error {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "llmcall-client",
		Version: "v1.0.0",
	}, nil)

	cmd := exec.Command(command, args...)
	transport := &mcp.CommandTransport{Command: cmd}

	return m.connectAndRegister(ctx, client, transport, name)
}

// ConnectToRemoteServer 连接到远程 MCP 服务器 (通过 SSE)
// 示例: ConnectToRemoteServer(ctx, "deepwiki", "https://mcp.deepwiki.com/sse")
func (m *MCPToolManager) ConnectToRemoteServer(ctx context.Context, name string, endpoint string) error {
	log.Printf("[MCP] 连接远程服务器: %s -> %s", name, endpoint)

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "llmcall-client",
		Version: "v1.0.0",
	}, nil)

	transport := &mcp.SSEClientTransport{
		Endpoint: endpoint,
	}

	return m.connectAndRegister(ctx, client, transport, name)
}

// connectAndRegister 通用的连接和注册工具逻辑
func (m *MCPToolManager) connectAndRegister(ctx context.Context, client *mcp.Client, transport mcp.Transport, name string) error {
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server %s: %w", name, err)
	}

	m.sessions[name] = session
	log.Printf("[MCP] 已连接到服务器: %s", name)

	// 获取该服务器提供的所有工具
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list tools from %s: %w", name, err)
	}

	for _, tool := range toolsResult.Tools {
		m.tools[tool.Name] = tool
		m.toolMap[tool.Name] = name
		log.Printf("[MCP] 注册工具: %s (from %s) - %s", tool.Name, name, tool.Description)
	}

	return nil
}

// CallTool 调用指定的 MCP 工具
func (m *MCPToolManager) CallTool(ctx context.Context, toolName string, arguments map[string]any) (string, error) {
	sessionName, ok := m.toolMap[toolName]
	if !ok {
		return "", fmt.Errorf("tool %s not found", toolName)
	}

	argsJSON, _ := json.Marshal(arguments)
	log.Printf("[MCP] 调用工具: %s (server=%s), 参数: %s", toolName, sessionName, string(argsJSON))

	session := m.sessions[sessionName]
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	})
	if err != nil {
		log.Printf("[MCP] 工具 %s 调用失败: %v", toolName, err)
		return "", err
	}

	if result.IsError {
		log.Printf("[MCP] 工具 %s 返回错误", toolName)
		return "", fmt.Errorf("tool %s returned error", toolName)
	}

	// 提取文本内容
	var contents []string
	for _, c := range result.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			contents = append(contents, tc.Text)
		}
	}

	output := fmt.Sprintf("%v", contents)
	log.Printf("[MCP] 工具 %s 返回: %s", toolName, output)
	return output, nil
}

// ToOpenAITools 将 MCP 工具转换为 OpenAI 工具格式
func (m *MCPToolManager) ToOpenAITools() []openai.ChatCompletionToolParam {
	var tools []openai.ChatCompletionToolParam

	for _, tool := range m.tools {
		// 将 MCP 的 InputSchema 转换为 OpenAI 的 FunctionParameters
		params := openai.FunctionParameters{"type": "object"}

		if tool.InputSchema != nil {
			// InputSchema 是 any 类型，需要通过 JSON 转换
			data, err := json.Marshal(tool.InputSchema)
			if err == nil {
				var schemaMap map[string]any
				if json.Unmarshal(data, &schemaMap) == nil {
					params = schemaMap
				}
			}
		}

		tools = append(tools, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: param.Opt[string]{Value: tool.Description},
				Parameters:  params,
			},
		})
	}

	return tools
}

// Close 关闭所有 MCP 会话
func (m *MCPToolManager) Close() {
	for name, session := range m.sessions {
		session.Close()
		log.Printf("Closed MCP session: %s", name)
	}
}

// 本地工具函数 (非 MCP)
func localGetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func localGenerateImage(prompt string) string {
	return fmt.Sprintf(`{"image_url": "https://example.com/generated-image.png", "prompt": "%s"}`, prompt)
}

func executeLocalTool(name string, arguments string) string {
	var args map[string]any
	json.Unmarshal([]byte(arguments), &args)

	switch name {
	case "getCurrentTime":
		return localGetCurrentTime()
	case "generateImage":
		prompt, _ := args["prompt"].(string)
		return localGenerateImage(prompt)
	default:
		return fmt.Sprintf(`{"error": "unknown local tool: %s"}`, name)
	}
}

func main() {
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	apiKey := os.Getenv("MOONSHOT_API_KEY")
	log.Printf("Using API Key: %s...%s", apiKey[:10], apiKey[len(apiKey)-4:])

	ctx := context.Background()

	// 创建 MCP 工具管理器
	mcpManager := NewMCPToolManager()
	defer mcpManager.Close()

	// ========== 连接远程 MCP 服务器 (DeepWiki) ==========
	// DeepWiki 提供 GitHub 仓库文档搜索功能
	log.Println("正在连接 DeepWiki MCP 服务器...")
	if err := mcpManager.ConnectToRemoteServer(ctx, "deepwiki", "https://mcp.deepwiki.com/sse"); err != nil {
		log.Printf("警告: 无法连接 DeepWiki: %v", err)
	}

	log.Printf("MCP 工具管理器初始化完成，已注册 %d 个 MCP 工具", len(mcpManager.tools))

	// 创建 OpenAI 客户端
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.moonshot.cn/v1"),
	)

	// 定义本地工具 (非 MCP)
	localTools := []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "getCurrentTime",
				Description: param.Opt[string]{Value: "获取当前时间"},
				Parameters:  openai.FunctionParameters{"type": "object"},
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

	// 合并 MCP 工具和本地工具
	allTools := append(localTools, mcpManager.ToOpenAITools()...)

	// 用户输入 - 这里可以改成你想问的问题
	userInput := "帮我查一下 openai/openai-go 这个 GitHub 仓库是做什么的，怎么使用？"
	log.Printf("\n用户输入: %s\n", userInput)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("你是一个智能助手。你可以使用以下工具来帮助用户：\n" +
			"- read_wiki_structure: 获取 GitHub 仓库的文档结构\n" +
			"- read_wiki_contents: 查看 GitHub 仓库的文档内容\n" +
			"- ask_question: 向 GitHub 仓库提问\n" +
			"- getCurrentTime: 获取当前时间\n" +
			"当用户询问关于 GitHub 仓库的问题时，请使用 DeepWiki 的工具来获取信息。"),
		openai.UserMessage(userInput),
	}

	// Agentic Loop
	for {
		log.Printf("发送请求，消息数: %d", len(messages))

		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    "moonshot-v1-8k",
			Messages: messages,
			Tools:    allTools,
		})
		if err != nil {
			panic(err)
		}

		choice := resp.Choices[0]
		log.Printf("收到响应，finish_reason: %s", choice.FinishReason)

		if choice.FinishReason != "tool_calls" || len(choice.Message.ToolCalls) == 0 {
			fmt.Println("\n===== 最终回复 =====")
			fmt.Println(choice.Message.Content)
			break
		}

		for _, toolCall := range choice.Message.ToolCalls {
			log.Printf("调用工具: %s (id=%s)", toolCall.Function.Name, toolCall.ID)

			var result string

			// 检查是否是 MCP 工具
			if _, isMCP := mcpManager.tools[toolCall.Function.Name]; isMCP {
				var args map[string]any
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				result, err = mcpManager.CallTool(ctx, toolCall.Function.Name, args)
				if err != nil {
					result = fmt.Sprintf(`{"error": "%s"}`, err.Error())
				}
			} else {
				// 本地工具
				result = executeLocalTool(toolCall.Function.Name, toolCall.Function.Arguments)
			}

			log.Printf("工具返回: %s", result)

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
