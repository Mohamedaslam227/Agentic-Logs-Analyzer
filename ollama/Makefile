.PHONY: deploy stop pull-qwen run-qwen port-forward status

# Start deployment: Check minikube, deploy pod.yaml, then ollama.yaml
deploy:
	@echo "Checking Minikube status..."
	@minikube status || minikube start --driver=docker
	@echo "Deploying pod.yaml..."
	kubectl apply -f pod.yaml
	@echo "Deploying Ollama..."
	kubectl apply -f ollama.yaml

# Stop deployment: Delete ollama and pod resources
stop:
	@echo "Deleting Ollama resources..."
	-kubectl delete -f ollama.yaml
	@echo "Deleting pod.yaml resources..."
	-kubectl delete -f pod.yaml
	@echo "Stopping Minikube..."
	minikube stop

# Check the status of the rollout
status:
	kubectl rollout status deploy/ollama

# Pull the Qwen model inside the Ollama pod
pull-qwen:
	@echo "Waiting for pod to be ready..."
	kubectl wait --for=condition=ready pod -l app=ollama --timeout=300s
	@echo "Pulling Qwen2.5:0.5b model..."
	kubectl exec -it deploy/ollama -- ollama pull qwen2.5:0.5b

# Run Qwen interactively (CLI)
run-qwen:
	kubectl exec -it deploy/ollama -- ollama run qwen2.5:0.5b

# Port forward to interface with the API locally (localhost:11434)
port-forward:
	@echo "Forwarding port 11434..."
	kubectl port-forward svc/ollama-service 11434:11434
