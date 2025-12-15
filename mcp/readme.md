# MCP

## MCP Server

Run MCP server (stdio transport):

```bash
python mcp/server.py
```

## Qwen-Agent 调用示例

可同时接入：

- 本地 stdio MCP Server（本仓库 `mcp/server.py`）
- ModelScope MCP 广场的 MCP Server（SSE URL：`https://mcp.api-inference.modelscope.net/.../sse`）

1. 配置 API Key：

```bash
export DASHSCOPE_API_KEY="sk-xxx"
```

2. 运行：

```bash
python mcp/qwen_agent_client.py
```

### 配置 ModelScope MCP 广场（SSE）

编辑 `mcp/modelscope_mcp_servers.json`，把你在 MCP 广场里“通过 SSE URL 连接服务”拿到的专属 URL 填进去：

```json
{
  "amap-maps": {
    "type": "sse",
    "url": "https://mcp.api-inference.modelscope.net/xxx/sse"
  }
}
```

3. 示例提问：

```text
请调用本地 MCP 工具 env_get，读取环境变量 DASHSCOPE_API_KEY 的值；如果不存在就返回 empty。
```

### 获取“当前位置”（方案 1：IP 定位 + amap-maps 逆地理编码）

本仓库本地 MCP Server 提供 `get_ip_location`（IP 粗略定位，精度有限）。你可以让模型：

1) 先调用本地 `get_ip_location` 得到经纬度

2) 再调用 `amap-maps` 做逆地理编码，得到更可读的地址

示例提问：

```text
请先调用本地 MCP 工具 get_ip_location 获取我的经纬度，然后使用 amap-maps 对该经纬度做逆地理编码，告诉我我现在大概在哪。
```

注：天气查询 `get_weather` 依赖外网访问 open-meteo。
