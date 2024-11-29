package job

import (
	"github.com/LabGroupware/go-measure-tui/internal/job/jobaction"
	"github.com/LabGroupware/go-measure-tui/internal/job/jobdata"
)

type CreateUserProfileJobBeganData JobBeganData[jobaction.CreateUserProfileAction]
type CreateUserProfileJobProcessedData JobProcessedData[jobaction.CreateUserProfileAction, jobdata.CreateUserProfileResultData]
type CreateUserProfileJobSuccessData JobSuccessData[jobdata.CreateUserProfileResultData, jobaction.CreateUserProfileAction]
type CreateUserProfileJobFailedData JobFailedData[jobaction.CreateUserProfileAction]
