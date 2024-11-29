package jobmodel

type OrganizationPlan string

const (
	OrganizationPlanBasic    OrganizationPlan = "BASIC"
	OrganizationPlanStandard OrganizationPlan = "STANDARD"
	OrganizationPlanPremium  OrganizationPlan = "PREMIUM"
)

type Organization struct {
	OrganizationID string           `json:"organizationId"`
	OwnerID        string           `json:"ownerId"`
	Name           string           `json:"name"`
	Plan           OrganizationPlan `json:"plan"`
	SiteURL        string           `json:"siteUrl"`
}

type OrganizationWithUsers struct {
	Organization Organization         `json:"organization"`
	Users        []UserOnOrganization `json:"users"`
}
