import json
import os
from pathlib import Path

from qwen_agent.agents import Assistant


def _load_modelscope_mcp_servers() -> dict:
    cfg_path = Path(__file__).with_name("modelscope_mcp_servers.json")
    if not cfg_path.exists():
        return {}

    with cfg_path.open("r", encoding="utf-8") as f:
        data = json.load(f)

    if not isinstance(data, dict):
        raise ValueError("modelscope_mcp_servers.json must be a JSON object")

    return data


def main() -> None:
    llm_cfg = {
        "model": "qwen-plus-latest",
        "model_server": "https://dashscope.aliyuncs.com/compatible-mode/v1",
        "api_key": os.getenv("DASHSCOPE_API_KEY"),
    }

    system = "你是一个助手，你可以通过 MCP 工具查询天气、获取时间、回显文本。"

    mcp_servers = {
        "llmcall": {
            "command": "python",
            "args": ["mcp/server.py"],
        }
    }
    mcp_servers.update(_load_modelscope_mcp_servers())

    tools = [{"mcpServers": mcp_servers}]

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
        is_tool_call = False
        tool_call_info = {}

        for response_chunk in bot.run(messages):
            new_response = response_chunk[-1]
            if "function_call" in new_response:
                is_tool_call = True
                tool_call_info = new_response["function_call"]
            elif "function_call" not in new_response and is_tool_call:
                is_tool_call = False
                print("\n" + "=" * 20)
                print("工具调用信息：", tool_call_info)
                print("工具调用结果：", new_response)
                print("=" * 20)
            elif new_response.get("role") == "assistant" and "content" in new_response:
                incremental_content = new_response["content"][len(bot_response) :]
                print(incremental_content, end="", flush=True)
                bot_response += incremental_content

            messages.extend(response_chunk)


if __name__ == "__main__":
    main()
