package jobmodel

type UserPreference struct {
	UserPreferenceID      string `json:"userPreferenceId"`
	UserID                string `json:"userId"`
	Timezone              string `json:"timezone"`
	Theme                 string `json:"theme"`
	Language              string `json:"language"`
	NotificationSettingID string `json:"notificationSettingId"`
}
