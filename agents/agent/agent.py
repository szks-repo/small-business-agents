from pydantic import BaseModel
from typing import Dict, Any

from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.pydantic_v1 import BaseModel as V1BaseModel, Field
from langgraph.graph import StateGraph, END

# --- Pydanticモデル定義 ---
class WebhookPayload(BaseModel):
    source: str
    content: str
    metadata: Dict[str, Any]

class ClassificationResponse(BaseModel):
    task_type: str

# --- LangChain/LangGraphでの分類エージェント定義 ---
class TaskClassifier(V1BaseModel):
    """Categorize the user's request into one of the available types."""
    task_type: str = Field(
        description="The type of task. Should be one of 'sales_inquiry', 'customer_support', or 'other'."
    )

def get_classifier_agent():
    llm = ChatOpenAI(model="gpt-4o", temperature=0)
    prompt = ChatPromptTemplate.from_messages([
        ("system", "You are an expert at classifying incoming business inquiries."),
        ("human", "Please classify the following inquiry: \n\n{inquiry}")
    ])
    return prompt | llm.with_structured_output(TaskClassifier)

# --- LangGraphでの専門エージェントの例（セールス問い合わせ） ---
class SalesInquiryState(V1BaseModel):
    original_inquiry: str
    draft_response: str = ""
    final_response: str = ""

def draft_sales_response(state: SalesInquiryState):
    llm = ChatOpenAI(model="gpt-4o", temperature=0.7)
    prompt = ChatPromptTemplate.from_template(
        "Draft a friendly and helpful response to the following sales inquiry:\n\n{inquiry}"
    )
    chain = prompt | llm
    response = chain.invoke({"inquiry": state.original_inquiry})
    return {"draft_response": response.content}

def finalize_response(state: SalesInquiryState):
    # ここでは単純にドラフトを最終結果とするが、レビューや情報追加のステップも可能
    return {"final_response": state.draft_response}

def build_sales_agent_graph():
    graph = StateGraph(SalesInquiryState)
    graph.add_node("draft", draft_sales_response)
    graph.add_node("finalize", finalize_response)
    graph.set_entry_point("draft")
    graph.add_edge("draft", "finalize")
    graph.add_edge("finalize", END)
    return graph.compile()
