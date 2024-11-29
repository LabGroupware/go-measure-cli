package jobmodel

type TaskStatus string

const (
	TaskStatusPrepare    TaskStatus = "prepare"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancel     TaskStatus = "cancel"
)

type Task struct {
	TaskID        string     `json:"taskId"`
	TeamID        string     `json:"teamId"`
	ChargeUserID  string     `json:"chargeUserId"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Status        TaskStatus `json:"status"`
	StartDateTime string     `json:"startDateTime"`
	DueDateTime   string     `json:"dueDateTime"`
}

type TaskWithAttachments struct {
	Task        Task               `json:"task"`
	Attachments []FileObjectOnTask `json:"attachments"`
}
