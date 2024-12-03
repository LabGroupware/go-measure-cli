package request

type UpdateUserPreferenceRequestDto struct {
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
}
