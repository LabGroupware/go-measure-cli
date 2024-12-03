package response

type ResponseDto[T any] struct {
	Success         bool              `json:"success"`
	Data            T                 `json:"data"`
	Code            string            `json:"code"`
	Caption         string            `json:"caption"`
	ErrorAttributes ErrorAttributeDto `json:"errorAttributes"`
}

type ErrorAttributeDto struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Value     any    `json:"value"`
}

type ErrorResponseDto ResponseDto[any]

type ListResponseData[T any] struct {
	ListData  []T       `json:"listData"`
	CountData CountData `json:"countData"`
}

type CountData struct {
	Count   int  `json:"count"`
	IsValid bool `json:"isValid"`
}

type ListResponseDto[T any] ResponseDto[ListResponseData[T]]

type CommandResponseData struct {
	JobID string `json:"jobId"`
}

type CommandResponseDto ResponseDto[CommandResponseData]
