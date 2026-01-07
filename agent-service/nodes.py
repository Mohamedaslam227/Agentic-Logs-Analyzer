from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_core.messages import SystemMessage, HumanMessage
from llm import get_llm
from tools import k8s_investigator_tools

# --- Investigator Node ---
def run_investigator(state):
    """
    The brain of the agent. It receives the state, checks if there are message history,
    and calls the LLM with tools bound.
    """
    llm = get_llm().bind_tools(k8s_investigator_tools)
    
    # If this is the first turn, initialize the conversation with context
    if not state.messages:
        sys_msg = SystemMessage(content="""You are a Senior SRE Agent intentionally designed to investigate Kubernetes incidents.

YOUR PROTOCOL:
1. REVIEW the incident details.
2. THOUGHT: Create a plan. What do I need to check? (e.g., "I need to see the logs to understand the crash.")
3. ACTION: Call the appropriate tool (e.g., `k8s_fetch_logs`).
4. OBSERVATION: The tool will return data. READ IT CAREFULLY.
5. ANALYSIS: Based on the tool output, determine the root cause.
   - If the output is "Connection Refused", the root cause is likely the dependency.
   - If the output is "OOMKilled", the root cause is Memory Limit.
6. FINAL ANSWER: Once you are confident, output a concise explanation of the Root Cause.

DO NOT stop after calling a tool. You MUST provide the final analysis based on the tool's result.
""")
        human_msg = HumanMessage(content=f"""
NEW INCIDENT DETECTED:
- Type: {state.event_type}
- Severity: {state.severity}
- Resource: {state.resource}
- Message: {state.message}

Please investigate.
""")
        # We start with these messages
        # Note: In LangGraph, we return the NEW messages to append.
        # But for the first call, we need to pass them to invoke() as well.
        messages = [sys_msg, human_msg]
        response = llm.invoke(messages)
        # Return the initial prompt AND the LLM's first response (which might be a tool call)
        return {"messages": [sys_msg, human_msg, response]}
    else:
        # Subsequent turns: just pass the history
        response = llm.invoke(state.messages)
        return {"messages": [response]}

# --- Decision Node ---
decision_prompt = ChatPromptTemplate.from_template("""
You are an SRE decision system.
Based on the investigation, choose the best action.

Incident:
- Type: {event_type}
- Severity: {severity}
- Root Cause: {root_cause}

Options:
- auto_mitigate (Only if likely safe and root cause is clear)
- require_human_approval (If dangerous or uncertain)

Answer with ONLY the option name.
""")

def decide_action(state):
    """
    Decides the next step after investigation is complete.
    """
    # The last message from the investigator loop should contain the explanation
    # We search backwards for the last AIMessage that has content (not just tool calls)
    root_cause_text = "Unknown"
    for msg in reversed(state.messages):
        if hasattr(msg, "content") and msg.content:
            root_cause_text = msg.content
            break
            
    chain = decision_prompt | get_llm() | StrOutputParser()
    action = chain.invoke({
        "event_type": state.event_type,
        "severity": state.severity,
        "root_cause": root_cause_text
    })
    return {"decision": action.strip().lower(), "root_cause": root_cause_text}