package model

type URLData struct {
	ID       int64  `json:"id" gorm:"PRIMARY_KEY"`
	URL      string `json:"url"`
	Interval int64  `json:"interval"`
}

type DownloadHistory struct {
	ID int64 `json:"id"  gorm:"PRIMARY_KEY"`
	URLDataID int64   `json:"urldata_id"`
	Response  *string `json:"response"`
	Duration  float64 `json:"duration"  gorm:"type:numeric(4,3)"`
	CreatedAt float64 `json:"created_at" gorm:"type:float"`
}

type RenderedDownloadHistory struct {
	Response  *string `json:"response"`
	Duration  float64 `json:"duration"`
	CreatedAt float64 `json:"created_at"`
}

type SaveURLDataResponse struct {
	ID int64 `json:"id"`
}

type GetAllURLDataResponse struct {
	ID int64 `json:"id"`
}
