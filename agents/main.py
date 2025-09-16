from pydantic import BaseModel
from fastapi import FastAPI
from typing import Dict, Any, List
import uvicorn

from agent import InboxAgent
from config import Settings
from logger import setup_logger

from ollama import chat

logger = setup_logger(__file__)

class InboxRequest(BaseModel):
    from_address: str
    sender_name: str
    to: List[str]
    cc: List[str]
    subject: str
    body: str
    # Header: Dict[str, List[str]]

class ExecutionResponse(BaseModel):
    result: str

app = FastAPI()

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.get("/ollama")
async def call_ollama():
    response = chat(
        model="gpt-oss:20b",
        messages=[{ 'role': 'user', 'content': 'こんにちは' }],
        think=True,
        options={ "temperature": 0, "num_ctx": 512 }
    )

    return response


@app.post("/inbox", response_model=ExecutionResponse)
async def execute_task(payload: InboxRequest):
    logger.info(f"/inbox: {payload.model_dump_json(indent=2)}")

    agent = InboxAgent(settings=Settings(
        openai_api_key="",
        openai_api_base="",
        openai_model="gpt4.1",
    ))
    result = agent.run_agent(question="""
todo: build question
""")

    return ExecutionResponse(result=result)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, log_level="debug")
