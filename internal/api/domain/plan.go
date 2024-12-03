package domain

type TaskDto struct {
	TaskID        string                `json:"taskId"`
	TeamID        string                `json:"teamId"`
	ChargeUserID  string                `json:"chargeUserId"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	Status        string                `json:"status"`
	StartDatetime string                `json:"startDatetime"`
	DueDatetime   string                `json:"dueDatetime"`
	Attachments   []FileObjectOnTaskDto `json:"attachments"`
	Team          TeamDto               `json:"team"`
	ChargeUser    UserProfileDto        `json:"chargeUser"`
}

type TaskOnFileObjectDto struct {
	TaskDto `json:",inline"`
}
