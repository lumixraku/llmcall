# MCP

## MCP Server

Run MCP server (stdio transport):

```bash
python mcp/server.py
```

## Qwen-Agent 调用示例

1. 配置 API Key：

```bash
export DASHSCOPE_API_KEY="sk-xxx"
```

2. 运行：

```bash
python mcp/qwen_agent_client.py
```

3. 示例提问：

```text
杭州现在天气怎么样？
上海现在气温多少？
现在几点了？
把 hello 原样返回
```

注：天气查询 `get_weather` 依赖外网访问 open-meteo。
