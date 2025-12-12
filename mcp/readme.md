# MCP 服务发现 + LLM 编排 Demo

## 这个程序做什么？

这是一个 **MCP 服务发现 + LLM 编排** 示例程序，从官方 MCP Registry 动态搜索可用服务：

```
用户需求 → LLM 提取关键词 → 搜索 MCP Registry → LLM 生成调用编排计划
```

### 核心流程

1. **提取搜索关键词** - 使用 LLM 从用户需求中提取英文关键词
2. **搜索 MCP Registry** - 从 `registry.modelcontextprotocol.io` 搜索相关 MCP 服务
3. **格式化服务信息** - 将搜索结果整理为 LLM 可读的格式
4. **生成编排计划** - LLM 根据可用服务生成调用编排方案

### MCP Registry API

程序使用官方 MCP Registry API 搜索服务：

```
GET https://registry.modelcontextprotocol.io/v0/servers?search={query}&limit={limit}
```

返回的服务信息包括：
- **name** - 服务名称
- **description** - 服务描述
- **version** - 版本号
- **repository** - 源码仓库地址
- **packages** - 安装方式（npm/pypi 等）

## 运行方式

### 1. 设置环境变量

创建 `.env` 文件（在项目根目录 `/Users/mac/repos/llmcall/`）：

```bash
MOONSHOT_API_KEY=your_api_key_here
```

### 2. 运行程序

```bash
cd /Users/mac/repos/llmcall/mcp
go run go_mcp.go
```

### 3. 修改问题

编辑代码中的 `userInput` 变量来改变用户需求。

## 运行结果示例

```
🗣️  用户需求: 我想找附近好吃的，最好适合今天的天气

🔍 正在分析需求，提取关键词...
📝 提取的关键词: [weather, restaurant, food, maps]

🌐 正在从 MCP Registry 搜索服务...
✅ 找到 12 个相关 MCP 服务:

## 搜索到的 MCP 服务

### 1. mcp-weather-free
- 版本: 1.0.0
- 描述: 免费天气查询服务
- 安装: `npx @example/mcp-weather-free`

### 2. mcp-server-google-maps
- 版本: 0.6.2
- 描述: Google Maps 位置服务
- 安装: `npx @modelcontextprotocol/server-google-maps`

...

📋 正在生成 MCP 调用编排计划...

=============================================================
📋 MCP 调用编排计划:
=============================================================

## 选用的 MCP 服务
- mcp-weather-free: 获取当前天气
- mcp-server-google-maps: 搜索附近餐厅

## 执行计划

### 步骤 1: 获取当前天气
- 调用服务: mcp-weather-free
- 输入参数: `get_weather(city)`
- 输出结果: 当前城市的天气情况

### 步骤 2: 搜索附近的餐厅
- 调用服务: mcp-server-google-maps
- 输入参数: `maps_search_places(query, location, radius)`
- 输出结果: 附近符合条件的餐厅列表

## 数据流转
用户请求 → 获取天气 → 确定餐厅类型 → 搜索餐厅 → 返回结果
```
