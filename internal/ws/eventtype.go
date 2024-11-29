package ws

type EventType string

const (
	EventTypesJobBegan     EventType = "JOB_BEGAN"
	EventTypesJobProcessed EventType = "JOB_PROCESSED"
	EventTypesJobSuccess   EventType = "JOB_SUCCESS"
	EventTypesJobFailed    EventType = "JOB_FAILED"
)
