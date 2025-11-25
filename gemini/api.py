from google import genai
from google.genai import types
import base64
import os

# The media_resolution parameter is currently only available in the v1alpha API version.
client = genai.Client(http_options={'api_version': 'v1alpha'})

def load_image_as_base64(image_path):
    """Load an image file and convert it to base64 encoded string"""
    abs_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), image_path)
    with open(abs_path, "rb") as image_file:
        return base64.b64encode(image_file.read())

# Load the image data
image_data = load_image_as_base64("assets/pic1.png")

response = client.models.generate_content(
    model="gemini-3-pro-preview",
    contents=[
        types.Content(
            parts=[
                types.Part(text="What is in this image?"),
                types.Part(
                    inline_data=types.Blob(
                        mime_type="image/png",
                        data=image_data,
                    ),
                    media_resolution={"level": "media_resolution_high"}
                )
            ]
        )
    ]
)

print(response.text)