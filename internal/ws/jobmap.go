package ws

import (
	"reflect"

	"github.com/LabGroupware/go-measure-tui/internal/job"
)

type JobDataTypeMapper map[job.JobEventType]reflect.Type

func NewJobDataTypeMapper() JobDataTypeMapper {
	return JobDataTypeMapper{
		job.UserProfileCreatedJobBegin:     reflect.TypeOf(EventResponseMessageWithData[job.CreateUserProfileJobBeganData]{}),
		job.UserPreferenceUpdatedJobBegin:  reflect.TypeOf(EventResponseMessageWithData[job.UpdateUserPreferenceJobBeganData]{}),
		job.TeamCreatedJobBegin:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTeamJobBeganData]{}),
		job.TeamUsersAddedJobBegin:         reflect.TypeOf(EventResponseMessageWithData[job.AddUsersTeamJobBeganData]{}),
		job.OrganizationCreatedJobBegin:    reflect.TypeOf(EventResponseMessageWithData[job.CreateOrganizationJobBeganData]{}),
		job.OrganizationUsersAddedJobBegin: reflect.TypeOf(EventResponseMessageWithData[job.AddUsersOrganizationJobBeganData]{}),
		job.TaskCreatedJobBegin:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTaskJobBeganData]{}),
		job.TaskStatusUpdatedJobBegin:      reflect.TypeOf(EventResponseMessageWithData[job.UpdateStatusTaskJobBeganData]{}),
		job.FileObjectCreatedJobBegin:      reflect.TypeOf(EventResponseMessageWithData[job.CreateFileObjectJobBeganData]{}),

		job.UserProfileCreatedJobProcessed:     reflect.TypeOf(EventResponseMessageWithData[job.CreateUserProfileJobProcessedData]{}),
		job.UserPreferenceUpdatedJobProcessed:  reflect.TypeOf(EventResponseMessageWithData[job.UpdateUserPreferenceJobProcessedData]{}),
		job.TeamCreatedJobProcessed:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTeamJobProcessedData]{}),
		job.TeamUsersAddedJobProcessed:         reflect.TypeOf(EventResponseMessageWithData[job.AddUsersTeamJobProcessedData]{}),
		job.OrganizationCreatedJobProcessed:    reflect.TypeOf(EventResponseMessageWithData[job.CreateOrganizationJobProcessedData]{}),
		job.OrganizationUsersAddedJobProcessed: reflect.TypeOf(EventResponseMessageWithData[job.AddUsersOrganizationJobProcessedData]{}),
		job.TaskCreatedJobProcessed:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTaskJobProcessedData]{}),
		job.TaskStatusUpdatedJobProcessed:      reflect.TypeOf(EventResponseMessageWithData[job.UpdateStatusTaskJobProcessedData]{}),
		job.FileObjectCreatedJobProcessed:      reflect.TypeOf(EventResponseMessageWithData[job.CreateFileObjectJobProcessedData]{}),

		job.UserProfileCreatedJobFailed:     reflect.TypeOf(EventResponseMessageWithData[job.CreateUserProfileJobFailedData]{}),
		job.UserPreferenceUpdatedJobFailed:  reflect.TypeOf(EventResponseMessageWithData[job.UpdateUserPreferenceJobFailedData]{}),
		job.TeamCreatedJobFailed:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTeamJobFailedData]{}),
		job.TeamUsersAddedJobFailed:         reflect.TypeOf(EventResponseMessageWithData[job.AddUsersTeamJobFailedData]{}),
		job.OrganizationCreatedJobFailed:    reflect.TypeOf(EventResponseMessageWithData[job.CreateOrganizationJobFailedData]{}),
		job.OrganizationUsersAddedJobFailed: reflect.TypeOf(EventResponseMessageWithData[job.AddUsersOrganizationJobFailedData]{}),
		job.TaskCreatedJobFailed:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTaskJobFailedData]{}),
		job.TaskStatusUpdatedJobFailed:      reflect.TypeOf(EventResponseMessageWithData[job.UpdateStatusTaskJobFailedData]{}),
		job.FileObjectCreatedJobFailed:      reflect.TypeOf(EventResponseMessageWithData[job.CreateFileObjectJobFailedData]{}),

		job.UserProfileCreatedJobSuccess:     reflect.TypeOf(EventResponseMessageWithData[job.CreateUserProfileJobSuccessData]{}),
		job.UserPreferenceUpdatedJobSuccess:  reflect.TypeOf(EventResponseMessageWithData[job.UpdateUserPreferenceJobSuccessData]{}),
		job.TeamCreatedJobSuccess:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTeamJobSuccessData]{}),
		job.TeamUsersAddedJobSuccess:         reflect.TypeOf(EventResponseMessageWithData[job.AddUsersTeamJobSuccessData]{}),
		job.OrganizationCreatedJobSuccess:    reflect.TypeOf(EventResponseMessageWithData[job.CreateOrganizationJobSuccessData]{}),
		job.OrganizationUsersAddedJobSuccess: reflect.TypeOf(EventResponseMessageWithData[job.AddUsersOrganizationJobSuccessData]{}),
		job.TaskCreatedJobSuccess:            reflect.TypeOf(EventResponseMessageWithData[job.CreateTaskJobSuccessData]{}),
		job.TaskStatusUpdatedJobSuccess:      reflect.TypeOf(EventResponseMessageWithData[job.UpdateStatusTaskJobSuccessData]{}),
		job.FileObjectCreatedJobSuccess:      reflect.TypeOf(EventResponseMessageWithData[job.CreateFileObjectJobSuccessData]{}),
	}
}
