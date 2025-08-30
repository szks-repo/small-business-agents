import operator
from typing import Annotated, Literal, Sequence, TypedDict

from langchain_core.utils.function_calling import convert_to_openai_tool
from langgraph.constants import Send
from langgraph.graph import END, START, StateGraph
from langgraph.pregel import Pregel
from openai import OpenAI
from openai.types.chat import ChatCompletionMessageParam

from models import (
    AgentResult,
    Plan,
    ReflectionResult,
    SearchOutput,
    Subtask,
    ToolResult,
)

from prompts import InboxAgentPrompt
from config import Settings

class AgentState(TypedDict):
    question: str
    plan: list[str]
    current_step: int
    subtask_results: Annotated[Sequence[Subtask], operator.add]
    last_answer: str

class AgentSubGraphState(TypedDict):
    question: str
    plan: list[str]
    subtask: str
    is_completed: bool
    messages: list[ChatCompletionMessageParam]
    challenge_count: int
    tool_results: Annotated[Sequence[Sequence[SearchOutput]], operator.add]
    reflection_results: Annotated[Sequence[ReflectionResult], operator.add]
    subtask_answer: str

class InboxAgent:
    def __init__(
        self,
        settings: Settings,
        tools: list = [],
        prompts: InboxAgentPrompt = InboxAgentPrompt(),
    ) -> None:
        self.settings = settings
        self.tools = tools
        self.tool_map = {tool.name: tool for tool in tools}
        self.prompts = prompts
        self.client = OpenAI(api_key=self.settings.openai_api_key)


    def run_agent(self, question: str = "") -> str:
        return "implement me"