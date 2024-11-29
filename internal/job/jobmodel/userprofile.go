package jobmodel

type UserProfile struct {
	UserProfileID string `json:"userProfileId"`
	UserID        string `json:"userId"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	MiddleName    string `json:"middleName"`
	Nickname      string `json:"nickname"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Website       string `json:"website"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	Birthdate     string `json:"birthdate"`
	Zoneinfo      string `json:"zoneinfo"`
	Locale        string `json:"locale"`
}

type UserOnOrganization struct {
	UserOrganizationID string `json:"userOrganizationId"`
	UserID             string `json:"userId"`
}

type UserOnTeam struct {
	UserTeamID string `json:"userTeamId"`
	UserID     string `json:"userId"`
}
