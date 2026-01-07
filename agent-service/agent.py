from langgraph.graph import StateGraph, END
from langgraph.prebuilt import ToolNode, tools_condition
from model import IncidentState
from nodes import run_investigator, decide_action
from tools import k8s_investigator_tools

def build_agent():
    """
    Constructs the Investigator Agent Graph.
    Flow: User -> Investigator -> (Loop: Tools -> Investigator) -> Decide -> End
    """
    graph = StateGraph(IncidentState)
    
    graph.add_node("investigator", run_investigator)
    graph.add_node("tools", ToolNode(k8s_investigator_tools))
    graph.add_node("decide", decide_action)

    graph.set_entry_point("investigator")
    
    # The 'investigator' node will either yield tool calls or a final text response
    graph.add_conditional_edges(
        "investigator",
        tools_condition,
        {
            "tools": "tools",  # If tool calls detected, go to tools node
            "__end__": "decide"  # If no tools, move to decision making
        }
    )
    
    graph.add_edge("tools", "investigator")
    graph.add_edge("decide", END)
    
    return graph.compile()
