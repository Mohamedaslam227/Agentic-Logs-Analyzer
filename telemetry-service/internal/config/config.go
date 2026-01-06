package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Create configs structure

type Config struct {
	ServiceName string
	Environment string
	//Kubernetsis
	ClusterName  string
	PollInterval time.Duration
	//Events / AI Service
	EventSinkURL string
	EventTimeout time.Duration
	//HTTP Server
	//HTTP Server
	HTTPPort string
	//Detectors
	CPUThreshold float64
}

func Load() *Config {
	cfg := &Config{
		ServiceName: getenv("SERVICE_NAME", "telemetry-service"),
		Environment: getenv("ENVIRONMENT", "development"),
		//ClusterName: getenv("CLUSTER_NAME", "local-cluster"),
		PollInterval: getDurationEnv("POLL_INTERVAL", 30*time.Second),
		EventSinkURL: getenv("EVENT_SINK_URL", "http://localhost:8080/events"),
		EventTimeout: getDurationEnv("EVENT_TIMEOUT", 180*time.Second),
		HTTPPort:     getenv("HTTP_PORT", "8080"),
		CPUThreshold: getFloatEnv("CPU_THRESHOLD", 50.0),
	}
	validate(cfg)
	logConfig(cfg)
	return cfg
}

func validate(cfg *Config) {
	if cfg.EventSinkURL == "" {
		log.Fatal("EVENT_SINK_URL is required")
	}
	if cfg.PollInterval <= 0 {
		log.Fatal("POLL_INTERVAL must be greater than 0")
	}
	if cfg.EventTimeout <= 0 {
		log.Fatal("EVENT_TIMEOUT must be greater than 0")
	}

	if cfg.HTTPPort == "" {
		log.Fatal("HTTP_PORT is required")
	}
	if cfg.CPUThreshold <= 0 {
		log.Fatal("CPU_THRESHOLD must be greater than 0")
	}
}

func logConfig(cfg *Config) {
	log.Println("------Configuration Loaded-------")
	log.Println("Service Name:", cfg.ServiceName)
	log.Println("Environment:", cfg.Environment)
	//log.Println("Cluster Name:", cfg.ClusterName)
	log.Println("Poll Interval:", cfg.PollInterval)
	log.Println("Event Sink URL:", cfg.EventSinkURL)
	log.Println("Event Timeout:", cfg.EventTimeout)
	log.Println("HTTP Port:", cfg.HTTPPort)
	log.Println("CPU Threshold:", cfg.CPUThreshold)
}

func getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
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

func getFloatEnv(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Invalid float value for %s: %s", key, value)
	}
	return f
}
