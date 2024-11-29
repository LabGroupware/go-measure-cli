package jobdata

import "github.com/LabGroupware/go-measure-tui/internal/job/jobmodel"

type CreateTaskResultData struct {
	Task jobmodel.TaskWithAttachments `json:"task"`
}

type UpdateStatusTaskResultData struct {
	Task jobmodel.Task `json:"task"`
}
