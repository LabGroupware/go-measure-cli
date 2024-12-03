package response

type JobResponseDto struct {
	JobId            string         `json:"jobId"`
	Initialized      bool           `json:"initialized"`
	Success          bool           `json:"success"`
	Process          bool           `json:"process"`
	IsValid          bool           `json:"isValid"`
	Data             any            `json:"data"`
	ScheduledActions []string       `json:"scheduledActions"`
	PendingAction    string         `json:"pendingAction"`
	CompletedActions []JobActionDto `json:"completedActions"`
	Code             string         `json:"code"`
	Caption          string         `json:"caption"`
	ErrorAttributes  any            `json:"errorAttributes"`
	StartedAt        string         `json:"startedAt"`
	ExpiredAt        string         `json:"expiredAt"`
	CompletedAt      string         `json:"completedAt"`
}

type JobActionDto struct {
	ActionCode      string `json:"actionCode"`
	Success         bool   `json:"success"`
	Data            any    `json:"data"`
	Code            string `json:"code"`
	Caption         string `json:"caption"`
	ErrorAttributes any    `json:"errorAttributes"`
	Datetime        string `json:"datetime"`
}
