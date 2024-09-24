package response

type SettingRes struct {
	SettingCode  string `json:"setting_code"`
	SetGroup     string `json:"set_group"`
	SetKey       string `json:"set_key"`
	SetLabel     string `json:"set_label"`
	SetOrder     int    `json:"set_order"`
	ContentType  string `json:"content_type"`
	ContentValue string `json:"content_value"`
	IsActive     bool   `json:"is_active"`
}
