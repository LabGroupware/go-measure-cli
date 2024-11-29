package job

import (
	"github.com/LabGroupware/go-measure-tui/internal/job/jobaction"
	"github.com/LabGroupware/go-measure-tui/internal/job/jobdata"
)

type UpdateUserPreferenceJobBeganData JobBeganData[jobaction.UpdateUserPreferenceAction]
type UpdateUserPreferenceJobProcessedData JobProcessedData[jobaction.UpdateUserPreferenceAction, jobdata.UpdateUserPreferenceResultData]
type UpdateUserPreferenceJobSuccessData JobSuccessData[jobdata.UpdateUserPreferenceResultData, jobaction.UpdateUserPreferenceAction]
type UpdateUserPreferenceJobFailedData JobFailedData[jobaction.UpdateUserPreferenceAction]
