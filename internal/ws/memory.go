package ws

type SubscribeMessageMemoryData struct {
	AggregateType AggregateType
	AggregateIDs  []string
	EventTypes    []EventType
}

type SubscribeMessageMemory struct {
	Memory map[string]SubscribeMessageMemoryData
}

type UnsubscribeMessageMemoryData struct {
	SubscriptionIDs []string
}

type UnsubscribeMessageMemory struct {
	Memory map[string]UnsubscribeMessageMemoryData
}

type SubscribeMemoryData struct {
	AggregateType AggregateType
	AggregateIDs  []string
	EventTypes    []EventType
}

type SubscribeMemory struct {
	Memory map[string]SubscribeMemoryData
}
