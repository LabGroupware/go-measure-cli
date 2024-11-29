package job

import (
	"github.com/LabGroupware/go-measure-tui/internal/job/jobaction"
	"github.com/LabGroupware/go-measure-tui/internal/job/jobdata"
)

type CreateTaskJobBeganData JobBeganData[jobaction.CreateTaskAction]
type CreateTaskJobProcessedData JobProcessedData[jobaction.CreateTaskAction, jobdata.CreateTaskResultData]
type CreateTaskJobSuccessData JobSuccessData[jobdata.CreateTaskResultData, jobaction.CreateTaskAction]
type CreateTaskJobFailedData JobFailedData[jobaction.CreateTaskAction]

type UpdateStatusTaskJobBeganData JobBeganData[jobaction.UpdateStatusTaskAction]
type UpdateStatusTaskJobProcessedData JobProcessedData[jobaction.UpdateStatusTaskAction, jobdata.UpdateStatusTaskResultData]
type UpdateStatusTaskJobSuccessData JobSuccessData[jobdata.UpdateStatusTaskResultData, jobaction.UpdateStatusTaskAction]
type UpdateStatusTaskJobFailedData JobFailedData[jobaction.UpdateStatusTaskAction]
