# AI DevOps SRE Assistant

An intelligent Kubernetes telemetry and incident response system using AI agents.

## Architecture

This project consists of three microservices:

1. **Telemetry Service** (Go) - Monitors Kubernetes cluster metrics
2. **Agent Service** (Python/FastAPI) - AI-powered incident analysis and decision-making
3. **Ollama Service** - Local LLM inference engine

```
┌─────────────────┐      ┌──────────────┐      ┌─────────────┐
│  Telemetry      │─────▶│    Agent     │─────▶│   Ollama    │
│  Service (Go)   │ HTTP │ Service (Py) │ HTTP │  (qwen2.5)  │
└─────────────────┘      └──────────────┘      └─────────────┘
        │                        │
        ▼                        ▼
   Kubernetes              AI Decision
    Metrics                  Making
```

## Quick Start

### Prerequisites
- Docker
- Minikube
- kubectl
- Make
- Python 3.11+ (for local development)
- Go 1.22+ (for local development)

### Deploy to Kubernetes

```bash
# Build all Docker images
make build-all

# Deploy everything to Kubernetes
make deploy-all
```

### Local Development

**Run Telemetry Service:**
```bash
make local-telemetry
```

**Run Agent Service:**
```bash
cd agent-service
python -m venv venv
.\venv\Scripts\activate
pip install -r requirements.txt
make local-agent
```

## Configuration

All services use ConfigMaps for configuration. You can modify the following files:

- `telemetry-service/deployments/configmap.yaml` - Telemetry service config
- `agent-service/deployments/configmap.yaml` - Agent service config
- `ollama/ollama.yaml` - Ollama service config

### Environment Variables

**Telemetry Service:**
- `SERVICE_NAME` - Service name (default: `telemetry-service`)
- `ENVIRONMENT` - Environment (default: `production`)
- `POLL_INTERVAL` - Metrics polling interval in seconds (default: `30`)
- `EVENT_SINK_URL` - Agent service URL (default: `http://agent-service/events`)
- `EVENT_TIMEOUT` - HTTP timeout in seconds (default: `180`)
- `HTTP_PORT` - HTTP server port (default: `8080`)
- `CPU_THRESHOLD` - CPU spike threshold percentage (default: `75.0`)

**Agent Service:**
- `OLLAMA_MODEL` - LLM model to use (default: `qwen2.5:0.5b`)
- `OLLAMA_BASE_URL` - Ollama service URL (default: `http://ollama-service:11434`)
- `OLLAMA_TEMPERATURE` - LLM temperature (default: `0.1`)
- `OLLAMA_NUM_CTX` - Context window size (default: `2048`)
- `OLLAMA_TIMEOUT` - Request timeout in seconds (default: `60`)

## Makefile Targets

### Local Development
- `make local-telemetry` - Run telemetry-service locally
- `make local-agent` - Run agent-service locally

### Docker Build
- `make build-telemetry` - Build telemetry-service image
- `make build-agent` - Build agent-service image
- `make build-all` - Build all images

### Kubernetes Deployment
- `make deploy-all` - Deploy all services
- `make deploy-ollama` - Deploy Ollama only
- `make deploy-telemetry` - Deploy telemetry-service only
- `make deploy-agent` - Deploy agent-service only

### Cleanup
- `make stop-all` - Stop all Kubernetes services
- `make clean` - Remove Docker images

## Testing

Send a test event to the agent service:
```bash
kubectl port-forward svc/agent-service 8080:80

# In another terminal
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d @agent-service/sample_payload_cpu.json
```

## Project Structure

```
.
├── Makefile                      # Main deployment orchestration
├── agent-service/                # AI agent service (Python/FastAPI)
│   ├── Dockerfile
│   ├── deployments/
│   │   ├── configmap.yaml
│   │   └── agent.yaml
│   ├── main.py
│   ├── agent.py
│   ├── nodes.py
│   ├── llm.py
│   └── requirements.txt
├── telemetry-service/            # Metrics collection (Go)
│   ├── Dockerfile
│   ├── deployments/
│   │   ├── configmap.yaml
│   │   └── telemetry.yaml
│   ├── cmd/server/
│   └── internal/
└── ollama/                       # LLM service
    ├── Makefile
    └── ollama.yaml
```
