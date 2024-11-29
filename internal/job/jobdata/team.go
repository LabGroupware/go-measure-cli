package jobdata

import "github.com/LabGroupware/go-measure-tui/internal/job/jobmodel"

type CreateTeamResultData struct {
	Team jobmodel.TeamWithUsers `json:"team"`
}

type AddUsersTeamResultData struct {
	AddedUsers []jobmodel.UserOnTeam `json:"addedUsers"`
}
