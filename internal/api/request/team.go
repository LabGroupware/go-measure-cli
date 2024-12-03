package request

type AddUsersTeamRequestDto struct {
	UserIDs []string `json:"userIds"`
}

type CreateTeamRequestDto struct {
	OrganizationID string   `json:"organizationId"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	UserIDs        []string `json:"userIds"`
}
