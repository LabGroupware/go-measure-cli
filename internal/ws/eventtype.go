package ws

type EventType string

const (
	EventTypesJobBegan     EventType = "JOB_BEGAN"
	EventTypesJobProcessed EventType = "JOB_PROCESSEDD"
	EventTypesJobSuccess   EventType = "JOB_SUCCESS"
	EventTypesJobFailed    EventType = "JOB_FAILED"
)

func (et EventType) String() string {
	return string(et)
}

func NewFromEventTypesString(eventTypes []string) []EventType {
	var ets []EventType = make([]EventType, 0, len(eventTypes))
	for _, et := range eventTypes {
		ets = append(ets, EventType(et))
	}
	return ets
}

func ContainsEventType(eventTypes []string, et EventType) bool {
	for _, eventType := range eventTypes {
		if eventType == string(et) {
			return true
		}
	}
	return false
}
