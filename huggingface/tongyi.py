from diffusers import ZImagePipeline
import torch

# 从本地缓存加载 (CPU + float32，MPS 有 bug 会产生黑图)
pipe = ZImagePipeline.from_pretrained(
    "/Users/mac/.cache/huggingface/hub/models--Tongyi-MAI--Z-Image-Turbo/snapshots/78771b7e11b922c868dd766476bda1f4fc6bfc96/",
    torch_dtype=torch.float32,
    local_files_only=True,
    low_cpu_mem_usage=True,
)
pipe.enable_attention_slicing("max")  # 最大程度减少内存

# 生成图片 (CPU 会慢一些，但结果正确)
with torch.no_grad():
    image = pipe(
        "Tokyo nightview, high quality, detailed, 8k",
        num_inference_steps=25,  # 增加步数提升质量
        guidance_scale=7.5,      # 提示词引导强度
    ).images[0]
image.save("./assets/output/tongyi.png")