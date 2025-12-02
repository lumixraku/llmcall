// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/modelcontextprotocol/go-sdk/mcp"
// )

// // 这是一个简单的示例，展示如何连接远程 MCP 服务器 (DeepWiki)
// // DeepWiki 提供 GitHub 仓库文档搜索功能

// func main() {
// 	ctx := context.Background()

// 	log.Println("=== DeepWiki MCP 远程服务演示 ===")
// 	log.Println("连接到 https://mcp.deepwiki.com/sse ...")

// 	// 创建 MCP 客户端
// 	client := mcp.NewClient(&mcp.Implementation{
// 		Name:    "deepwiki-demo",
// 		Version: "v1.0.0",
// 	}, nil)

// 	// 使用 SSE Transport 连接远程 MCP 服务器
// 	transport := &mcp.SSEClientTransport{
// 		Endpoint: "https://mcp.deepwiki.com/sse",
// 	}

// 	// 连接到服务器
// 	session, err := client.Connect(ctx, transport, nil)
// 	if err != nil {
// 		log.Fatalf("连接失败: %v", err)
// 	}
// 	defer session.Close()

// 	log.Println("✅ 连接成功!")

// 	// 列出可用的工具
// 	log.Println("\n=== 获取可用工具列表 ===")
// 	toolsResult, err := session.ListTools(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("获取工具列表失败: %v", err)
// 	}

// 	for _, tool := range toolsResult.Tools {
// 		log.Printf("📦 工具: %s", tool.Name)
// 		log.Printf("   描述: %s", tool.Description)
// 		if tool.InputSchema != nil {
// 			schemaJSON, _ := json.MarshalIndent(tool.InputSchema, "   ", "  ")
// 			log.Printf("   参数: %s", string(schemaJSON))
// 		}
// 		log.Println()
// 	}

// 	// 调用工具示例
// 	log.Println("=== 调用 MCP 工具 ===")

// 	// 1. 获取 openai/openai-go 仓库的文档结构
// 	log.Println("\n📖 调用 read_wiki_structure 获取文档结构...")
// 	result, err := session.CallTool(ctx, &mcp.CallToolParams{
// 		Name: "read_wiki_structure",
// 		Arguments: map[string]any{
// 			"repoName": "openai/openai-go",
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("❌ 调用失败: %v", err)
// 	} else {
// 		log.Println("✅ 调用成功!")
// 		for _, content := range result.Content {
// 			if tc, ok := content.(*mcp.TextContent); ok {
// 				text := tc.Text
// 				if len(text) > 800 {
// 					text = text[:800] + "\n... (更多内容省略)"
// 				}
// 				fmt.Printf("\n%s\n", text)
// 			}
// 		}
// 	}

// 	// 2. 提问关于仓库的问题
// 	log.Println("\n❓ 调用 ask_question 提问...")
// 	result, err = session.CallTool(ctx, &mcp.CallToolParams{
// 		Name: "ask_question",
// 		Arguments: map[string]any{
// 			"repoName": "openai/openai-go",
// 			"question": "How to create a chat completion?",
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("❌ 调用失败: %v", err)
// 	} else {
// 		log.Println("✅ 调用成功!")
// 		for _, content := range result.Content {
// 			if tc, ok := content.(*mcp.TextContent); ok {
// 				text := tc.Text
// 				if len(text) > 1500 {
// 					text = text[:1500] + "\n... (更多内容省略)"
// 				}
// 				fmt.Printf("\n%s\n", text)
// 			}
// 		}
// 	}

// 	log.Println("\n=== 演示结束 ===")
// }
