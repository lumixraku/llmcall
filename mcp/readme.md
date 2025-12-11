# MCP Client + LLM Agent Demo

## 这个程序做什么？

这是一个 **MCP 客户端 + LLM Agent** 示例程序：

```
用户问题 → LLM (moonshot) → 决定调用工具 → MCP 服务器执行 → 返回结果 → LLM 生成回答
                ↑                                              |
                └──────────────── Agentic Loop ────────────────┘
```

### 核心流程

1. **连接 MCP 服务器** - 连接 DeepWiki (`https://mcp.deepwiki.com/sse`)，获取可用工具
2. **发送用户问题给 LLM** - moonshot 模型决定是否需要调用工具
3. **执行工具调用** - 如果 LLM 返回 `tool_calls`，程序通过 MCP 协议调用工具
4. **循环直到完成** - 将工具结果返回给 LLM，直到 LLM 给出最终回答

### DeepWiki 提供的工具

- `read_wiki_structure` - 获取 GitHub 仓库的文档结构
- `read_wiki_contents` - 查看 GitHub 仓库的文档内容
- `ask_question` - 向 GitHub 仓库提问

## 运行方式

### 1. 设置环境变量

创建 `.env` 文件（在项目根目录 `/Users/mac/repos/llmcall/`）：

```bash
MOONSHOT_API_KEY=your_api_key_here
```

### 2. 运行程序

```bash
cd /Users/mac/repos/llmcall
go run mcp/weather_mcp.go
```

### 3. 修改问题

编辑代码中的 `userInput` 变量来改变问题：

### 4. 运行结果
go run go_mcp.go
🗣️  用户需求: 我想找附近好吃的，最好适合今天的天气

📋 执行计划:
## 执行计划

### 步骤 1: 获取当前天气
- 调用服务: mcp-weather-free
- 输入参数: `get_weather(city)`
- 输出结果: 返回当前城市的天气情况，包括温度、降水概率等。

### 步骤 2: 根据天气选择合适的餐厅类型
- 调用服务: 无（内部逻辑判断）
- 输入参数: 根据步骤1获取的天气情况
- 输出结果: 确定适合当前天气的餐厅类型（例如，晴天可能适合户外烧烤，雨天可能适合室内火锅）。

### 步骤 3: 搜索附近的餐厅
- 调用服务: mcp-server-google-maps
- 输入参数: `maps_search_places(query, location, radius)`
  - query: 根据步骤2确定的餐厅类型
  - location: 用户当前位置（可以通过`maps_reverse_geocode(lat, lng)`获取）
  - radius: 搜索半径，例如5公里
- 输出结果: 返回附近符合条件的餐厅列表。

## 数据流转
1. 用户请求 -> 获取当前天气（mcp-weather-free）-> 确定适合天气的餐厅类型
2. 确定适合天气的餐厅类型 -> 搜索附近的餐厅（mcp-server-google-maps）

## 最终输出
最终返回给用户的结果将是一个附近的餐厅列表，这些餐厅被筛选为适合当前天气的类型。例如，如果今天是晴天，可能会推荐户外烧烤餐厅；如果是雨天，可能会推荐室内火锅餐厅。餐厅列表将包括餐厅名称、地址、距离等信息。
