.PHONY: help local-telemetry local-agent build-telemetry build-agent deploy-all deploy-telemetry deploy-agent deploy-ollama stop-all clean

# Default target
help:
	@echo "Available targets:"
	@echo "  Local Development:"
	@echo "    local-telemetry    - Run telemetry-service locally"
	@echo "    local-agent        - Run agent-service locally (requires Python venv)"
	@echo ""
	@echo "  Docker Build:"
	@echo "    build-telemetry    - Build telemetry-service Docker image"
	@echo "    build-agent        - Build agent-service Docker image"
	@echo "    build-all          - Build all Docker images"
	@echo ""
	@echo "  Kubernetes Deployment:"
	@echo "    deploy-all         - Deploy all services to Kubernetes"
	@echo "    deploy-ollama      - Deploy Ollama service"
	@echo "    deploy-telemetry   - Deploy telemetry-service"
	@echo "    deploy-agent       - Deploy agent-service"
	@echo ""
	@echo "  Cleanup:"
	@echo "    stop-all           - Stop all Kubernetes services"
	@echo "    clean              - Clean Docker images"

# ========== LOCAL DEVELOPMENT ==========
local-telemetry:
	@echo "Running telemetry-service locally..."
	cd telemetry-service && go run cmd/telemetry/main.go

local-agent:
	@echo "Setting up agent-service..."
	@if not exist "agent-service\venv" ( \
		echo Creating virtual environment... && \
		cd agent-service && python -m venv venv \
	) else ( \
		echo Virtual environment already exists \
	)
	@echo "Installing dependencies..."
	@cd agent-service && .\venv\Scripts\pip.exe install -r requirements.txt
	@echo "Starting agent-service..."
	@echo "Make sure Ollama is running on http://localhost:11434"
	@cd agent-service && .\venv\Scripts\python.exe -m uvicorn main:app --host 0.0.0.0 --port 8080 --reload

# ========== DOCKER BUILD ==========
build-telemetry:
	@echo "Building telemetry-service Docker image..."
	cd telemetry-service && docker build -t telemetry-service:latest .

build-agent:
	@echo "Building agent-service Docker image..."
	cd agent-service && docker build -t agent-service:latest .

build-all: build-telemetry build-agent
	@echo "All Docker images built successfully!"

# ========== KUBERNETES DEPLOYMENT ==========
deploy-all: deploy-ollama deploy-agent deploy-telemetry
	@echo "All services deployed!"
	@echo "Waiting for deployments to be ready..."
	kubectl wait --for=condition=available --timeout=300s deployment/ollama deployment/agent-service deployment/telemetry-service
	@echo "✅ All services are ready!"

deploy-ollama:
	@echo "Checking Minikube status..."
	@minikube status || minikube start --driver=docker
	@echo "Deploying Ollama service..."
	cd ollama && kubectl apply -f ollama.yaml && kubectl apply -f pod.yaml
	@echo "Waiting for Ollama to be ready..."
	kubectl wait --for=condition=ready pod -l app=ollama --timeout=300s
	@echo "Pulling qwen2.5:0.5b model..."
	kubectl exec -it deploy/ollama -- ollama pull qwen2.5:0.5b

deploy-agent:
	@echo "Deploying agent-service..."
	cd agent-service/deployments && kubectl apply -f configmap.yaml && kubectl apply -f agent.yaml

deploy-telemetry:
	@echo "Deploying telemetry-service..."
	cd telemetry-service/deployments && kubectl apply -f configmap.yaml && kubectl apply -f telemetry.yaml

# ========== CLEANUP ==========
stop-all:
	@echo "Stopping all services..."
	-kubectl delete -f telemetry-service/deployments/telemetry.yaml
	-kubectl delete -f telemetry-service/deployments/configmap.yaml
	-kubectl delete -f agent-service/deployments/agent.yaml
	-kubectl delete -f agent-service/deployments/configmap.yaml
	-kubectl delete -f ollama/ollama.yaml
	@echo "All services stopped!"

clean:
	@echo "Removing Docker images..."
	-docker rmi telemetry-service:latest
	-docker rmi agent-service:latest
	@echo "Docker images removed!"

port-forward:
	@echo "Forwarding port 11434..."
	@kubectl port-forward svc/ollama-service 11434:11434