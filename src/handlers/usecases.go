package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/logging"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/worker"
)

const (
	NotFoundKey              = "NOT_FOUND"
	idKey                    = "id"
	maxPayloadSize           = 1024 * 1024
	tooLargeRequestBodyError = "RequestEntity Too Large"
	renderingDataError       = "error during rendering url data: %s"
	savingDataError          = "error during saving url data: %s"
	deletingDataError        = "error during deleting url data: %s"
	gettingDataError         = "error during getting url data: %s"
)

func SaveURLData(repo database.Repository, wrk worker.Worker) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rawData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logging.Log.Errorf(gettingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(rawData) > maxPayloadSize {
			logging.Log.Error(tooLargeRequestBodyError)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		var urlData model.URLData
		if err = json.Unmarshal(rawData, &urlData); err != nil {
			logging.Log.Errorf(gettingDataError, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		savedID, err := repo.SaveURLData(&urlData)
		if err != nil {
			logging.Log.Errorf(savingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(time.Duration(urlData.Interval) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					dwnldHistory := wrk.DoRequest(urlData.URL, savedID)
					err = wrk.SaveDownloadHistory(dwnldHistory)
					if err != nil {
						logging.Log.Errorf(savingDataError, err.Error())
					}
				}
			}
		}()

		if err = renderJSON(w, model.SaveURLDataResponse{ID: savedID}); err != nil {
			logging.Log.Errorf(renderingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

func DeleteURLData(repo database.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stringID := vars[idKey]

		ID, err := strconv.ParseInt(stringID, 10, 64)
		if err != nil {
			logging.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = repo.DeleteURLData(ID); err != nil {
			if err.Error() == NotFoundKey {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logging.Log.Errorf(deletingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetAllURLData(repo database.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := repo.GetAllURLData()
		if err != nil {
			logging.Log.Errorf(gettingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = renderJSON(w, data); err != nil {
			logging.Log.Errorf(renderingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

func GetDownloadHistoriesByURLDataID(repo database.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stringID := vars[idKey]
		data, err := repo.GetDownloadHistoriesByURLDataID(stringID)
		if err != nil {
			logging.Log.Errorf(gettingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		renderedData := []model.RenderedDownloadHistory{}
		for _, elem := range data {
			renderedElem := model.RenderedDownloadHistory{
				Response:  elem.Response,
				Duration:  elem.Duration,
				CreatedAt: elem.CreatedAt,
			}
			renderedData = append(renderedData, renderedElem)
		}
		if err = renderJSON(w, renderedData); err != nil {
			logging.Log.Errorf(renderingDataError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

func renderJSON(w http.ResponseWriter, response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
