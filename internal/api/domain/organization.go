package domain

type OrganizationDto struct {
	OrganizationID string                         `json:"organizationId"`
	OwnerID        string                         `json:"ownerId"`
	Name           string                         `json:"name"`
	Plan           string                         `json:"plan"`
	SiteURL        string                         `json:"siteUrl"`
	Users          []UserProfileOnOrganizationDto `json:"users"`
	Owner          UserProfileDto                 `json:"owner"`
	Teams          []TeamDto                      `json:"teams"`
}

type OrganizationOnUserProfileDto struct {
	OrganizationDto `json:",inline"`
}
