import os

from qwen_agent.agents import Assistant


def main() -> None:
    llm_cfg = {
        "model": "qwen-plus-latest",
        "model_server": "https://dashscope.aliyuncs.com/compatible-mode/v1",
        "api_key": os.getenv("DASHSCOPE_API_KEY"),
    }

    system = "你是一个助手，你可以通过 MCP 工具获取时间、回显文本。"

    tools = [
        {
            "mcpServers": {
                "llmcall": {
                    "command": "python",
                    "args": ["mcp/server.py"],
                }
            }
        }
    ]

    bot = Assistant(
        llm=llm_cfg,
        name="助手",
        description="MCP demo",
        system_message=system,
        function_list=tools,
    )

    messages = []
    while True:
        query = input("\nuser question: ")
        if not query.strip():
            print("user question cannot be empty！")
            continue

        messages.append({"role": "user", "content": query})
        bot_response = ""

        for response_chunk in bot.run(messages):
            new_response = response_chunk[-1]
            if new_response.get("role") == "assistant" and "content" in new_response:
                incremental_content = new_response["content"][len(bot_response) :]
                print(incremental_content, end="", flush=True)
                bot_response += incremental_content

            messages.extend(response_chunk)


if __name__ == "__main__":
    main()
