// +build integration

package rest_test

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"io/ioutil"
//	"log"
//	"mime/multipart"
//	"net/http"
//	"strconv"
//	"strings"
//	"testing"
//
//	_ "github.com/lib/pq"
//	"github.com/stretchr/testify/assert"
//	"github.com/volatiletech/sqlboiler/v4/boil"
//
//	"github.com/bxcodec/faker/v3"
//
//	"github.com/selmison/code-micro-videos/pkg/api/rest"
//	"github.com/selmison/code-micro-videos/pkg/storage/files/memory"
//	"github.com/selmison/code-micro-videos/pkg/video"
//	"github.com/selmison/code-micro-videos/testdata"
//)
//
//func Test_RestApi_Post_Videos(t *testing.T) {
//	cfg, teardownTestCase, err := setupTestCase(t, nil)
//	if err != nil {
//		t.Fatalf("test: failed to setup test case: %v\n", err)
//	}
//	defer teardownTestCase(t)
//	seed, teardownTestCase, err := memory.SetupFileTestCase()
//	if err != nil {
//		t.Fatalf("test: failed to setup files test case: %v\n", err)
//	}
//	const (
//		fakeCategoryIndex = 0
//		fakeGenreIndex    = 0
//	)
//	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "videos")
//	fakeExistGenreName := testdata.FakeNewGenres[fakeGenreIndex].Name
//	fakeExistCategoryName := testdata.FakeCategories[fakeCategoryIndex].Name
//	fakeExistCategoryDescription := testdata.FakeCategories[fakeCategoryIndex].Description
//	fakeDoesNotExistGenreName := faker.FirstName()
//	fakeDoesNotExistCategoryName := faker.FirstName()
//	fakeDoesNotExistCategoryDescription := faker.Sentence()
//	type request struct {
//		url             string
//		fields          map[string]string
//		fieldNameToFile string
//		pathFile        string
//	}
//	type response struct {
//		status int
//		body   string
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    response
//		wantErr bool
//	}{
//		{
//			name: "When the title in body is blank",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    "",
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusBadRequest,
//				body:   http.StatusText(http.StatusBadRequest),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When the title in body is omitted",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusBadRequest,
//				body:   http.StatusText(http.StatusBadRequest),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When the year_launched in body is omitted",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusBadRequest,
//				body:   http.StatusText(http.StatusBadRequest),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When the opened in body is omitted",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusCreated,
//				body:   http.StatusText(http.StatusCreated),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When the rating in body is omitted",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusBadRequest,
//				body:   http.StatusText(http.StatusBadRequest),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When the duration in body is blank",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusBadRequest,
//				body:   http.StatusText(http.StatusBadRequest),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When Video is with wrong categories and genres",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeDoesNotExistGenreName,
//					"categories.0.name":        fakeDoesNotExistCategoryName,
//					"categories.0.description": fakeDoesNotExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusNotFound,
//				body:   http.StatusText(http.StatusNotFound),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When everything is right",
//			req: request{
//				url: fakeUrl,
//				fields: map[string]string{
//					"title":                    faker.Name(),
//					"description":              faker.Sentence(),
//					"year_launched":            "2020",
//					"opened":                   "false",
//					"rating":                   strconv.Itoa(int(video.TenRating)),
//					"duration":                 "250",
//					"genres.0.name":            fakeExistGenreName,
//					"categories.0.name":        fakeExistCategoryName,
//					"categories.0.description": fakeExistCategoryDescription,
//				},
//				fieldNameToFile: rest.VideoFileField,
//				pathFile:        seed.FakeTmpFilePath,
//			},
//			want: response{
//				status: http.StatusCreated,
//				body:   http.StatusText(http.StatusCreated),
//			},
//			wantErr: false,
//		},
//	}
//	ctx := context.Background()
//	fakeExistCategory := testdata.FakeCategories[fakeCategoryIndex]
//	err = fakeExistCategory.InsertG(ctx, boil.Infer())
//	if err != nil {
//		t.Errorf("test: insert category: %s", err)
//		return
//	}
//	fakeExistGenre := testdata.FakeGenres[fakeGenreIndex]
//	err = fakeExistGenre.InsertG(ctx, boil.Infer())
//	if err != nil {
//		t.Errorf("test: insert genre: %s", err)
//		return
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := postMultipartWithFile(
//				tt.req.url,
//				tt.req.fields,
//				seed,
//				tt.req.fieldNameToFile,
//				tt.req.pathFile,
//			)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				bs, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Errorf("read body: %v", err)
//					return
//				}
//				data := strings.TrimSpace(string(bs))
//				if data != tt.want.body {
//					t.Errorf("body: %v, want: %v", data, tt.want.body)
//				}
//			}
//		})
//	}
//}
//
//func Test_RestApi_Get_Videos(t *testing.T) {
//	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
//	if err != nil {
//		t.Errorf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	fakeUrl := fmt.Sprintf("http://%s/%s", cfg.AddressServer, "videos")
//	type request struct {
//		url         string
//		contentType string
//	}
//	type response struct {
//		status int
//		body   []byte
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    response
//		wantErr bool
//	}{
//		{
//			name: "When everything is right",
//			req: request{
//				url:         fakeUrl,
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: response{
//				status: http.StatusOK,
//				body:   toJSON(testdata.FakeNewVideos),
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := http.Get(tt.req.url)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				data, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Errorf("read body: %v", err)
//					return
//				}
//				assert.JSONEq(
//					t,
//					strings.TrimSpace(string(data)),
//					strings.TrimSpace(string(tt.want.body)),
//					"they should be equal",
//				)
//			}
//		})
//	}
//}
//
//func Test_RestApi_Get_Video(t *testing.T) {
//	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
//	if err != nil {
//		t.Errorf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	fakeExistTitle := testdata.FakeVideos[0].Title
//	fakeDoesNotExistTitle := "fakeDoesNotExistTitle"
//	fakeExistVideo := testdata.FakeNewVideos[0]
//	fakeUrl := func(title string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
//	}
//	type request struct {
//		url         string
//		contentType string
//	}
//	type response struct {
//		status int
//		body   interface{}
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    response
//		wantErr bool
//	}{
//		{
//			name: "When title doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeDoesNotExistTitle),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: response{
//				status: http.StatusNotFound,
//				body:   []byte("Not Found"),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When title exists",
//			req: request{
//				url:         fakeUrl(fakeExistTitle),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: response{
//				status: http.StatusOK,
//				body:   fakeExistVideo,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := http.Get(tt.req.url)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				data, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Errorf("read body: %v", err)
//					return
//				}
//				if tt.name == "When title doesn't exist" {
//					assert.Equal(
//						t,
//						strings.TrimSpace(string(data)),
//						strings.TrimSpace(string(tt.want.body.([]byte))),
//						"they should be equal",
//					)
//					return
//				}
//				if tt.name == "When title exists" {
//					videoBody := video.Video{}
//					if err := json.Unmarshal(data, &videoBody); err != nil {
//						t.Errorf("unmarshal data: %v", err)
//						return
//					}
//					assert.ObjectsAreEqualValues(videoBody, fakeExistVideo)
//					return
//				}
//			}
//		})
//	}
//}
//
//func Test_RestApi_Delete_Video(t *testing.T) {
//	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
//	if err != nil {
//		t.Errorf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	fakeExistTitle := testdata.FakeVideos[0].Title
//	fakeDoesNotExistTitle := "doesNotExistTitle"
//	fakeUrl := func(title string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
//	}
//	type request struct {
//		url         string
//		contentType string
//	}
//	type response struct {
//		status int
//		body   string
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    response
//		wantErr bool
//	}{
//		{
//			name: "When title doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeDoesNotExistTitle),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: response{
//				status: http.StatusNotFound,
//				body:   "Not Found",
//			},
//			wantErr: false,
//		},
//		{
//			name: "When title exists",
//			req: request{
//				url:         fakeUrl(fakeExistTitle),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: response{
//				status: http.StatusOK,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			client := &http.Client{}
//			req, err := http.NewRequest(http.MethodDelete, tt.req.url, nil)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			got, err := client.Do(req)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				bs, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Errorf("read body: %v", err)
//					return
//				}
//				data := strings.TrimSpace(string(bs))
//				if data != tt.want.body {
//					t.Errorf("\nbody: %v\nwant: %v", data, tt.want.body)
//				}
//			}
//		})
//	}
//}
//
//func Test_RestApi_Update_Video(t *testing.T) {
//	cfg, teardownTestCase, err := setupTestCase(t, testdata.FakeVideos)
//	if err != nil {
//		t.Errorf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	const (
//		fakeVideosIndex   = 0
//		fakeCategoryIndex = 0
//		fakeGenreIndex    = 0
//	)
//	fakeExistentVideo := testdata.FakeVideos[fakeVideosIndex]
//	fakeExistTitle := fakeExistentVideo.Title
//	fakeDoesNotExistTitle := "doesNotExistTitle"
//
//	fakeUrl := func(title string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "videos", title)
//	}
//	fakeExistentCategoryId := toJSON([]string{
//		fakeExistentVideo.CategoriesId[fakeCategoryIndex],
//	})
//	fakeDoesNotExistentCategory := toJSON([]string{
//		testdata.FakeNonExistentCategoryId,
//	})
//	fakeExistentGenreId := toJSON([]string{
//		fakeExistentVideo.GenresId[fakeGenreIndex],
//	})
//	fakeDoesNotExistentGenre := toJSON([]string{
//		testdata.FakeNonExistentGenreId,
//	})
//	type request struct {
//		url         string
//		contentType string
//		body        io.Reader
//	}
//	type response struct {
//		status int
//		body   string
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    response
//		wantErr bool
//	}{
//		{
//			name: "When title doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeDoesNotExistTitle),
//				contentType: "application/json; charset=UTF-8",
//				body: strings.NewReader(
//					fmt.Sprintf(
//						`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d, "genres":%s, "categories":%s}`,
//						faker.Name(),
//						faker.Sentence(),
//						2020,
//						false,
//						video.TenRating,
//						250,
//						fakeExistentGenreId,
//						fakeExistentCategoryId,
//					)),
//			},
//			want: response{
//				status: http.StatusNotFound,
//				body:   "Not Found",
//			},
//			wantErr: false,
//		},
//		{
//			name: "When Video is with wrong categories and genres",
//			req: request{
//				fakeUrl(fakeExistTitle),
//				"application/json; charset=UTF-8",
//				strings.NewReader(
//					fmt.Sprintf(
//						`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d, "genres":%s, "categories":%s}`,
//						faker.Name(),
//						faker.Sentence(),
//						2020,
//						false,
//						video.TenRating,
//						250,
//						fakeDoesNotExistentCategory,
//						fakeDoesNotExistentGenre,
//					)),
//			},
//			want: response{
//				status: http.StatusNotFound,
//				body:   http.StatusText(http.StatusNotFound),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When everything is right",
//			req: request{
//				fakeUrl(fakeExistTitle),
//				"application/json; charset=UTF-8",
//				strings.NewReader(
//					fmt.Sprintf(
//						`{"title": "%s", "description": "%s", "year_launched": %d, "opened": %t, "rating": %d, "duration": %d, "genres":%s, "categories":%s}`,
//						faker.Name(),
//						faker.Sentence(),
//						2020,
//						false,
//						video.TenRating,
//						250,
//						fakeExistentGenreId,
//						fakeExistentCategoryId,
//					)),
//			},
//			want: response{
//				status: http.StatusOK,
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			client := &http.Client{}
//			req, err := http.NewRequest(http.MethodPut, tt.req.url, tt.req.body)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			got, err := client.Do(req)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Errorf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				bs, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Errorf("read body: %v", err)
//					return
//				}
//				data := strings.TrimSpace(string(bs))
//				if data != tt.want.body {
//					t.Errorf("\nbody: %v\nwant: %v", data, tt.want.body)
//				}
//			}
//		})
//	}
//}
//
//func postMultipartWithFile(
//	uri string,
//	fields map[string]string,
//	fileSeed *memory.FileSeed,
//	fieldNameToFile,
//	pathFile string,
//) (*http.Response, error) {
//	var buffer bytes.Buffer
//	multipartWriter := multipart.NewWriter(&buffer)
//	for fieldName, value := range fields {
//		if err := multipartWriter.WriteField(fieldName, value); err != nil {
//			return nil, fmt.Errorf("could not write field: %v\n", err)
//		}
//	}
//	fileWriter, err := multipartWriter.CreateFormFile(fieldNameToFile, pathFile)
//	if err != nil {
//		return nil, err
//	}
//	file, err := fileSeed.FakeAfero.ReadFile(pathFile)
//	if err != nil {
//		return nil, fmt.Errorf("could not read file: %v\n", err)
//	}
//	if _, err = fileWriter.Write(file); err != nil {
//		return nil, fmt.Errorf("could not copy to file writer: %v\n", err)
//	}
//	if err := multipartWriter.Close(); err != nil {
//		log.Printf("could not close multipartWriter: %v\n", err)
//	}
//	req, err := http.NewRequest(http.MethodPost, uri, &buffer)
//	if err != nil {
//		return nil, fmt.Errorf("could not new request: %v\n", err)
//	}
//	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
//	client := &http.Client{}
//	res, err := client.Do(req)
//	return res, err
//}
