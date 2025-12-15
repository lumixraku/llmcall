import os
import json
from datetime import datetime
from urllib.parse import urlencode
from urllib.request import Request, urlopen

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
def get_weather(city: str) -> str:
    if not city.strip():
        raise ValueError("city cannot be empty")

    geo_params = urlencode({"name": city, "count": 1, "language": "zh"})
    geo_url = f"https://geocoding-api.open-meteo.com/v1/search?{geo_params}"
    geo = _http_get_json(geo_url)
    results = geo.get("results") or []
    if not results:
        return f"未找到城市：{city}"

    r0 = results[0]
    lat = r0.get("latitude")
    lon = r0.get("longitude")
    resolved_name = r0.get("name") or city
    country = r0.get("country") or ""
    admin1 = r0.get("admin1") or ""

    weather_params = urlencode(
        {
            "latitude": lat,
            "longitude": lon,
            "current": "temperature_2m,relative_humidity_2m,wind_speed_10m",
        }
    )
    weather_url = f"https://api.open-meteo.com/v1/forecast?{weather_params}"
    w = _http_get_json(weather_url)
    current = w.get("current") or {}
    temp = current.get("temperature_2m")
    humidity = current.get("relative_humidity_2m")
    wind = current.get("wind_speed_10m")
    t = current.get("time")

    location = ", ".join([x for x in [resolved_name, admin1, country] if x])
    return (
        f"{location} 当前天气（{t}）: "
        f"气温 {temp}°C, 相对湿度 {humidity}%, 风速 {wind} km/h"
    )


if __name__ == "__main__":
    app.run(transport="stdio")
