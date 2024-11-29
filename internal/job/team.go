package job

import (
	"github.com/LabGroupware/go-measure-tui/internal/job/jobaction"
	"github.com/LabGroupware/go-measure-tui/internal/job/jobdata"
)

type AddUsersTeamJobBeganData JobBeganData[jobaction.AddUsersTeamAction]
type AddUsersTeamJobProcessedData JobProcessedData[jobaction.AddUsersTeamAction, jobdata.AddUsersTeamResultData]
type AddUsersTeamJobSuccessData JobSuccessData[jobdata.AddUsersTeamResultData, jobaction.AddUsersTeamAction]
type AddUsersTeamJobFailedData JobFailedData[jobaction.AddUsersTeamAction]

type CreateTeamJobBeganData JobBeganData[jobaction.CreateTeamAction]
type CreateTeamJobProcessedData JobProcessedData[jobaction.CreateTeamAction, jobdata.CreateTeamResultData]
type CreateTeamJobSuccessData JobSuccessData[jobdata.CreateTeamResultData, jobaction.CreateTeamAction]
type CreateTeamJobFailedData JobFailedData[jobaction.CreateTeamAction]
