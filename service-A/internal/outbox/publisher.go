package outbox

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"service-A/internal/kafka"
)

type Publisher struct {
	repository    *Repository
	kafkaProducer *kafka.Producer
	ticker        *time.Ticker
	done          chan bool
	batchSize     int
	workerCount   int
}

type ProcessingResult struct {
	EventID int
	Error   error
}

func NewPublisher(repository *Repository, kafkaProducer *kafka.Producer, intervalSeconds int) *Publisher {
	return &Publisher{
		repository:    repository,
		kafkaProducer: kafkaProducer,
		ticker:        time.NewTicker(time.Duration(intervalSeconds) * time.Second),
		done:          make(chan bool),
		batchSize:     50,
		workerCount:   5,
	}
}

func (p *Publisher) Start(ctx context.Context) {
	log.Println("Starting outbox publisher...")

	go func() {
		for {
			select {
			case <-p.done:
				log.Println("Stopping outbox publisher...")
				return
			case <-p.ticker.C:
				p.processOutboxEvents(ctx)
			}
		}
	}()
}

func (p *Publisher) Stop() {
	p.ticker.Stop()
	p.done <- true
}

func (p *Publisher) processOutboxEvents(ctx context.Context) {
	events, err := p.repository.GetUnprocessedEvents(ctx)
	if err != nil {
		log.Printf("Failed to get unprocessed events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("Processing %d unprocessed events", len(events))

	for i := 0; i < len(events); i += p.batchSize {
		end := i + p.batchSize
		if end > len(events) {
			end = len(events)
		}
		batch := events[i:end]
		p.processBatch(ctx, batch)
	}
}

func (p *Publisher) processBatch(ctx context.Context, events []Event) {
	results := make(chan ProcessingResult, len(events))

	jobs := make(chan Event, len(events))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(jobs, results, &wg)
	}

	for _, event := range events {
		jobs <- event
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	successfulEventIDs := make([]int, 0, len(events))
	failedCount := 0

	for result := range results {
		if result.Error != nil {
			log.Printf("Failed to publish event %d: %v", result.EventID, result.Error)
			failedCount++
		} else {
			successfulEventIDs = append(successfulEventIDs, result.EventID)
		}
	}

	if len(successfulEventIDs) > 0 {
		p.markEventsProcessedBatch(ctx, successfulEventIDs)
	}

	if failedCount > 0 {
		log.Printf("Batch processing completed: %d successful, %d failed",
			len(successfulEventIDs), failedCount)
	} else {
		log.Printf("Batch processing completed successfully: %d events processed",
			len(successfulEventIDs))
	}
}

func (p *Publisher) worker(jobs <-chan Event, results chan<- ProcessingResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for event := range jobs {
		err := p.publishEvent(event)
		results <- ProcessingResult{
			EventID: event.ID,
			Error:   err,
		}
	}
}

func (p *Publisher) markEventsProcessedBatch(ctx context.Context, eventIDs []int) {
	if len(eventIDs) == 0 {
		return
	}

	query := `UPDATE outbox SET processed = TRUE WHERE id = ANY($1)`

	_, err := p.repository.db.Exec(ctx, query, eventIDs)
	if err != nil {
		log.Printf("Failed to mark events as processed in batch: %v", err)

		p.markEventsProcessedIndividually(ctx, eventIDs)
		return
	}

	log.Printf("Successfully marked %d events as processed", len(eventIDs))
}

func (p *Publisher) markEventsProcessedIndividually(ctx context.Context, eventIDs []int) {
	log.Println("Falling back to individual event marking...")

	for _, eventID := range eventIDs {
		if err := p.repository.MarkEventProcessed(ctx, eventID); err != nil {
			log.Printf("Failed to mark event %d as processed: %v", eventID, err)
		}
	}
}

func (p *Publisher) publishEvent(event Event) error {
	switch event.EventType {
	case "addition_result":
		var additionEvent AdditionEvent
		if err := json.Unmarshal(event.Payload, &additionEvent); err != nil {
			return err
		}

		return p.kafkaProducer.SendNumber(additionEvent.Number)
	default:
		log.Printf("Unknown event type: %s", event.EventType)
		return nil
	}
}
