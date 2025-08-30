SYSTEM_PROMPT = """
todo system prompt
"""

class InboxAgentPrompt:
    def init(
        self,   
        system_prompt: str = SYSTEM_PROMPT 
    ) -> None:
        self.system_prompt = system_prompt
