package jobdata

import "github.com/LabGroupware/go-measure-tui/internal/job/jobmodel"

type AddUsersOrganizationResultData struct {
	AddedUsers []jobmodel.UserOnOrganization `json:"addedUsers"`
}

type CreateOrganizationResultData struct {
	Organization jobmodel.OrganizationWithUsers `json:"organization"`
	DefaultTeam  jobmodel.Team                  `json:"defaultTeam"`
}
