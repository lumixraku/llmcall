from pathlib import Path


def read_prompt(prompt_name: str) -> str:
    """
    读取 prompts 目录下的 prompt 文件内容

    Args:
        prompt_name: prompt 文件名（不含路径）

    Returns:
        文件内容字符串
    """
    prompts_dir = Path(__file__).parent.parent / "prompts"
    prompt_path = prompts_dir / prompt_name

    if not prompt_path.exists():
        raise FileNotFoundError(f"Prompt file not found: {prompt_path}")

    return prompt_path.read_text(encoding="utf-8")
