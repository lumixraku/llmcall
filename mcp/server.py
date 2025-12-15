import os
from datetime import datetime

from mcp.server.fastmcp import FastMCP


app = FastMCP("llmcall")


@app.tool()
def ping() -> str:
    return "pong"


@app.tool()
def echo(text: str) -> str:
    return text


@app.tool()
def get_time(tz: str = "local") -> str:
    if tz != "local":
        raise ValueError("Only tz=local is supported")
    return datetime.now().isoformat(timespec="seconds")


@app.tool()
def env_get(name: str, default: str = "") -> str:
    return os.getenv(name, default)


if __name__ == "__main__":
    app.run(transport="stdio")
