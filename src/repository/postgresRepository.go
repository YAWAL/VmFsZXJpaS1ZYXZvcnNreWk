package repository

import (
	"fmt"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/handlers"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"

	"github.com/jinzhu/gorm"
)

// PostgresRepository is a struct for Postgres repository which satisfies Repository interface
type PostgresRepository struct {
	conn *gorm.DB
}

// NewPostgresRepository creates PostgresRepository
func NewPostgresRepository(conn *gorm.DB) PostgresRepository {
	return PostgresRepository{conn: conn}
}

func (pg PostgresRepository) SaveURLData(data *model.URLData) (savedID int64, err error) {
	// check if record exists
	var existedData model.URLData
	pg.conn.Where(data).First(&existedData)
	// if exists - perform updating
	if existedData.ID > 0 {
		if err = pg.conn.Save(existedData).Error; err != nil {
			return 0, err
		}
		return existedData.ID, nil
	}

	if err = pg.conn.Create(data).Error; err != nil {
		return 0, err
	}
	// retrieve newly created URLData's ID
	if err = pg.conn.Where(data).First(&existedData).Error; err != nil {
		return 0, err
	}
	return existedData.ID, nil
}

func (pg PostgresRepository) DeleteURLData(URLDataID int64) (err error) {
	urlData := model.URLData{ID: URLDataID}
	if err = pg.conn.First(&model.URLData{ID: URLDataID}).Error; err != nil {
		return fmt.Errorf(handlers.NotFoundKey)
	}

	if err = pg.conn.Delete(&urlData).Error; err != nil {
		return err
	}

	return nil
}

func (pg PostgresRepository) GetAllURLData() (data []model.URLData, err error) {
	if err = pg.conn.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (pg PostgresRepository) GetDownloadHistoriesByURLDataID(ID string) (downloadHistory []model.DownloadHistory, err error) {
	// find by url_data_id
	if err = pg.conn.Where("url_data_id = ?", ID).Find(&downloadHistory).Error; err != nil {
		return nil, err
	}
	return downloadHistory, nil
}

func (pg PostgresRepository) SaveDownloadHistory(downloadHistory model.DownloadHistory) (err error) {
	return pg.conn.Create(&downloadHistory).Error
}
