package domain

type UserProfileDto struct {
	UserProfileID      string                         `json:"userProfileId"`
	UserID             string                         `json:"userId"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
	GivenName          string                         `json:"givenName"`
	FamilyName         string                         `json:"familyName"`
	MiddleName         string                         `json:"middleName"`
	Nickname           string                         `json:"nickname"`
	Profile            string                         `json:"profile"`
	Picture            string                         `json:"picture"`
	Website            string                         `json:"website"`
	Phone              string                         `json:"phone"`
	Gender             string                         `json:"gender"`
	Birthdate          string                         `json:"birthdate"`
	Zoneinfo           string                         `json:"zoneinfo"`
	Locale             string                         `json:"locale"`
	UserPreference     UserPreferenceDto              `json:"userPreference"`
	Organizations      []OrganizationOnUserProfileDto `json:"organizations"`
	Teams              []TeamOnUserProfileDto         `json:"teams"`
	OwnedOrganizations []OrganizationDto              `json:"ownedOrganizations"`
	ChargeTasks        []TaskDto                      `json:"chargeTasks"`
}

type UserProfileOnOrganizationDto struct {
	UserProfileDto `json:",inline"`
}

type UserProfileOnTeamDto struct {
	UserProfileDto `json:",inline"`
}
