package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// MCPServer 表示从官方 MCP Registry 返回的服务信息
type MCPServer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Repository  struct {
		URL    string `json:"url"`
		Source string `json:"source"`
	} `json:"repository"`
	Packages []struct {
		RegistryType string `json:"registryType"`
		Identifier   string `json:"identifier"`
		Version      string `json:"version"`
		RuntimeHint  string `json:"runtimeHint"`
	} `json:"packages"`
}

// MCPRegistryResponse represents the response structure returned by the official MCP Registry API.
type MCPRegistryResponse struct {
	Servers []struct {
		Server MCPServer `json:"server"`
	} `json:"servers"`
	Metadata struct {
		NextCursor string `json:"nextCursor"`
		Count      int    `json:"count"`
	} `json:"metadata"`
}

// searchMCPServers queries the official MCP Registry for servers matching the given search query.
// API docs: https://github.com/modelcontextprotocol/registry
func searchMCPServers(query string, limit int) ([]MCPServer, error) {
	baseURL := "https://registry.modelcontextprotocol.io/v0/servers"
	params := url.Values{}
	params.Set("search", query)
	params.Set("limit", fmt.Sprintf("%d", limit))

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result MCPRegistryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	servers := make([]MCPServer, 0, len(result.Servers))
	for _, s := range result.Servers {
		servers = append(servers, s.Server)
	}
	return servers, nil
}

// formatMCPServersForLLM formats an MCP server list into LLM-friendly markdown text.
func formatMCPServersForLLM(servers []MCPServer) string {
	if len(servers) == 0 {
		return "未找到相关的 MCP 服务"
	}

	var sb strings.Builder
	sb.WriteString("## 搜索到的 MCP 服务\n\n")

	for i, server := range servers {
		sb.WriteString(fmt.Sprintf("### %d. %s\n", i+1, server.Name))
		sb.WriteString(fmt.Sprintf("- 版本: %s\n", server.Version))
		sb.WriteString(fmt.Sprintf("- 描述: %s\n", server.Description))
		if server.Repository.URL != "" {
			sb.WriteString(fmt.Sprintf("- 仓库: %s\n", server.Repository.URL))
		}
		// 显示安装方式
		for _, pkg := range server.Packages {
			if pkg.Identifier != "" {
				runtime := pkg.RuntimeHint
				if runtime == "" {
					if pkg.RegistryType == "npm" {
						runtime = "npx"
					} else if pkg.RegistryType == "pypi" {
						runtime = "uvx"
					}
				}
				sb.WriteString(fmt.Sprintf("- 安装: `%s %s`\n", runtime, pkg.Identifier))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// extractSearchKeywords uses an LLM to extract English MCP Registry search keywords from the user request.
func extractSearchKeywords(ctx context.Context, client openai.Client, userInput string) ([]string, error) {
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "moonshot-v1-8k",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`你是一个关键词提取专家。根据用户的需求，提取出用于搜索 MCP (Model Context Protocol) 服务的英文关键词。

MCP 服务是一些可以被 AI 调用的工具，比如：
- weather: 天气查询
- maps/location: 地图、位置服务
- search: 搜索引擎
- restaurant/food: 餐厅、美食
- time: 时间服务
- fetch: 网页抓取
- filesystem: 文件系统

请直接输出 2-4 个英文关键词，用逗号分隔，不要有其他内容。
例如: weather, restaurant, maps`),
			openai.UserMessage(userInput),
		},
	})
	if err != nil {
		return nil, err
	}

	content := resp.Choices[0].Message.Content
	keywords := strings.Split(content, ",")
	for i := range keywords {
		keywords[i] = strings.TrimSpace(keywords[i])
	}
	return keywords, nil
}

func main() {
	godotenv.Overload()

	apiKey := os.Getenv("MOONSHOT_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 MOONSHOT_API_KEY")
	}

	ctx := context.Background()
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.moonshot.cn/v1"),
	)

	// 用户需求
	userInput := "我想找附近好吃的，最好适合今天的天气"

	fmt.Printf("🗣️  用户需求: %s\n\n", userInput)

	// 第一阶段：提取搜索关键词
	fmt.Println("🔍 正在分析需求，提取关键词...")
	keywords, err := extractSearchKeywords(ctx, client, userInput)
	if err != nil {
		log.Fatalf("提取关键词失败: %v", err)
	}
	fmt.Printf("📝 提取的关键词: %v\n\n", keywords)

	// 第二阶段：搜索 MCP 服务
	fmt.Println("🌐 正在从 modelcontextprotocol.io 搜索 MCP 服务...")
	var allServers []MCPServer
	seenServers := make(map[string]bool)

	for _, keyword := range keywords {
		servers, err := searchMCPServers(keyword, 5)
		if err != nil {
			fmt.Printf("⚠️  搜索 '%s' 失败: %v\n", keyword, err)
			continue
		}
		// 去重
		for _, server := range servers {
			if !seenServers[server.Name] {
				seenServers[server.Name] = true
				allServers = append(allServers, server)
			}
		}
	}

	if len(allServers) == 0 {
		log.Fatal("未找到任何相关的 MCP 服务")
	}

	mcpDescription := formatMCPServersForLLM(allServers)
	fmt.Printf("✅ 找到 %d 个相关 MCP 服务:\n\n", len(allServers))
	fmt.Println(mcpDescription)

	// 第三阶段：生成调用编排计划
	fmt.Println("📋 正在生成 MCP 调用编排计划...\n")

	systemPrompt := fmt.Sprintf(`你是一个 MCP 服务编排专家。

以下是根据用户需求搜索到的可用 MCP 服务：

%s

用户会提出一个需求，你需要：
1. 从上述服务中选择合适的服务来完成任务
2. 分析需求，列出实现步骤
3. 说明每一步调用哪个 MCP 服务
4. 说明服务之间的数据流转关系

请用以下格式输出：

## 选用的 MCP 服务
[列出你选择使用的服务及原因]

## 执行计划

### 步骤 1: [步骤名称]
- 调用服务: [服务包名]
- 输入参数: [参数说明]
- 输出结果: [结果说明]

### 步骤 2: ...
(依此类推)

## 数据流转
[用箭头说明数据如何从一个服务传递到下一个服务]

## 最终输出
[说明最终返回给用户的结果]
`, mcpDescription)

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "moonshot-v1-8k",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userInput),
		},
	})
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("📋 MCP 调用编排计划:")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println(resp.Choices[0].Message.Content)
}
