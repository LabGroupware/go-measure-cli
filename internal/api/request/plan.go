package request

type CreateTaskRequestDto struct {
	TeamID        string   `json:"teamId"`
	ChargeUserID  string   `json:"chargeUserId"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	StartDatetime string   `json:"startDatetime"`
	DueDatetime   string   `json:"dueDatetime"`
	AttachmentIDs []string `json:"attachmentIds"`
}

type UpdateStatusTaskRequestDto struct {
	Status string `json:"status"`
}
