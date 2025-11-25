# Google Gemini API Test

This project demonstrates how to use Google Gemini 3 API to generate content.

## Setup

1. Create a virtual environment using uv with Python 3.11:
```bash
uv venv --python=3.11
```

2. Activate the virtual environment:
```bash
source .venv/bin/activate
```

3. Install dependencies using uv:
```bash
uv pip install -e .
```

This will install the project and all its dependencies based on pyproject.toml.

4. Set up your Google API key:
```bash
export GOOGLE_API_KEY="your-api-key-here"
```

## Running the Program

```bash
python gemini/api.py
```

This will generate an infographic of the current weather in Tokyo and save it as `weather_tokyo.png`.
