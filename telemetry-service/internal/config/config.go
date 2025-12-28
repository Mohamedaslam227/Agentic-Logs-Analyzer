package config
import(
	"os"
	"log"
	"strconv"
	"time"
)

// Create configs structure

type Config struct {
	ServiceName string
	Environment string
	//Kubernetsis
	ClusterName string
	PollInterval time.Duration
	//Events / AI Service
	EventSinkURL string
	EventTimeout time.Duration
	//HTTP Server
	HTTPPort string
}

func Load() *Config {
	cfg := &Config{
		ServiceName: getenv("SERVICE_NAME", "telemetry-service"),
		Environment: getenv("ENVIRONMENT", "development"),
		//ClusterName: getenv("CLUSTER_NAME", "local-cluster"),
		PollInterval: getDurationEnv("POLL_INTERVAL", 30*time.Second),
		EventSinkURL: getenv("EVENT_SINK_URL", "http://localhost:8080/events"),
		EventTimeout: getDurationEnv("EVENT_TIMEOUT", 10*time.Second),
		HTTPPort: getenv("HTTP_PORT", "8080"),

		}
		validate(cfg)
		logConfig(cfg)
		return cfg
}

func validate(cfg *Config) {
	if cfg.EventSinkURL == "" {
		log.Fatal("EVENT_SINK_URL is required")
	}
	if cfg.PollInterval <=0 {
		log.Fatal("POLL_INTERVAL must be greater than 0")
	}
	if cfg.EventTimeout <=0 {
		log.Fatal("EVENT_TIMEOUT must be greater than 0")
	}

	if cfg.HTTPPort == "" {
		log.Fatal("HTTP_PORT is required")
	}
}


func logConfig(cfg *Config) {
	log.Println("------Configuration Loaded-------")
	log.Println("Service Name:", cfg.ServiceName)
	log.Println("Environment:", cfg.Environment)
	log.Println("Cluster Name:", cfg.ClusterName)
	log.Println("Poll Interval:", cfg.PollInterval)
	log.Println("Event Sink URL:", cfg.EventSinkURL)
	log.Println("Event Timeout:", cfg.EventTimeout)
	log.Println("HTTP Port:", cfg.HTTPPort)
}

func getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}


func getDurationEnv(key string,defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	seconds, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid duration value for %s: %s", key, value)
	}
	return time.Duration(seconds) * time.Second
}