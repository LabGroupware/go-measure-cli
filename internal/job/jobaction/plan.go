package jobaction

type CreateTaskAction string

const (
	CREATE_TASK_VALIDATE_TASK        CreateTaskAction = "VALIDATE_TASK"
	CREATE_TASK_VALIDATE_USER        CreateTaskAction = "VALIDATE_USER"
	CREATE_TASK_VALIDATE_TEAM        CreateTaskAction = "VALIDATE_TEAM"
	CREATE_TASK_VALIDATE_FILE_OBJECT CreateTaskAction = "VALIDATE_FILE_OBJECT"
	CREATE_TASK_CREATE_TASK          CreateTaskAction = "CREATE_TASK_AND_ATTACH_FILE_OBJECT"
)

type UpdateStatusTaskAction string

const (
	UPDATE_STATUS_TASK_VALIDATE_TASK UpdateStatusTaskAction = "VALIDATE_TASK"
	UPDATE_STATUS_TASK_UPDATE_TASK   UpdateStatusTaskAction = "UPDATE_TASK_STATUS"
)
