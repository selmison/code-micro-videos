// +build integration

package rest_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/bxcodec/faker/v3"

	"github.com/selmison/code-micro-videos/testdata"
)

func Test_RestApi_Post_Categories(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "categories")
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
			name: "When name field is empty",
			req: request{
				fakeUrl,
				"application/json; charset=UTF-8",
				strings.NewReader(fmt.Sprintf(
					`{"name": "%s" }`,
					"",
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

func Test_RestApi_Get_Categories(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "categories")
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
				body:   toJSON(testdata.FakeCategoriesDTO),
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
				bs, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Errorf("read body: %v", err)
					return
				}
				if bytes.Equal(bs, tt.want.body) {
					t.Errorf("\nbody: %v\nwant: %v", string(bs), string(tt.want.body))
				}
			}
		})
	}
}

func Test_RestApi_Get_Category(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[0]
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "categories", name)
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
			name: "When name doesn't exist",
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
			name: "When name exists",
			req: request{
				url:         fakeUrl(fakeExistName),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusOK,
				body:   toJSON(fakeExistCategoryDTO),
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
				if bytes.Equal(bs, tt.want.body) {
					t.Errorf("\nbody: %v\nwant: %v", string(bs), string(tt.want.body))
				}
			}
		})
	}
}

func Test_RestApi_Delete_Category(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "categories", name)
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
			name: "When name doesn't exist",
			req: request{
				url:         fakeUrl(fakeDoesNotExistName),
				contentType: "application/json; charset=UTF-8",
			},
			want: response{
				status: http.StatusNotFound,
				body:   "Not Found",
			},
			wantErr: false,
		},
		{
			name: "When name exists",
			req: request{
				url:         fakeUrl(fakeExistName),
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

func Test_RestApi_Update_Category(t *testing.T) {
	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "categories", name)
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
			name: "When name doesn't exist",
			req: request{
				url:         fakeUrl(fakeDoesNotExistName),
				contentType: "application/json; charset=UTF-8",
				body: strings.NewReader(fmt.Sprintf(
					`{"name": "%s"}`,
					faker.Name(),
				)),
			},
			want: response{
				status: http.StatusNotFound,
				body:   "Not Found",
			},
			wantErr: false,
		},
		{
			name: "When name exists",
			req: request{
				url:         fakeUrl(fakeExistName),
				contentType: "application/json; charset=UTF-8",
				body: strings.NewReader(fmt.Sprintf(
					`{"name": "%s", "description": "%s"}`,
					faker.Name(),
					faker.Sentence(),
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
