package queryreqbatch

import (
	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
)

type ExecutorFactory interface {
	Factory(
		ctr *app.Container,
		id int,
		request *ValidatedQueryRequest,
		termChan chan<- TerminateType,
		authToken *auth.AuthToken,
		apiEndpoint string,
		consumer ResponseDataConsumer,
	) (queryreq.QueryExecutor, func(), error)
}

var TypeFactoryMap = map[QueryType]ExecutorFactory{
	FindJob:            FindJobFactory{},
	FindTask:           FindTaskFactory{},
	FindTeam:           FindTeamFactory{},
	FindUserPreference: FindUserPreferenceFactory{},
	FindFileObject:     FindFileObjectFactory{},
	FindOrganization:   FindOrganizationFactory{},
	FindUser:           FindUserFactory{},
	GetTasks:           GetTasksFactory{},
	GetTeams:           GetTeamsFactory{},
	GetFileObjects:     GetFileObjectsFactory{},
	GetOrganizations:   GetOrganizationsFactory{},
	GetUsers:           GetUsersFactory{},
}
