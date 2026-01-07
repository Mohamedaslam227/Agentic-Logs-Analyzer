import os
from langchain_core.tools import tool
from kubernetes import client, config
from kubernetes.client.rest import ApiException


_k8s_client_initialized = False

def check_k8s_client():
    """
    Safely initializes the K8s client. 
    Works for both local (minikube/docker-desktop) and in-cluster (Pod) execution.
    """
    global _k8s_client_initialized
    if _k8s_client_initialized:
        return
    try:
        config.load_kube_config()
        print("Loaded cluster configs...")
    except Exception as e:
        print("Failed to load cluster configs: " + str(e))
        raise e
    _k8s_client_initialized = True


@tool
def k8s_get_pod_health(pod_name: str,namespace: str = "default") -> str:
    """
    Advanced Health Check.
    Retrieves the status, restart count, and recent events for a specific Pod.
    Useful for diagnosing 'Pending', 'CrashLoopBackOff', or 'Error' states.
    """
    check_k8s_client()
    v1= client.CoreV1Api()
    try:
        pod = v1.read_namespaced_pod(name=pod_name,namespace=namespace)
        status = pod.status
        container_statuses = []
        if status.container_statuses:
            for c in status.container_statuses:
                state = "Unknown"
                if c.state.running:
                    state = "Running"
                elif c.state.waiting:
                    state = f"Waiting (Reason: {c.state.waiting.reason}) Message: {c.state.waiting.message}"
                elif c.state.terminated:
                    state = f"Terminated (Reason: {c.state.terminated.reason}), ExitCode: {c.state.terminated.exit_code}, Message: {c.state.terminated.message}"
                container_statuses.append(f"- Container '{c.name}': {state} (Restarts: {c.restart_count})")
            
        event_resp = v1.list_namespaced_event(namespace=namespace, field_selector=f"involvedObject.name={pod_name}")
        events = [f"[{e.type}] {e.reason}: {e.message}" for e in event_resp.items]
        
        c_stats_text = "\n".join(container_statuses)
        e_stats_text = "\n".join(events[-5:]) if events else "No recent events found."

        report = f"""
--- Pod Health Report: {pod_name} ---
Phase: {status.phase}
IP: {status.pod_ip}
Node: {status.host_ip}
Containers:
{c_stats_text}
Recent Events:
{e_stats_text}
-------------------------------------
"""
        return report
    except ApiException as e:
        if e.status == 404:
            return f"Pod {pod_name} not found in namespace {namespace}"
        else:
            return f"API Error: {e.reason}"

@tool
def k8s_fetch_logs(pod_name:str,namespace:str="default", lines:int=50) -> str:
    """
    Fetches the last N lines of logs from a Pod.
    Automatically detects if there are multiple containers and fetches logs for the first one.
    """
    

    check_k8s_client()
    v1 = client.CoreV1Api()
    try:
        pod = v1.read_namespaced_pod(name=pod_name,namespace=namespace)
        container_name = pod.spec.containers[0].name
        logs = v1.read_namespaced_pod_log(name=pod_name,namespace=namespace,container=container_name,tail_lines=lines)
        
        if not logs:
            return f"No logs found for pod {pod_name} in namespace {namespace}"
        return logs
    except ApiException as e:
        return f"Failed to fetch logs for pod {pod_name} in namespace {namespace}: {e.reason}"


@tool
def k8s_list_pods(namespace:str="default") -> str:
    """
    Lists all pods in a namespace with their current status.
    Use this to identify which pod is failing if you don't know the exact name.
    """

    check_k8s_client()
    v1 = client.CoreV1Api()
    try:
        pods = v1.list_namespaced_pod(namespace=namespace)
        summary = []
        for p in pods.items:
            restart_count = sum(c.restart_count for c in p.status.container_statuses) if p.status.container_statuses else 0
            summary.append(f"{p.metadata.name} | Status: {p.status.phase} | Restarts: {restart_count}")
        return "\n".join(summary)
    except ApiException as e:
        return f"Failed to list pods in namespace {namespace}: {e.reason}"


k8s_investigator_tools =[k8s_fetch_logs,k8s_list_pods,k8s_get_pod_health]