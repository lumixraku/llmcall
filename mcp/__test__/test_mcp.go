package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("连接 DeepWiki MCP 服务器...")

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "v1.0.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: "https://mcp.deepwiki.com/mcp",
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer session.Close()

	log.Println("连接成功，列出工具...")

	// 列出工具
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("列出工具失败: %v", err)
	}

	for _, tool := range tools.Tools {
		schema, _ := json.MarshalIndent(tool.InputSchema, "", "  ")
		fmt.Printf("\n工具: %s\n描述: %s\n参数: %s\n", tool.Name, tool.Description, string(schema))
	}

	// 测试调用 read_wiki_structure (这个之前成功过)
	log.Println("\n测试调用 read_wiki_structure...")
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "read_wiki_structure",
		Arguments: map[string]any{
			"repoName": "openai/openai-go",
		},
	})
	if err != nil {
		log.Printf("调用失败: %v", err)
	} else {
		for _, c := range result.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				fmt.Printf("结果: %s\n", tc.Text[:min(500, len(tc.Text))])
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
