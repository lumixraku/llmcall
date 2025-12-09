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

```go
userInput := "帮我查一下 openai/openai-go 这个 GitHub 仓库是做什么的，怎么使用？"
```

## 常见问题

### MCP 协议版本

DeepWiki 支持两种协议：
- `/sse` - 旧版 SSE 协议（2024-11-05），连接不稳定
- `/mcp` - 新版 Streamable HTTP 协议（2025-03-26），推荐使用

代码中使用 `StreamableClientTransport` 连接 `/mcp` 端点。

### "无法访问仓库" 的原因

如果 LLM 回复说"无法访问仓库"，可能是：

1. **MCP 连接失败** - 检查日志是否有 `警告: 无法连接 DeepWiki`
2. **工具未注册** - 日志应显示 `已注册 X 个 MCP 工具`，如果是 0 说明连接失败
3. **LLM 没有调用工具** - moonshot 模型可能没有选择使用工具，而是直接回答
4. **工具调用失败** - DeepWiki 服务器可能返回错误（如 JSON 解析失败）

### 调试方法

观察日志输出：
- `[MCP] 注册工具: xxx` - 确认工具已注册
- `调用工具: xxx` - 确认 LLM 决定调用工具
- `finish_reason: tool_calls` - 表示 LLM 要调用工具
- `finish_reason: stop` - 表示 LLM 直接给出回答（没调用工具）
