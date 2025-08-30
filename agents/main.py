from pydantic import BaseModel
from fastapi import FastAPI
from typing import Dict, Any, List
import uvicorn

from agent import InboxAgent
from config import Settings

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
def read_root():
    return {"status": "healthy"}

@app.post("/inbox", response_model=ExecutionResponse)
async def execute_task(payload: InboxRequest):
    print("--- Received full ---")
    print(payload.model_dump_json(indent=2))
    print("---------------------------------")

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
