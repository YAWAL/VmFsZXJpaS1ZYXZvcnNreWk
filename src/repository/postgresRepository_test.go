package repository

import (
	"database/sql"
	"testing"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type PostgresRepoSuite struct {
	suite.Suite
	DB              *gorm.DB
	mock            sqlmock.Sqlmock
	repository      PostgresRepository
	downloadHistory *model.DownloadHistory
}

func TestNewPostgresRepository(t *testing.T) {
	got := NewPostgresRepository(nil)
	t.Log(got)
}

func (s *PostgresRepoSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *PostgresRepoSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = NewPostgresRepository(s.DB)
}

func (s *PostgresRepoSuite) TestGetDownloadHistoriesByURLDataID() {
	var (
		id        int64 = 1
		URLDataID int64 = 1
	)

	s.mock.ExpectQuery(`SELECT (.+) FROM "download_histories" WHERE (.+)`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"url_data_id", "1"}).AddRow(URLDataID, id))
	_, err := s.repository.GetDownloadHistoriesByURLDataID("1")
	require.NoError(s.T(), err)
}

func (s *PostgresRepoSuite) TestPostgresRepository_GetAllURLData() {
	s.mock.ExpectQuery(`SELECT (.+) FROM "url_data"`).WillReturnRows(sqlmock.NewRows([]string{}))
	_, err := s.repository.GetAllURLData()
	require.NoError(s.T(), err)
}

func (s *PostgresRepoSuite) TestPostgresRepository_SaveDownloadHistory() {
	err := s.repository.SaveDownloadHistory(model.DownloadHistory{ID: 1, URLDataID: 1})
	require.Error(s.T(), err)
}

func (s *PostgresRepoSuite) TestPostgresRepository_DeleteURLData() {
	err := s.repository.DeleteURLData(1)
	require.Error(s.T(), err)
}

func (s *PostgresRepoSuite) TestPostgresRepository_SaveURLData() {
	_, err := s.repository.SaveURLData(&model.URLData{})
	require.Error(s.T(), err)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}
