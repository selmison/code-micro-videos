package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/bxcodec/faker/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	var code int
	dbContainer, err := testdata.NewDBContainer(context.Background())
	if err != nil {
		log.Fatal(err)
		return 1
	}
	dbConnStr := dbContainer.ConnStr
	config.DBConnStr = dbConnStr
	if err := config.InitDB(dbConnStr); err != nil {
		log.Fatalln(err, "init db")
		return 1
	}
	defer func() {
		if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
			code = 1
			log.Fatalln(err)
		}
	}()
	if code > 0 {
		return code
	}
	go func() {
		if err := InitApp(context.Background(), dbConnStr); err != nil {
			code = 1
			log.Fatalln(err, "init app")
		}
	}()
	if code > 0 {
		return code
	}
	time.Sleep(1 * time.Second)
	code = m.Run()
	return code
}

func Test_integration_CategoryCreate(t *testing.T) {
	fakeCategory := `{"name": "action", "description": "actions films"}`
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
			name: "create a category",
			req: request{
				fmt.Sprintf("http://%s/%s", config.AddressServer, "categories"),
				"application/json; charset=UTF-8",
				strings.NewReader(fakeCategory),
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
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
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
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func Test_RestApi_Post_Categories(t *testing.T) {
	fakeUrl := fmt.Sprintf("http://%s/%s", config.AddressServer, "categories")
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
					`{"name": "%s", "avatar": "%s", "whatsapp": "%s", "bio": "%s" }`,
					"",
					faker.URL(),
					faker.Phonenumber(),
					faker.Sentence(),
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
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func Test_RestApi_Get_Categories(t *testing.T) {
	fakeUrl := fmt.Sprintf("http://%s/%s", config.AddressServer, "categories")
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
	db, err := sql.Open(config.DBDrive, config.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	for _, c := range testdata.FakeCategories {
		err = c.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
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
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("\nresponse: %v\nwant: %v", data, tt.want.body)
				}
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func Test_RestApi_Get_Category(t *testing.T) {
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[0]
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", config.AddressServer, "categories", name)
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
				body:   toJSON(fakeExistCategoryDTO),
			},
		},
	}
	db, err := sql.Open(config.DBDrive, config.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	for _, c := range testdata.FakeCategories {
		err = c.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
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
				data := strings.TrimSpace(string(bs))
				if data != tt.want.body {
					t.Errorf("\nbody: %v\nwant: %v", data, tt.want.body)
				}
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func Test_RestApi_Delete_Category(t *testing.T) {
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", config.AddressServer, "categories", name)
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
	db, err := sql.Open(config.DBDrive, config.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	for _, c := range testdata.FakeCategories {
		err = c.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
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
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func Test_RestApi_Update_Category(t *testing.T) {
	fakeExistName := testdata.FakeCategories[0].Name
	fakeDoesNotExistName := "doesNotExistName"
	fakeUrl := func(name string) string {
		return fmt.Sprintf("http://%s/%s/%s", config.AddressServer, "categories", name)
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
					`{"name": "%s", "avatar": "%s", "whatsapp": "%s", "bio": "%s" }`,
					faker.Name(),
					faker.URL(),
					faker.Phonenumber(),
					faker.Sentence(),
				)),
			},
			want: response{
				status: http.StatusOK,
			},
		},
	}
	db, err := sql.Open(config.DBDrive, config.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	for _, c := range testdata.FakeCategories {
		err = c.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
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
	if err := config.ClearCategoriesTable(config.DBDrive, config.DBConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func toJSON(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
