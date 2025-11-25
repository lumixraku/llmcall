from setuptools import setup, find_packages

setup(
    name="api_test",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "google-generativeai>=0.5.0",
        "pillow>=10.0.0",
        "python-dotenv>=1.0.0",
    ],
    python_requires=">=3.9,<3.12",
)
