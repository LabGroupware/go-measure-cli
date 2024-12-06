package ws

type SubscribeMessageMemoryData struct {
	ConsumerID    string
	AggregateType AggregateType
	AggregateIDs  []string
	EventTypes    []EventType
}

type SubscribeMessageMemory struct {
	Memory map[string]SubscribeMessageMemoryData
}

func NewSubscribeMessageMemory() *SubscribeMessageMemory {
	return &SubscribeMessageMemory{
		Memory: make(map[string]SubscribeMessageMemoryData),
	}
}

type UnsubscribeMessageMemoryData struct {
	SubscriptionIDs []string
}

type UnsubscribeMessageMemory struct {
	Memory map[string]UnsubscribeMessageMemoryData
}

func NewUnsubscribeMessageMemory() *UnsubscribeMessageMemory {
	return &UnsubscribeMessageMemory{
		Memory: make(map[string]UnsubscribeMessageMemoryData),
	}
}

type SubscribeMemoryData struct {
	AggregateType AggregateType
	AggregateIDs  []string
	EventTypes    []EventType
}

type SubscribeMemory struct {
	Memory map[string]SubscribeMemoryData
}

func NewSubscribeMemory() *SubscribeMemory {
	return &SubscribeMemory{
		Memory: make(map[string]SubscribeMemoryData),
	}
}

type SubscribeConsumerMemory struct {
	Memory map[string]string
}

func NewSubscribeConsumerMemory() *SubscribeConsumerMemory {
	return &SubscribeConsumerMemory{
		Memory: make(map[string]string),
	}
}

type SubscribeConsumerMemorySubIndex struct {
	Memory map[string]string
}

func NewSubscribeConsumerMemorySubIndex() *SubscribeConsumerMemorySubIndex {
	return &SubscribeConsumerMemorySubIndex{
		Memory: make(map[string]string),
	}
}
