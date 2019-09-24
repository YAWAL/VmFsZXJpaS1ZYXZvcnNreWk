package database

import (
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
)

// Repository is an interface for different DB storages to proceed with url data
type Repository interface {
	SaveURLData(data *model.URLData) (savedID int64, err error)
	DeleteURLData(URLDataID int64) error
	GetAllURLData() (data []model.URLData, err error)
	GetDownloadHistoriesByURLDataID(ID string) (downloadHistory []model.DownloadHistory, err error)
	SaveDownloadHistory(downloadHistory model.DownloadHistory) error
}
