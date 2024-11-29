package job

import (
	"github.com/LabGroupware/go-measure-tui/internal/job/jobaction"
	"github.com/LabGroupware/go-measure-tui/internal/job/jobdata"
)

type CreateFileObjectJobBeganData JobBeganData[jobaction.CreateFileObjectAction]
type CreateFileObjectJobProcessedData JobProcessedData[jobaction.CreateFileObjectAction, jobdata.CreateFileObjectResultData]
type CreateFileObjectJobSuccessData JobSuccessData[jobdata.CreateFileObjectResultData, jobaction.CreateFileObjectAction]
type CreateFileObjectJobFailedData JobFailedData[jobaction.CreateFileObjectAction]
