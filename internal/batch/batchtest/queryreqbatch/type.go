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

type TerminateType int

const (
	_ TerminateType = iota
	ByContext
	ByCount
	BySystemError
	ByCreateRequestError
	ByParseResponseError
	ByWriteError
	ByTimeout
	ByResponseBodyMatch
	ByStatusCodeMatch
)

func NewTerminateTypeFromString(s string) TerminateType {
	switch s {
	case "context":
		return ByContext
	case "count":
		return ByCount
	case "sysError":
		return BySystemError
	case "createRequestError":
		return ByCreateRequestError
	case "parseError":
		return ByParseResponseError
	case "writeError":
		return ByWriteError
	case "time":
		return ByTimeout
	case "responseBody":
		return ByResponseBodyMatch
	case "statusCode":
		return ByStatusCodeMatch
	default:
		return 0
	}
}

func (t TerminateType) String() string {
	switch t {
	case ByContext:
		return "ByContext"
	case ByCount:
		return "ByCount"
	case BySystemError:
		return "BySystemError"
	case ByCreateRequestError:
		return "ByCreateRequestError"
	case ByParseResponseError:
		return "ByParseResponseError"
	case ByWriteError:
		return "ByWriteError"
	case ByTimeout:
		return "ByTimeout"
	case ByResponseBodyMatch:
		return "ByResponseBodyMatch"
	case ByStatusCodeMatch:
		return "ByStatusCodeMatch"
	default:
		return "Unknown"
	}
}
