package request

type AddUsersOrganizationRequestDto struct {
	UserIDs []string `json:"userIds"`
}

type CreateOrganizationRequestDto struct {
	Name    string   `json:"name"`
	Plan    string   `json:"plan"`
	UserIDs []string `json:"userIds"`
}
