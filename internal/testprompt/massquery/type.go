package massquery

type QueryType int

const (
	_ QueryType = iota
	FindJob
	FindOrganization
	GetOrganizations
	FindTask
	GetTasks
	FindFileObject
	GetFileObjects
	FindTeam
	GetTeams
	FindUserPreference
	FindUserProfile
	GetUserProfiles
)

type QueryTypeData struct {
	ID          QueryType
	Name        string
	Description string
}

var QueryTypes = []QueryTypeData{
	{ID: FindJob, Name: "Find Job", Description: "Find a job by ID"},
	{ID: FindOrganization, Name: "Find Organization", Description: "Find an organization by ID"},
	{ID: GetOrganizations, Name: "Get Organizations", Description: "Get all organizations"},
	{ID: FindTask, Name: "Find Task", Description: "Find a task by ID"},
	{ID: GetTasks, Name: "Get Tasks", Description: "Get all tasks"},
	{ID: FindFileObject, Name: "Find File Object", Description: "Find a file object by ID"},
	{ID: GetFileObjects, Name: "Get File Objects", Description: "Get all file objects"},
	{ID: FindTeam, Name: "Find Team", Description: "Find a team by ID"},
	{ID: GetTeams, Name: "Get Teams", Description: "Get all teams"},
	{ID: FindUserPreference, Name: "Find User Preference", Description: "Find a user preference by ID"},
	{ID: FindUserProfile, Name: "Find User Profile", Description: "Find a user profile by ID"},
	{ID: GetUserProfiles, Name: "Get User Profiles", Description: "Get all user profiles"},
}
