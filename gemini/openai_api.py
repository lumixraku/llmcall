import base64
import os
from openai import OpenAI
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# 使用中转站
sonetto_api_key = os.getenv("SONETTO_API_KEY")
if not sonetto_api_key:
    raise RuntimeError("SONETTO_API_KEY is not set in the environment")

client = OpenAI(
    base_url="https://vip.sonetto.top/v1",
    api_key=sonetto_api_key,
)

def load_image_as_base64(image_path):
    """Load an image file and convert it to base64 encoded string"""
    abs_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), image_path)
    with open(abs_path, "rb") as image_file:
        # OpenAI API 需要字符串格式
        return base64.b64encode(image_file.read()).decode('utf-8')

# Load the image data
image_data = load_image_as_base64("assets/two.png")

response = client.chat.completions.create(
    model="gemini-3-pro-image-preview-16",  # 对方中转站支持的模型名
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "What is in this image?"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/png;base64,{image_data}"
                    }
                }
            ],
        }
    ]
)

message = response.choices[0].message
content = message.content

if isinstance(content, list):
    text = "".join(part.text for part in content if hasattr(part, "text"))
else:
    text = content

print(text)