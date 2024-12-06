package socketsubscribe

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/LabGroupware/go-measure-tui/internal/ws"
	"github.com/jmespath/go-jmespath"
)

type sock struct {
	mu                         *sync.Mutex
	s                          *Socket
	actions                    []SocketSubscribeActionConfig
	actionsFileMap             *sync.Map
	actionIndexToConsumerID    *sync.Map
	consumerIndexToActionID    *sync.Map
	ConsumerSelfEventFilterMap map[string]SelfEventFilter
	consumerTermChanMap        map[string]chan<- string
	Consumers                  []string
}

type SelfEventFilter struct {
	JMESPath *jmespath.JMESPath
}

type globalSock struct {
	mu      *sync.Mutex
	sockMap map[string]*sock
}

var GlobalSock = &globalSock{
	mu:      &sync.Mutex{},
	sockMap: make(map[string]*sock),
}

func NewSock(s *Socket, outputEnabled bool) *sock {
	return &sock{
		mu:                         &sync.Mutex{},
		s:                          s,
		actions:                    []SocketSubscribeActionConfig{},
		actionsFileMap:             &sync.Map{},
		actionIndexToConsumerID:    &sync.Map{},
		consumerIndexToActionID:    &sync.Map{},
		ConsumerSelfEventFilterMap: make(map[string]SelfEventFilter),
		consumerTermChanMap:        make(map[string]chan<- string),
	}
}

func (s *sock) Subscribe(
	ctx context.Context,
	consumerId string,
	aggregateType ws.AggregateType,
	aggregateIDs []string,
	eventTypes []ws.EventType,
	actions []SocketSubscribeActionConfig,
	selfEventFilter SelfEventFilter,
) (<-chan string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.s.Subscribe(ctx, consumerId, aggregateType, aggregateIDs, eventTypes); err != nil {
		return nil, fmt.Errorf("failed to send subscribe message: %v", err)
	}
	notifyChan := make(chan string)

	s.addActions(consumerId, actions)
	s.consumerTermChanMap[consumerId] = notifyChan
	s.ConsumerSelfEventFilterMap[consumerId] = selfEventFilter
	s.Consumers = append(s.Consumers, consumerId)

	return notifyChan, nil
}

func (s *sock) addActions(consumerID string, actions []SocketSubscribeActionConfig) {

	s.actions = append(s.actions, actions...)

	for _, action := range actions {
		s.actionIndexToConsumerID.Store(action.ID, consumerID)
		s.consumerIndexToActionID.Store(consumerID, action.ID)
	}
}

func (s *sock) UnsubscribeNotifyByAction(
	ctx context.Context,
	actionId string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("Retrieved unsubscribe notify by action: lock", actionId)
	defer fmt.Println("Released unsubscribe notify by action: unlock", actionId)

	var consumerId string
	var ok bool
	if consumerId, ok = s.getConsumerIDByActionID(actionId); ok {
		if v, ok := s.consumerTermChanMap[consumerId]; ok {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled")
			case v <- actionId: // Notify
			}
			return nil
		}
	}
	return nil
}

func (s *sock) getConsumerIDByActionID(actionID string) (string, bool) {

	if consumerID, ok := s.actionIndexToConsumerID.Load(actionID); ok {
		return consumerID.(string), true
	}

	return "", false
}

func (s *sock) GetConsumerIDByActionID(actionID string) (string, bool) {

	if consumerID, ok := s.actionIndexToConsumerID.Load(actionID); ok {
		return consumerID.(string), true
	}

	return "", false
}

func (s *sock) Unsubscribe(
	ctx context.Context,
	consumerId string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("Retrieved unsubscribe: lock")
	defer fmt.Println("Released unsubscribe: unlock")

	if err := s.s.ws.UnsubscribeByConsumerID(consumerId); err != nil {
		return fmt.Errorf("failed to send unsubscribe message: %v", err)
	}
	s.removeActions(s.getActionsByConsumerID(consumerId))

	for i, c := range s.Consumers {
		if c == consumerId {
			s.Consumers = append(s.Consumers[:i], s.Consumers[i+1:]...)
		}
	}

	return nil
}

func (s *sock) getActionsByConsumerID(consumerID string) []SocketSubscribeActionConfig {
	actions := []SocketSubscribeActionConfig{}
	if actionID, ok := s.consumerIndexToActionID.Load(consumerID); ok {
		for _, action := range s.actions {
			if action.ID == actionID {
				actions = append(actions, action)
			}
		}
	}

	return actions
}

func (s *sock) removeActions(actions []SocketSubscribeActionConfig) {
	for _, action := range actions {
		for i, a := range s.actions {
			if a.ID == action.ID {
				s.actions = append(s.actions[:i], s.actions[i+1:]...)
			}
		}

		if consumerID, ok := s.actionIndexToConsumerID.Load(action.ID); ok {
			s.actionIndexToConsumerID.Delete(action.ID)
			s.consumerIndexToActionID.Delete(consumerID)
		}
	}
}

func (s *sock) AddActionsFileMap(actionId string, file *os.File) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.actionsFileMap.Store(actionId, file)
}

func (s *sock) GetActionsFileMap(actionId string) (*os.File, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if file, ok := s.actionsFileMap.Load(actionId); ok {
		return file.(*os.File), true
	}

	return nil, false
}

func (s *sock) RemoveActionsFileMap(actionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.actionsFileMap.Delete(actionId)
}

func (s *sock) RemovePluralActionsFileMap(actionIds []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, actionId := range actionIds {
		s.actionsFileMap.Delete(actionId)
	}
}

func (s *sock) removeAllActionsFileMap() {
	s.actionsFileMap = &sync.Map{}
}

func (s *sock) removeAllActions() {
	s.actions = []SocketSubscribeActionConfig{}
	s.actionIndexToConsumerID = &sync.Map{}
	s.consumerIndexToActionID = &sync.Map{}
}

func (s *sock) GetActions() []SocketSubscribeActionConfig {
	return s.actions
}

func (s *globalSock) FindSocket(id string) (*sock, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sock, ok := s.sockMap[id]; ok {
		return sock, nil
	}

	return nil, fmt.Errorf("socket not found: %s", id)
}

func (s *globalSock) AddSocket(id string, sock *sock) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sockMap[id] = sock
}

func (s *globalSock) CloseSocket(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sock, ok := s.sockMap[id]; ok {
		sock.removeAllActions()
		sock.removeAllActionsFileMap()
		sock.s.Close()
		delete(s.sockMap, id)
	}
}
