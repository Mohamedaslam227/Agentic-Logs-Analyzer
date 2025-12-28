package k8s
import (
	"log"
	"os"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset *kubernetes.Clientset
}

func NewClient() (*Client,error) {
	config, err := buildConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully created k8s client")
	return &Client{Clientset: clientset}, nil


}

func buildConfig() (*rest.Config, error) {
	config,err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}
	kubeconfig := kubeconfigPath()
	log.Println("Using kubeconfig path:", kubeconfig)
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil

	
}

func kubeconfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory", err)
	}
	return filepath.Join(home, ".kube", "config")

	
}