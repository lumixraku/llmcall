# https://ai.google.dev/gemini-api/docs/image-generation
# 地址api端口：https://vip.sonetto.top/v1
# https://docs.qq.com/sheet/DS1dTa0dZdVpnbnRr?tab=BB08J2

import base64
import os
import re
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

# model name
# https://docs.qq.com/sheet/DS1dTa0dZdVpnbnRr?tab=BB08J2
response = client.chat.completions.create(
    model="gemini-3-pro-image-preview-16",  # 对方中转站支持的模型名
    messages=[
        {
            "role": "user",
            "content": [
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/png;base64,{image_data}"
                    }
                },
                {"type": "text", "text": "Japanese Animate Style"},  # 编辑指令
            ],
        }
    ]
)

message = response.choices[0].message
content = message.content

# 处理返回的图片数据
# 先打印原始响应结构以便调试
# print("Response:", response)
# print("Message:", message)
# print("Content:", content)
# print("Content type:", type(content))

# 尝试不同的响应格式
if hasattr(response, 'generated_images'):
    # Gemini 原生格式
    for idx, img in enumerate(response.generated_images):
        image_bytes = base64.b64decode(img.image.image_bytes)
        output_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), f"assets/edited_{idx}.png")
        with open(output_path, "wb") as f:
            f.write(image_bytes)
        print(f"编辑后的图片已保存到: {output_path}")
elif isinstance(content, list):
    for idx, part in enumerate(content):
        print(f"Part {idx}:", part, type(part))
        # OpenAI 兼容格式可能返回的结构
        if hasattr(part, "image"):
            image_bytes = base64.b64decode(part.image.image_bytes)
            output_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), f"assets/edited_{idx}.png")
            with open(output_path, "wb") as f:
                f.write(image_bytes)
            print(f"编辑后的图片已保存到: {output_path}")
        elif hasattr(part, "text"):
            print("Text:", part.text)
else:
    # 返回的是 markdown 格式: ![image](data:image/png;base64,...)
    pattern = r'!\[.*?\]\(data:image/(\w+);base64,([A-Za-z0-9+/=]+)\)'
    matches = re.findall(pattern, content)
    if matches:
        for idx, (img_format, b64_data) in enumerate(matches):
            image_bytes = base64.b64decode(b64_data)
            output_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), f"assets/output/edited_{idx}.{img_format}")
            with open(output_path, "wb") as f:
                f.write(image_bytes)
            print(f"编辑后的图片已保存到: {output_path}")
    else:
        print("未能解析图片数据:", content[:200])