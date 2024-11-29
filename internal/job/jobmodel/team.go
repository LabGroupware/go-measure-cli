package jobmodel

type Team struct {
	TeamID         string `json:"teamId"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsDefault      bool   `json:"isDefault"`
}

type TeamWithUsers struct {
	Team  Team         `json:"team"`
	Users []UserOnTeam `json:"users"`
}
