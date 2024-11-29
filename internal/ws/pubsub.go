package ws

type SubscriptionManager struct {
	subscriptions map[EventType]map[string]chan *EventResponseMessage
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[EventType]map[string]chan *EventResponseMessage),
	}
}

func (sm *SubscriptionManager) Subscribe(subscriptionID string, eventType EventType) chan *EventResponseMessage {
	if _, ok := sm.subscriptions[eventType]; !ok {
		sm.subscriptions[eventType] = make(map[string]chan *EventResponseMessage)
	}

	ch := make(chan *EventResponseMessage)
	sm.subscriptions[eventType][subscriptionID] = ch
	return ch
}
