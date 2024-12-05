package socketsubscribe

type SocketSubscribeConfig struct {
	Type        string                             `json:"type"`
	Output      BatchTestOutput                    `json:"output"`
	Subscribes  []SocketSubscribeSubscribeConfig   `json:"subscribes"`
	Actions     []SocketSubscribeActionConfig      `json:"actions"`
	SuccessTerm []string                           `json:"successTerm"`
	Term        SocketSubscribeTermConditionConfig `json:"termCondition"`
}

type BatchTestOutput struct {
	Enabled bool `yaml:"enabled"`
}

type SocketSubscribeSubscribeConfig struct {
	AggregateType string   `json:"aggregateType"`
	AggregateId   []string `json:"aggregateId"`
	EventTypes    []string `json:"eventTypes"`
}

type SocketSubscribeActionConfig struct {
	ID         string                            `json:"id"`
	Types      []string                          `json:"types"`
	EventTypes []string                          `json:"eventTypes"`
	Data       []SocketSubscribeActionDataConfig `json:"data"`
}

type SocketSubscribeActionDataConfig struct {
	Key      string `json:"key"`
	JMESPath string `json:"jmesPath"`
	OnNil    string `json:"onNil"`
	OnError  string `json:"onError"`
}

type SocketSubscribeTermConditionConfig struct {
	Time  *string                                `json:"time"`
	Error []string                               `json:"error"`
	Event []string                               `json:"event"`
	Data  SocketSubscribeTermConditionDataConfig `json:"data"`
}

type SocketSubscribeTermConditionDataConfig struct {
	JMESPath *string `json:"jmesPath"`
}

type ErrorTypeForTerm string

const (
	ErrorTypeForTermParseError     ErrorTypeForTerm = "parse_error"
	ErrorTypeForTermUnmarshalError ErrorTypeForTerm = "unmarshal_error"
	ErrorTypeForTermReadError      ErrorTypeForTerm = "read_error"
	ErrorTypeForTermSendError      ErrorTypeForTerm = "send_error"
)

func ContainsTermError(conditions []string, termError ErrorTypeForTerm) bool {
	for _, condition := range conditions {
		if condition == string(termError) {
			return true
		}
	}
	return false
}

type SuccessTerm string

const (
	SuccessTermClose SuccessTerm = "close"
	SuccessTermTime  SuccessTerm = "time"
	SuccessTermError SuccessTerm = "error"
	SuccessTermEvent SuccessTerm = "event"
	SuccessTermData  SuccessTerm = "data"
)

func ContainsSuccessTerm(terms []string, term SuccessTerm) bool {
	for _, t := range terms {
		if t == string(term) {
			return true
		}
	}
	return false
}

type SocketActionType string

const (
	SocketActionTypeStore  SocketActionType = "store"
	SocketActionTypeOutput SocketActionType = "output"
)

func ContainsSocketActionType(types []string, actionType ...SocketActionType) bool {
	for _, t := range types {
		for _, at := range actionType {
			if t == string(at) {
				return true
			}
		}
	}
	return false
}
