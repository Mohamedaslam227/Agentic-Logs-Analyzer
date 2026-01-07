import os
from langchain_core.language_models import BaseLanguageModel
from langchain_ollama import ChatOllama

def get_llm() -> BaseLanguageModel:
    llm = ChatOllama(
        model=os.getenv("OLLAMA_MODEL", "qwen2.5:0.5b"),
        base_url=os.getenv("OLLAMA_BASE_URL", "http://localhost:11434"),
        temperature=float(os.getenv("OLLAMA_TEMPERATURE", "0.1")),
        num_ctx=int(os.getenv("OLLAMA_NUM_CTX", "2048")),
        timeout=int(os.getenv("OLLAMA_TIMEOUT", "60"))
    )
    return llm
