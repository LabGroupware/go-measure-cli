package request

type CreateFileObjectRequestDto struct {
	BucketID string `json:"bucketId"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}
