package domain

type TeamDto struct {
	TeamID         string                 `json:"teamId"`
	OrganizationID string                 `json:"organizationId"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	IsDefault      bool                   `json:"isDefault"`
	Users          []UserProfileOnTeamDto `json:"users"`
	Organization   OrganizationDto        `json:"organization"`
	Tasks          []TaskDto              `json:"tasks"`
}

type TeamOnUserProfileDto struct {
	TeamDto `json:",inline"`
}
