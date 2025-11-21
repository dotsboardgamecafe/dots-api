package response

type UserPointHistoryResponse struct {
	SourceId       int64  `json:"source_id"`
	SourceUserCode string `json:"source_user_code"`
	SourceType     string `json:"source_type"`
	SourceCode     string `json:"source_code"`
	SourceName     string `json:"source_name"`
	Point          int    `json:"point"`
	CreatedDate    string `json:"created_date"`
}
