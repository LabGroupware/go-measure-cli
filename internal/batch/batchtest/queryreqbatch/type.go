package queryreqbatch

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
	FindUser
	GetUsers
)

func NewQueryTypeFromString(s string) QueryType {
	switch s {
	case "FindJob":
		return FindJob
	case "FindOrganization":
		return FindOrganization
	case "GetOrganizations":
		return GetOrganizations
	case "FindTask":
		return FindTask
	case "GetTasks":
		return GetTasks
	case "FindFileObject":
		return FindFileObject
	case "GetFileObjects":
		return GetFileObjects
	case "FindTeam":
		return FindTeam
	case "GetTeams":
		return GetTeams
	case "FindUserPreference":
		return FindUserPreference
	case "FindUser":
		return FindUser
	case "GetUsers":
		return GetUsers
	default:
		return 0
	}
}
