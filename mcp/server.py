import os
import json
from pathlib import Path
from datetime import datetime
from urllib.parse import urlencode
from urllib.request import Request, urlopen

from dotenv import load_dotenv
from mcp.server.fastmcp import FastMCP


_dotenv_path = Path(__file__).resolve().parents[1] / ".env"
load_dotenv(dotenv_path=_dotenv_path, override=False)


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


def _http_get_json(url: str, timeout_seconds: int = 20) -> dict:
    req = Request(
        url,
        headers={
            "User-Agent": "llmcall-mcp/1.0",
            "Accept": "application/json",
        },
        method="GET",
    )
    with urlopen(req, timeout=timeout_seconds) as resp:
        charset = resp.headers.get_content_charset() or "utf-8"
        return json.loads(resp.read().decode(charset))


@app.tool()
def get_ip_location() -> dict:
    """Best-effort approximate location based on public IP (no API key).

    Returns latitude/longitude and coarse address fields if available.
    """

    data = _http_get_json("https://ipapi.co/json/")
    return {
        "ip": data.get("ip"),
        "latitude": data.get("latitude"),
        "longitude": data.get("longitude"),
        "city": data.get("city"),
        "region": data.get("region"),
        "country": data.get("country_name") or data.get("country"),
        "raw": data,
    }

if __name__ == "__main__":
    app.run(transport="stdio")
