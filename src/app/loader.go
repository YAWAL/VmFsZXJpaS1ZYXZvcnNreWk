package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/config"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/logging"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/repository"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/router"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/worker"
)

type URLservis struct {
	Server *http.Server
	db     *gorm.DB
}

// LoadApp performs reading of  config file, initializes connection to database,
// performs initial table migrations, initializes Postgres repository, initializes routers
// and creates server
func LoadApp(file string) (servis *URLservis, err error) {
	// read, config
	conf, err := config.ReadConfig(file)
	if err != nil {
		logging.Log.Errorf("cannot load config: %v", err.Error())
		return nil, err
	}
	// establish connections to DB
	db, err := database.PGconn(conf.Database)
	if err != nil {
		logging.Log.Errorf("cannot connect to DB: %s", err.Error())
		return nil, err
	}
	// migrate tables
	db.AutoMigrate(&model.DownloadHistory{}, &model.URLData{})
	logging.Log.Infof("Auto migration")

	// init repos
	repo := repository.NewPostgresRepository(db)
	logging.Log.Infof("Initializing Postgres repository")

	// create HTTP client with timeout
	cl := &http.Client{
		Timeout: worker.Timeout * time.Second,
	}
	// init worker
	worker := worker.New(cl, repo)

	// init routers
	r := router.New(repo, worker)
	logging.Log.Infof("Initializing routers")

	srv := &http.Server{
		Handler:      r,
		Addr:         conf.Host,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	logging.Log.Infof("Application is running on %s ", conf.Host)

	return &URLservis{db: db, Server: srv}, nil
}

// GracefullShutdown performs grecefull shutdown of server and loggs to stdout
// info about operation success
func GracefullShutdown(servis *URLservis, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logging.Log.Info("Servis is shutting down")

	if err := servis.db.Close(); err != nil {
		logging.Log.Errorf("can not gracefully close connections to database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	servis.Server.SetKeepAlivesEnabled(false)
	if err := servis.Server.Shutdown(ctx); err != nil {
		logging.Log.Errorf("can not gracefully shutdown the server: %v", err)
	}
	close(done)
}
