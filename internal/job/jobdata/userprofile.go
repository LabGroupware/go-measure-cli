package jobdata

import "github.com/LabGroupware/go-measure-tui/internal/job/jobmodel"

type CreateUserProfileResultData struct {
	UserProfile    jobmodel.UserProfile    `json:"userProfile"`
	UserPreference jobmodel.UserPreference `json:"userPreference"`
}
