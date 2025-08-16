package outbox

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        int
	EventType string
	Payload   json.RawMessage // this for jsonb storage! i need it i guess
	CreatedAt time.Time
	Processed bool
}

type AdditionEvent struct {
	Number int
}
