package queryreqbatch

import (
	"os"

	"github.com/LabGroupware/go-measure-tui/internal/api/request/queryreq"
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
)

type ExecutorFactory interface {
	Factory(
		ctr *app.Container,
		id int,
		request *ValidatedQueryRequest,
		termChan chan<- struct{},
		authToken *auth.AuthToken,
		apiEndpoint string,
		outputFile *os.File,
	) (queryreq.QueryExecutor, error)
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
