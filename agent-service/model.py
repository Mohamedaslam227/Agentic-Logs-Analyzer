from pydantic import BaseModel,Field
from typing import Optional,Dict,Any,Annotated
from langgraph.graph.message import add_messages
from datetime import datetime

class IncidentSignal(BaseModel):
    id: str = Field(...,description="Unique Event Id")
    type: str = Field(...,description="Event Type",example="cpu_spike")
    severity: str = Field(...,description="Event Severity",example="medium")
    namespace: Optional[str]
    resource: str
    message: str
    timestamp: datetime
    metadata: Optional[Dict[str,str]] = {}
    source: str = Field(...,description="Event Source")


class IncidentState(BaseModel):
    event_type: str
    severity: str
    resource: str
    message: str

    root_cause: Optional[str] = None
    decision: Optional[str] = None
    messages: Annotated[list[Any], add_messages] = []
    