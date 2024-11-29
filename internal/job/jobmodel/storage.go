package jobmodel

type FileObject struct {
	FileObjectID string `json:"fileObjectId"`
	BucketID     string `json:"bucketId"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mimeType"`
	Checksum     string `json:"checksum"`
}

type FileObjectOnTask struct {
	TaskAttachmentID string `json:"taskAttachmentId"`
	FileObjectID     string `json:"fileObjectId"`
}
