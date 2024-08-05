package v1

type ShengCaiListRequest struct {
	SpreadsheetToken string `json:"sheet_id" validate:"required"`
	Page             int    `json:"page" validate:"required,min=1"`
}

type ShengCaiListResponse struct {
	TotalCount int `json:"total_count"`
	List       []struct {
		SheetID     string `json:"sheet_id"`
		Title       string `json:"title"`
		Link        string `json:"link"`
		ReleaseDate string `json:"release_date"`
		Abstract    string `json:"abstract"`
		Keyword     string `json:"keyword"`
	} `json:"list"`
}

type ShengCaiGetMetaDataRequest struct {
	SpreadsheetToken string `json:"sheet_id" validate:"required"`
}

type ShengCaiGetMetaDataResponse struct {
	SheetName string `json:"sheet_name"`
	UpdateLog string `json:"update_log"`
}
