package worker

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/logging"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
)

const WorkerTimeout = 5

// Worker retrieves data from urls, downloads data from "url" with the given "interval" in background.
// In case of timeout or another error, null should be written to the "response" fieldt he whole response 
// shall be saved to the "response" field.
type Worker struct {
	cl   *http.Client
	repo database.Repository
}

func NewWorker(cl *http.Client, repo database.Repository) Worker {
	return Worker{cl: cl, repo: repo}
}

func (w Worker) DoRequest(url string, urlDataID int64) model.DownloadHistory {
	
	startTime := float64(time.Now().UnixNano()) / float64(time.Second)
	response, err := w.cl.Get(url)
	endTime := float64(time.Now().UnixNano()) / float64(time.Second)	
	if err != nil {
		logging.Log.Errorf("error during request: %s", err.Error())
		return model.DownloadHistory{
			URLDataID: urlDataID,
			Response:  nil,
			CreatedAt: setFloatPrecision(float64(time.Now().Unix())/float64(time.Second), 3),
			Duration:  WorkerTimeout,
		}
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Log.Errorf("error during readind response body: %s", err.Error())
		return model.DownloadHistory{
			URLDataID: urlDataID,
			Response:  nil,
			CreatedAt: setFloatPrecision(float64(time.Now().Unix())/float64(time.Second), 3),
			Duration:  WorkerTimeout,
		}
	}
	respBody := string(data)
	return model.DownloadHistory{
		URLDataID: urlDataID,
		Response:  &respBody,
		CreatedAt: setFloatPrecision(float64(time.Now().Unix())/float64(time.Second), 3),
		Duration:  setFloatPrecision((endTime - startTime), 3),
	}
}

func (w Worker) SaveDownloadHispory(downloadHistory model.DownloadHistory) error {	
	return w.repo.SaveDownloadHistory(downloadHistory)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func setFloatPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
