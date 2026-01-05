package scheduler

import (
	"context"
	"time"
	"log"
	"sync"
	"telemetry-service/internal/config"
	"telemetry-service/internal/k8s"
	"telemetry-service/internal/detectors"
	"telemetry-service/internal/metrics"
	"telemetry-service/internal/events"
)

type Scheduler struct {
	cfg *config.Config
	client *k8s.Client
	collectors []metrics.Collector
	detectors []detectors.Detector
	ctx context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
	publisher *events.Publisher

}


func New(cfg *config.Config, client *k8s.Client) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		cfg: cfg,
		client: client,
		ctx: ctx,
		cancel: cancel,
	}

	s.collectors = []metrics.Collector{
		metrics.NewCPUCollector(client),
	}
	s.detectors = []detectors.Detector{
		detectors.NewCPUSpikeDetector(600),
	}
	s.publisher = events.NewPublisher(cfg)
	return s

}

func (s *Scheduler) Start() {
	log.Println("Starting Scheduler....!")
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping Scheduler....!")
	s.cancel()
	s.wg.Wait()
	log.Println("Scheduler stopped.")
}

func (s *Scheduler) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.cfg.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.executeCycle()
		}
	}
}

func (s *Scheduler) executeCycle() {
	start := time.Now()
	log.Println("Executing cycle at", start)
	// Collect Metrics

	var allMetrics []metrics.Metric
	for _, collector := range s.collectors {
		collected, err := collector.Collect(s.ctx)
		if err != nil {
			log.Printf("Failed to collect metrics from %s: %v", collector.Name(), err)
			continue
		}
		log.Printf("Collected %d metrics from %s", len(collected), collector.Name())
		allMetrics = append(allMetrics, collected...)
	}
	if len(allMetrics) == 0 {
		log.Println("No metrics collected in this cycle")
		return
	}
	AggregatedMetrics := metrics.AggregateMetrics(allMetrics)
	for _, detectors := range s.detectors {
		signal,ok := detectors.Detect(AggregatedMetrics)
		if !ok {
			continue
		}
		err := s.publisher.Publish(s.ctx, signal)
		if err != nil {
			log.Printf("❌ Failed to publish event: %v", err)
			continue
		}

		log.Printf(
		"📤 Event published [%s] severity=%s resource=%s",
		signal.Type,
		signal.Severity,
		signal.Resource,
	)

	}
	elapsed := time.Since(start)
	log.Println("Cycle completed in", elapsed)

}