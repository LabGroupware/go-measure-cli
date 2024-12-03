package domain

type FileObjectDto struct {
	FileObjectID  string                `json:"fileObjectId"`
	BucketID      string                `json:"bucketId"`
	Name          string                `json:"name"`
	Path          string                `json:"path"`
	MimeType      string                `json:"mimeType"`
	Size          int64                 `json:"size"`
	Checksum      string                `json:"checksum"`
	AttachedTasks []TaskOnFileObjectDto `json:"attachedTasks"`
}

type FileObjectOnTaskDto struct {
	FileObjectDto `json:",inline"`
}
