from diffusers import ZImagePipeline
import torch
import os
from utils import read_prompt

# 完全禁用 MPS，强制使用 CPU
os.environ['PYTORCH_ENABLE_MPS_FALLBACK'] = '0'
os.environ['PYTORCH_MPS_HIGH_WATERMARK_RATIO'] = '0'

# 强制使用 CPU
device = "cpu"
dtype = torch.float32
print("强制使用 CPU（避免 MPS 问题）")

pipe = ZImagePipeline.from_pretrained(
    "/Users/mac/.cache/huggingface/hub/models--Tongyi-MAI--Z-Image-Turbo/snapshots/78771b7e11b922c868dd766476bda1f4fc6bfc96/",
    torch_dtype=dtype,
    local_files_only=True,
    low_cpu_mem_usage=True,
)

# 只用 ZImagePipeline 支持的优化方法
pipe.enable_attention_slicing("max")
pipe.enable_model_cpu_offload()
# 删除 pipe.enable_vae_slicing()  # 这个不支持

prompt = read_prompt("weather_card")
print(f"Prompt: {prompt[:100]}...")

print("开始生成图片...")
with torch.no_grad():
    image = pipe(
        prompt,
        height=512,
        width=512,
        num_inference_steps=4,
        guidance_scale=0.0,
        generator=torch.Generator("cpu").manual_seed(42),
    ).images[0]

os.makedirs("./assets/output", exist_ok=True)
image.save("./assets/output/tongyi.png")
print("图片已保存到 ./assets/output/tongyi.png")