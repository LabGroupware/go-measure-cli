package job

type Job[T, U any] struct {
	JobID            string         `json:"jobId"`
	Initialized      bool           `json:"initialized"`
	Success          bool           `json:"success"`
	Process          bool           `json:"process"`
	IsValid          bool           `json:"isValid"`
	Data             T              `json:"data"`
	ScheduledActions []U            `json:"scheduledActions"`
	PendingAction    U              `json:"pendingAction"`
	CompletedActions []JobAction[U] `json:"completedActions"`
	Code             string         `json:"code"`
	Caption          string         `json:"caption"`
	ErrorAttributes  any            `json:"errorAttributes"`
	StartedAt        string         `json:"startedAt"`
	ExpiredAt        string         `json:"expiredAt"`
	CompletedAt      string         `json:"completedAt"`
}

type JobAction[U any] struct {
	ActionCode      U      `json:"actionCode"`
	Success         bool   `json:"success"`
	Data            any    `json:"data"`
	Code            string `json:"code"`
	Caption         string `json:"caption"`
	ErrorAttributes any    `json:"errorAttributes"`
	Datetime        string `json:"datetime"`
}

type JobEventType string

const (
	UserProfileCreatedJobBegin     JobEventType = "org.cresplanex.nova.service.userprofile.Event.UserProfile.CreateJob.Begin"
	UserProfileCreatedJobProcessed JobEventType = "org.cresplanex.nova.service.userprofile.Event.UserProfile.CreateJob.Processed"
	UserProfileCreatedJobSuccess   JobEventType = "org.cresplanex.nova.service.userprofile.Event.UserProfile.CreateJob.Success"
	UserProfileCreatedJobFailed    JobEventType = "org.cresplanex.nova.service.userprofile.Event.UserProfile.CreateJob.Failed"

	UserPreferenceUpdatedJobBegin     JobEventType = "org.cresplanex.nova.service.userpreference.Event.UserPreference.UpdateJob.Begin"
	UserPreferenceUpdatedJobProcessed JobEventType = "org.cresplanex.nova.service.userpreference.Event.UserPreference.UpdateJob.Processed"
	UserPreferenceUpdatedJobSuccess   JobEventType = "org.cresplanex.nova.service.userpreference.Event.UserPreference.UpdateJob.Success"
	UserPreferenceUpdatedJobFailed    JobEventType = "org.cresplanex.nova.service.userpreference.Event.UserPreference.UpdateJob.Failed"

	TeamCreatedJobBegin     JobEventType = "org.cresplanex.nova.service.team.Event.Team.CreateJob.Begin"
	TeamCreatedJobProcessed JobEventType = "org.cresplanex.nova.service.team.Event.Team.CreateJob.Processed"
	TeamCreatedJobSuccess   JobEventType = "org.cresplanex.nova.service.team.Event.Team.CreateJob.Success"
	TeamCreatedJobFailed    JobEventType = "org.cresplanex.nova.service.team.Event.Team.CreateJob.Failed"

	TeamUsersAddedJobBegin     JobEventType = "org.cresplanex.nova.service.team.Event.Team.AddUsersJob.Begin"
	TeamUsersAddedJobProcessed JobEventType = "org.cresplanex.nova.service.team.Event.Team.AddUsersJob.Processed"
	TeamUsersAddedJobSuccess   JobEventType = "org.cresplanex.nova.service.team.Event.Team.AddUsersJob.Success"
	TeamUsersAddedJobFailed    JobEventType = "org.cresplanex.nova.service.team.Event.Team.AddUsersJob.Failed"

	OrganizationCreatedJobBegin     JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.CreateJob.Begin"
	OrganizationCreatedJobProcessed JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.CreateJob.Processed"
	OrganizationCreatedJobSuccess   JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.CreateJob.Success"
	OrganizationCreatedJobFailed    JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.CreateJob.Failed"

	OrganizationUsersAddedJobBegin     JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.AddUsersJob.Begin"
	OrganizationUsersAddedJobProcessed JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.AddUsersJob.Processed"
	OrganizationUsersAddedJobSuccess   JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.AddUsersJob.Success"
	OrganizationUsersAddedJobFailed    JobEventType = "org.cresplanex.nova.service.organization.Event.Organization.AddUsersJob.Failed"

	TaskCreatedJobBegin     JobEventType = "org.cresplanex.nova.service.plan.Event.Task.CreateJob.Begin"
	TaskCreatedJobProcessed JobEventType = "org.cresplanex.nova.service.plan.Event.Task.CreateJob.Processed"
	TaskCreatedJobSuccess   JobEventType = "org.cresplanex.nova.service.plan.Event.Task.CreateJob.Success"
	TaskCreatedJobFailed    JobEventType = "org.cresplanex.nova.service.plan.Event.Task.CreateJob.Failed"

	TaskStatusUpdatedJobBegin     JobEventType = "org.cresplanex.nova.service.plan.Event.Task.UpdateStatusJob.Begin"
	TaskStatusUpdatedJobProcessed JobEventType = "org.cresplanex.nova.service.plan.Event.Task.UpdateStatusJob.Processed"
	TaskStatusUpdatedJobSuccess   JobEventType = "org.cresplanex.nova.service.plan.Event.Task.UpdateStatusJob.Success"
	TaskStatusUpdatedJobFailed    JobEventType = "org.cresplanex.nova.service.plan.Event.Task.UpdateStatusJob.Failed"

	FileObjectCreatedJobBegin     JobEventType = "org.cresplanex.nova.service.storage.Event.FileObject.CreateJob.Begin"
	FileObjectCreatedJobProcessed JobEventType = "org.cresplanex.nova.service.storage.Event.FileObject.CreateJob.Processed"
	FileObjectCreatedJobSuccess   JobEventType = "org.cresplanex.nova.service.storage.Event.FileObject.CreateJob.Success"
	FileObjectCreatedJobFailed    JobEventType = "org.cresplanex.nova.service.storage.Event.FileObject.CreateJob.Failed"
)

type JobRawData struct {
	JobID        string       `json:"jobId"`
	JobEventType JobEventType `json:"jobEventType"`
}

type JobBeganData[T any] struct {
	JobID            string       `json:"jobId"`
	JobEventType     JobEventType `json:"jobEventType"`
	PendingAction    T            `json:"pendingAction"`
	ScheduledActions []T          `json:"scheduledActions"`
	Timestamp        string       `json:"timestamp"`
}

type JobProcessedData[T, U any] struct {
	JobID            string         `json:"jobId"`
	JobEventType     JobEventType   `json:"jobEventType"`
	PendingAction    T              `json:"pendingAction"`
	CompletedActions []JobAction[U] `json:"completedActions"`
	ScheduledActions []T            `json:"scheduledActions"`
	Timestamp        string         `json:"timestamp"`
}

type JobSuccessData[T, U any] struct {
	JobID            string         `json:"jobId"`
	JobEventType     JobEventType   `json:"jobEventType"`
	CompletedActions []JobAction[U] `json:"completedActions"`
	Data             T              `json:"data"`
	Timestamp        string         `json:"timestamp"`
}

type JobFailedData[U any] struct {
	JobID            string         `json:"jobId"`
	JobEventType     string         `json:"jobEventType"`
	CompletedActions []JobAction[U] `json:"completedActions"`
	ErrorAttributes  any            `json:"errorAttributes"`
	Timestamp        string         `json:"timestamp"`
}
