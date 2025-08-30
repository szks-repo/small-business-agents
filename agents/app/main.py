from pydantic import BaseModel
from fastapi import FastAPI
from typing import Dict, Any, List
import uvicorn

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

# classifier_agent = get_class()
# sales_agent = build_sales_agent_graph()
# @app.post("/classify", response_model=ClassificationResponse)
# async def classify_task(payload: WebhookPayload):
#     """タスクの内容を分類する"""
#     result = classifier_agent.invoke({"inquiry": payload.content})
#     return ClassificationResponse(task_type=result.task_type)

@app.get("/health")
def read_root():
    return {"status": "healthy"}

@app.post("/inbox", response_model=ExecutionResponse)
async def execute_task(payload: InboxRequest):
    print("--- Received full ---")
    print(payload.model_dump_json(indent=2))
    print("---------------------------------")

    return ExecutionResponse(result="implement me.")

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, log_level="debug")
