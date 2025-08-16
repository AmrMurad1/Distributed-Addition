package outbox

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// need this for batch operations
func (r *Repository) GetDB() *pgxpool.Pool {
	return r.db
}

func (r *Repository) SaveEvent(ctx context.Context, eventType string, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `INSERT INTO outbox (event_type, payload) VALUES ($1, $2)`

	_, err = r.db.Exec(ctx, query, eventType, payloadJSON)
	if err != nil {
		log.Printf("Failed to save event to outbox: %v", err)
		return err
	}

	log.Printf("Successfully saved event to outbox: %s", eventType)
	return nil
}

func (r *Repository) GetUnprocessedEvents(ctx context.Context) ([]Event, error) {
	query := `
		SELECT id, event_type, payload, created_at, processed 
		FROM outbox 
		WHERE processed = FALSE 
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.EventType, &event.Payload, &event.CreatedAt, &event.Processed)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *Repository) MarkEventProcessed(ctx context.Context, eventID int) error {
	query := `UPDATE outbox SET processed = TRUE WHERE id = $1`
	_, err := r.db.Exec(ctx, query, eventID)
	if err != nil {
		log.Printf("Failed to mark event %d as processed: %v", eventID, err)
		return err
	}

	log.Printf("Successfully marked event %d as processed", eventID)
	return nil
}
