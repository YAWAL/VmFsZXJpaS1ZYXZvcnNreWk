package worker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database/mocks"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/model"
)

func Test_round(t *testing.T) {

	tests := []struct {
		name string
		num  float64
		want int
	}{
		{
			name: "with 3 extra digits rounded up",
			num:  476746.584,
			want: 476747,
		},
		{
			name: "with 3 extra digits rounded up",
			num:  476746.484,
			want: 476746,
		},
		{
			name: "with no extra digits",
			num:  476746,
			want: 476746,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := round(tt.num); got != tt.want {
				t.Errorf("round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setFloatPrecision(t *testing.T) {

	tests := []struct {
		name      string
		num       float64
		precision int
		want      float64
	}{
		{
			name:      "1",
			num:       123.556,
			precision: 2,
			want:      123.56,
		},
		{
			name:      "2",
			num:       123.555,
			precision: 2,
			want:      123.56,
		},
		{
			name:      "3",
			num:       123.546,
			precision: 2,
			want:      123.55,
		},
		{
			name:      "4",
			num:       123.544,
			precision: 2,
			want:      123.54,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setFloatPrecision(tt.num, tt.precision); got != tt.want {
				t.Errorf("setFloatPrecision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {

	var cl *http.Client
	var repo database.Repository

	tests := []struct {
		name string
		want Worker
	}{
		{
			name: "valid worker creation",
			want: Worker{cl: cl, repo: repo},
		},
		{
			name: "valid worker with nil fields",
			want: Worker{cl: nil, repo: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(cl, repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_SaveDownloadHispory(t *testing.T) {

	repoMock := &mocks.Repository{}
	var cl *http.Client
	var responce string
	tests := []struct {
		name            string
		w               Worker
		downloadHistory model.DownloadHistory
		wantErr         bool
	}{
		{
			name: "valid saving of downloaded history",
			w:    Worker{cl: cl, repo: repoMock},
			downloadHistory: model.DownloadHistory{
				ID:        1,
				URLDataID: 1,
				Response:  &responce,
				Duration:  1,
				CreatedAt: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		repoMock.On("SaveDownloadHistory", tt.downloadHistory).Return(nil)
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.w.SaveDownloadHistory(tt.downloadHistory); (err != nil) != tt.wantErr {
				t.Errorf("Worker.SaveDownloadHispory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestWorker_DoRequest(t *testing.T) {

// 	repoMock := &mocks.Repository{}
// 	var cl *http.Client
// 	var responce string

// 	tests := []struct {
// 		name      string
// 		w         Worker
// 		url       string
// 		urlDataID int64
// 		want      model.DownloadHistory
// 	}{
// 		{
// 			name:      "string",
// 			url:       "http://blog.golang.org/cover",
// 			urlDataID: 1,
// 			want: model.DownloadHistory{
// 				ID:        1,
// 				URLDataID: 1,
// 				Response:  &responce,
// 				Duration:  1,
// 				CreatedAt: 1,
// 			},
// 		},
// 	}
// 	w := Worker{cl: cl, repo: repoMock}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			fmt.Println(tt.url)
// 			fmt.Println(tt.urlDataID)

// 			got := w.DoRequest(tt.url, tt.urlDataID)

// 			t.Log(got)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Worker.DoRequest() = %v, want %v", got, tt.want)
// 			}

// 		})
// 	}
// }

// go test -coverprofile=cov.out
// go tool cover -html=cov.out

func TestWorker_DoRequest_(t *testing.T) {

	repoMock := &mocks.Repository{}

	tests := []struct {
		name          string
		w             Worker
		url           string
		urlDataID     int64
		want          model.DownloadHistory
		wantURLDataID int64
	}{
		{
			name:      "Invalid URL case 1 -> Duration set to 5s",
			url:       "/some/url",
			urlDataID: 1,
			want: model.DownloadHistory{
				ID:        0,
				URLDataID: 1,
				Response:  nil,
				Duration:  5,
			},
			wantURLDataID: 1,
		},
		{
			name:      "Invalid URL case 2 -> Duration set to 5s",
			url:       "https://unknown",
			urlDataID: 1,
			want: model.DownloadHistory{
				ID:        0,
				URLDataID: 2,
				Response:  nil,
				Duration:  5,
			},
			wantURLDataID: 2,
		},
		{
			name:      "Invalid URL case 3 -> Duration set to 5s",
			url:       "https://unknown",
			urlDataID: 1,
			want: model.DownloadHistory{
				ID:        0,
				URLDataID: 3,
			},
			wantURLDataID: 3,
		},
		{
			name:      "Valid URL",
			url:       "https://golang.org/pkg/fmt/",
			urlDataID: 1,
			want: model.DownloadHistory{
				ID:        0,
				URLDataID: 4,
			},
			wantURLDataID: 4,
		},
	}
	for _, tt := range tests {

		// Start a local HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			equals(t, req.URL.String(), tt.url)
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}))
		// Close the server when test finishes
		defer server.Close()
		w := Worker{cl: server.Client(), repo: repoMock}

		got := w.DoRequest(tt.url, tt.want.URLDataID)

		// resp := got.Response

		// fmt.Printf("%v", *resp)
		equals(t, tt.want.ID, got.ID)
		equals(t, tt.want.URLDataID, got.URLDataID)

	}

}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
