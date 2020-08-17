// +build integration

package rest_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/bxcodec/faker/v3"

	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/testdata"
)

func Test_integration_VideoCreate(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	type request struct {
		url         string
		contentType string
		body        io.Reader
	}
	type response struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "create a video",
			req: request{
				fmt.Sprintf("http://%s/%s", cfg.AddressServer, "videos"),
				"application/json; charset=UTF-8",
				strings.NewReader(fmt.Sprintf(
					`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d}`,
					faker.Name(),
					faker.Sentence(),
					2020,
					false,
					crud.TenRating,
					250,
				)),
			},
			want: response{
				status: http.StatusCreated,
				body:   `"Created"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Post(tt.req.url, tt.req.contentType, tt.req.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
				}
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
				}
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("Body: %v, want: %v", data, tt.want.body)
				}
			}
		})
	}
}

func Test_RestApi_Post_Videos(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "videos")
	type request struct {
		url         string
		contentType string
		body        io.Reader
	}
	type response struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "When title field is empty",
			req: request{
				fakeUrl,
				"application/json; charset=UTF-8",
				strings.NewReader(fmt.Sprintf(
					`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d}`,
					"",
					faker.Sentence(),
					2020,
					false,
					crud.TenRating,
					250,
				)),
			},
			want: response{
				status: http.StatusBadRequest,
				body:   http.StatusText(http.StatusBadRequest),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Post(tt.req.url, tt.req.contentType, tt.req.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
					return
				}
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("body: %v, want: %v", data, tt.want.body)
				}
			}
		})
	}
}

func Test_RestApi_Get_Videos(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "videos")
	type request struct {
		url         string
		contentType string
	}
	type response struct {
		status int
		body   []byte
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "When everything is right",
			req: request{
				url:         fakeUrl,
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusOK,
				body:   toJSON(testdata.FakeVideosDTO),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Get(tt.req.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
					return
				}
				data, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				comp, err := JSONBytesEqual(data, tt.want.body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				if !comp {
					t.Errorf("\nresponse: %vwant: %v", string(data), string(tt.want.body))
				}
			}
		})
	}
}

func Test_RestApi_Get_Video(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	fakeDoesNotExistName := "doesNotExistName"
	fakeExistVideoDTO := testdata.FakeVideosDTO[0]
	fakeUrl := func(title string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
	}
	type request struct {
		url         string
		contentType string
	}
	type response struct {
		status int
		body   []byte
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "When title doesn't exist",
			req: request{
				url:         fakeUrl(fakeDoesNotExistName),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusNotFound,
				body:   toJSON("Not Found"),
			},
			wantErr: false,
		},
		{
			name: "When title exists",
			req: request{
				url:         fakeUrl(fakeExistTitle),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusOK,
				body:   toJSON(fakeExistVideoDTO),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Get(tt.req.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
					return
				}
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				bodyResponse, err := json.Marshal(bs)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				comp, err := JSONBytesEqual(bodyResponse, tt.want.body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				if comp {
					t.Errorf("\nresponse: %vwant: %v", string(bodyResponse), string(tt.want.body))
				}
			}
		})
	}
}

func Test_RestApi_Delete_Video(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	fakeDoesNotExistTitle := "doesNotExistTitle"
	fakeUrl := func(title string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
	}
	type request struct {
		url         string
		contentType string
	}
	type response struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "When title doesn't exist",
			req: request{
				url:         fakeUrl(fakeDoesNotExistTitle),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusNotFound,
				body:   "Not Found",
			},
			wantErr: false,
		},
		{
			name: "When title exists",
			req: request{
				url:         fakeUrl(fakeExistTitle),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodDelete, tt.req.url, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			got, err := client.Do(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
					return
				}
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("\nbody: %v\nwant: %v", data, tt.want.body)
				}
			}
		})
	}
}

func Test_RestApi_Update_Video(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	fakeDoesNotExistTitle := "doesNotExistTitle"
	fakeUrl := func(title string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
	}
	type request struct {
		url         string
		contentType string
		body        io.Reader
	}
	type response struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		req     request
		want    response
		wantErr bool
	}{
		{
			name: "When title doesn't exist",
			req: request{
				url:         fakeUrl(fakeDoesNotExistTitle),
				contentType: "application/json; charset=UTF-8",
				body: strings.NewReader(fmt.Sprintf(
					`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d}`,
					faker.Name(),
					faker.Sentence(),
					2020,
					false,
					crud.TenRating,
					250,
				)),
			},
			want: response{
				status: http.StatusNotFound,
				body:   "Not Found",
			},
			wantErr: false,
		},
		{
			name: "When title exists",
			req: request{
				url:         fakeUrl(fakeExistTitle),
				contentType: "application/json; charset=UTF-8",
				body: strings.NewReader(fmt.Sprintf(
					`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d}`,
					faker.Sentence(),
					faker.Sentence(),
					2020,
					false,
					crud.TenRating,
					250,
				)),
			},
			want: response{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, tt.req.url, tt.req.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			got, err := client.Do(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.StatusCode != tt.want.status {
					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
					return
				}
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("\nbody: %v\nwant: %v", data, tt.want.body)
				}
			}
		})
	}
}
