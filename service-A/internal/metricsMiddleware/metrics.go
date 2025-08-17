package metricsMiddleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// gRPC metrics
	GRPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Kafka metrics
	KafkaMessagesProduced = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total number of messages produced to Kafka",
		},
		[]string{"topic", "status"},
	)

	// Outbox metrics
	OutboxEventsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_events_created_total",
			Help: "Total number of outbox events created",
		},
	)

	OutboxEventsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outbox_events_processed_total",
			Help: "Total number of outbox events processed",
		},
		[]string{"status"},
	)

	OutboxEventsInQueue = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "outbox_events_in_queue",
			Help: "Number of unprocessed events in outbox",
		},
	)

	// Database metrics
	DatabaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"},
	)
)
