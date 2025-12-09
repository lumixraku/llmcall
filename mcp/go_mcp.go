package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// 真实的 MCP 服务列表 (来自 mcp.so 和 github.com/modelcontextprotocol/servers)
const mcpServicesDescription = `
## 可用的 MCP 服务 (真实服务)

### 1. @modelcontextprotocol/server-fetch
- 来源: https://github.com/modelcontextprotocol/servers/tree/main/src/fetch
- 启动: npx -y @modelcontextprotocol/server-fetch
- 功能: 抓取网页内容并转换为 LLM 友好格式
- 工具:
  - fetch(url) - 获取网页内容

### 2. @modelcontextprotocol/server-time
- 来源: https://github.com/modelcontextprotocol/servers/tree/main/src/time
- 启动: npx -y @modelcontextprotocol/server-time
- 功能: 时间和时区转换
- 工具:
  - get_current_time(timezone?) - 获取当前时间
  - convert_time(time, from_tz, to_tz) - 时区转换

### 3. @anthropic/mcp-server-brave-search
- 来源: https://github.com/anthropics/mcp-servers
- 启动: npx -y @anthropic/mcp-server-brave-search (需要 BRAVE_API_KEY)
- 功能: Brave 搜索引擎
- 工具:
  - brave_web_search(query) - 网页搜索
  - brave_local_search(query, location) - 本地商家搜索

### 4. mcp-server-google-maps
- 来源: https://github.com/modelcontextprotocol/servers
- 启动: npx -y @anthropic/mcp-server-google-maps (需要 GOOGLE_MAPS_API_KEY)
- 功能: Google 地图服务
- 工具:
  - maps_geocode(address) - 地址转经纬度
  - maps_reverse_geocode(lat, lng) - 经纬度转地址
  - maps_search_places(query, location, radius) - 搜索附近地点
  - maps_directions(origin, destination) - 获取路线

### 5. mcp-weather-free (Open-Meteo)
- 来源: https://github.com/microagents/mcp-weather-free
- 启动: npx -y mcp-weather-free
- 功能: 免费天气服务 (无需 API Key)
- 工具:
  - get_weather(city) - 获取城市天气
  - get_weather_by_coords(lat, lon) - 通过经纬度获取天气

### 6. @anthropic/mcp-server-filesystem
- 来源: https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem
- 启动: npx -y @modelcontextprotocol/server-filesystem /path/to/allowed/dir
- 功能: 文件系统操作
- 工具:
  - read_file(path) - 读取文件
  - write_file(path, content) - 写入文件
  - list_directory(path) - 列出目录
`

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

	systemPrompt := fmt.Sprintf(`你是一个 MCP 服务编排专家。

%s

用户会提出一个需求，你需要：
1. 分析需求，列出实现步骤
2. 说明每一步调用哪个 MCP 服务
3. 说明服务之间的数据流转关系

请用以下格式输出：

## 执行计划

### 步骤 1: [步骤名称]
- 调用服务: [服务名]
- 输入参数: [参数说明]
- 输出结果: [结果说明]

### 步骤 2: ...
(依此类推)

## 数据流转
[用箭头说明数据如何从一个服务传递到下一个服务]

## 最终输出
[说明最终返回给用户的结果]
`, mcpServicesDescription)

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

	fmt.Println("📋 执行计划:")
	fmt.Println(resp.Choices[0].Message.Content)
}
