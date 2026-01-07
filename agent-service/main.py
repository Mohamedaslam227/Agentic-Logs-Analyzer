from fastapi import FastAPI
from fastapi.responses import JSONResponse
from model import IncidentSignal
from logger import setup_logger
from agent import build_agent

app = FastAPI()
logger = setup_logger()

@app.get("/")
def root():
    return JSONResponse(status_code=200, content={"message": "Welcome to the Agent Service"})

@app.post("/events")
def receive_event(event: IncidentSignal):
    logger.info(
        f"Event received | type={event.type} "
        f"severity={event.severity} "
        f"resource={event.resource} "
        f"source={event.source}"
    )
    agent = build_agent()
    result = agent.invoke({
        "event_type": event.type,
        "severity": event.severity,
        "resource": event.resource,
        "message": event.message
    })
    
    logger.info(f"Agent Decision: {result.get('decision')}")
    logger.info(f"Root Cause: {result.get('root_cause')}")

    return JSONResponse(status_code=200, content={"message": "Event received", "decision": result.get("decision")})