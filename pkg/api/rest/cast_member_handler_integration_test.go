// +build integration

package rest_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"

	"github.com/selmison/code-micro-videos/pkg/api/rest"
	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/storage/test"
	"github.com/selmison/code-micro-videos/testdata"
)

var fakeCastMembers = testdata.FakeCastMembers

func Test_RestApi_Post_CastMembers(t *testing.T) {
	teardownTestCase, _, err := test.SetupTestCase(t, nil)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	srv := httptest.NewServer(rest.NewServer())
	defer srv.Close()
	fakeURL := fmt.Sprintf("%s/%s", srv.URL, "cast_members")
	fakeNewCastMemberDTO := cast_member.NewCastMemberDTO{
		Name: faker.Name(),
		Type: cast_member.Director,
	}
	fakeNewCastMemberJSONWithNameBlank := fmt.Sprintf(
		`{"name": "%v", "type": %d}`,
		"     ",
		fakeNewCastMemberDTO.Type,
	)
	fakeNewCastMemberJSON := fmt.Sprintf(
		`{"name": "%v", "type": %d}`,
		fakeNewCastMemberDTO.Name,
		fakeNewCastMemberDTO.Type,
	)
	type request struct {
		url  string
		body io.Reader
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
			name: "When name field is blank",
			req: request{
				fakeURL,
				strings.NewReader(fakeNewCastMemberJSONWithNameBlank),
			},
			want: response{
				status: http.StatusBadRequest,
				body:   http.StatusText(http.StatusBadRequest),
			},
			wantErr: false,
		},
		{
			name: "When everything is right",
			req: request{
				fakeURL,
				strings.NewReader(fakeNewCastMemberJSON),
			},
			want: response{
				status: http.StatusCreated,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Post(tt.req.url, "application/json; charset=UTF-8", tt.req.body)
			if (err != nil) != tt.wantErr {
				require.Error(t, err)
			}
			if got != nil {
				require.Equal(t, tt.want.status, got.StatusCode, "they should be equal")
				switch tt.want.status {
				case http.StatusCreated:
					createResponse := cast_member.CreateResponse{}
					defer got.Body.Close()
					if err := json.NewDecoder(got.Body).Decode(&createResponse); err != nil {
						t.Fatalf("could not decode response: %v", err)
					}
					require.Equal(t, createResponse.CastMember, createResponse.CastMember, "they should be equal")
				}
			}
		})
	}
}

func Test_RestApi_Get_CastMembers(t *testing.T) {
	teardownTestCase, _, err := test.SetupTestCase(t, fakeCastMembers)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	srv := httptest.NewServer(rest.NewServer())
	defer srv.Close()
	fakeUrl := fmt.Sprintf("%s/%s", srv.URL, "cast_members")
	type response struct {
		status int
		body   []byte
	}
	tests := []struct {
		name    string
		url     string
		want    response
		wantErr bool
	}{
		{
			name: "When everything is right",
			url:  fakeUrl,
			want: response{
				status: http.StatusOK,
				body:   toJSON(testdata.FakeNewCastMembers),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := http.Get(tt.url)
			if (err != nil) != tt.wantErr {
				require.Error(t, err)
			}
			if got != nil {
				require.Equal(t, tt.want.status, got.StatusCode, "they should be equal")
				switch tt.want.status {
				case http.StatusOK:
					listResponse := cast_member.ListResponse{}
					defer got.Body.Close()
					if err := json.NewDecoder(got.Body).Decode(&listResponse); err != nil {
						t.Fatalf("could not decode response: %v", err)
					}
					require.ElementsMatch(
						t,
						listResponse.CastMembers,
						testdata.FakeCastMemberDTOs,
						"they should be equal",
					)
				}
			}
		})
	}
}

//func Test_RestApi_Get_CastMember(t *testing.T) {
//	teardownTestCase, _, err := test.SetupTestCase(t, testdata.FakeCastMembers)
//	if err != nil {
//		t.Fatalf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	fakeExistentId := testdata.FakeCastMembers[0].Id()
//	fakeDoesNonExistentId := "fakeDoesNonExistentId"
//	fakeExistCastMemberDTO := testdata.FakeNewCastMembers[0]
//	fakeUrl := func(name string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "cast_members", name)
//	}
//	type request struct {
//		url         string
//		contentType string
//	}
//	type statusResponse struct {
//		status int
//		body   []byte
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    statusResponse
//		wantErr bool
//	}{
//		{
//			name: "When name doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeDoesNonExistentId),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: statusResponse{
//				status: http.StatusNotFound,
//				body:   toJSON("Not Found"),
//			},
//			wantErr: false,
//		},
//		{
//			name: "When name exists",
//			req: request{
//				url:         fakeUrl(fakeExistentId),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: statusResponse{
//				status: http.StatusOK,
//				body:   toJSON(fakeExistCastMemberDTO),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := http.Get(tt.req.url)
//			if (err != nil) != tt.wantErr {
//				t.Fatalf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Fatalf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				body, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Fatalf("read body: %v", err)
//					return
//				}
//				bodyResponse, err := json.Marshal(body)
//				if err != nil {
//					t.Fatalf("read body: %v", err)
//					return
//				}
//				comp, err := JSONBytesEqual(bodyResponse, tt.want.body)
//				if err != nil {
//					t.Fatalf("read body: %v", err)
//					return
//				}
//				if comp {
//					t.Fatalf("\nresponse: %v\nwant: %v", string(bodyResponse), string(tt.want.body))
//				}
//			}
//		})
//	}
//}
//
//func Test_RestApi_Delete_CastMember(t *testing.T) {
//	teardownTestCase, _, err := test.SetupTestCase(t, testdata.FakeCastMembers)
//	if err != nil {
//		t.Fatalf("test: failed to setup test case: %v\n", err)
//	}
//	defer teardownTestCase(t)
//	fakeExistentId := testdata.FakeCastMembers[0].Id()
//	fakeNonExistentId := "fakeNonExistentId"
//	fakeUrl := func(name string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "cast_members", name)
//	}
//	type request struct {
//		url         string
//		contentType string
//	}
//	type statusResponse struct {
//		status int
//		body   string
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    statusResponse
//		wantErr bool
//	}{
//		{
//			name: "When name doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeNonExistentId),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: statusResponse{
//				status: http.StatusNotFound,
//				body:   "Not Found",
//			},
//			wantErr: false,
//		},
//		{
//			name: "When name exists",
//			req: request{
//				url:         fakeUrl(fakeExistentId),
//				contentType: "application/json; charset=UTF-8",
//			},
//			want: statusResponse{
//				status: http.StatusOK,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			client := &http.Client{}
//			req, err := http.NewRequest(http.MethodDelete, tt.req.url, nil)
//			if (err != nil) != tt.wantErr {
//				t.Fatalf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			got, err := client.Do(req)
//			if (err != nil) != tt.wantErr {
//				t.Fatalf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Fatalf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				bs, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Fatalf("read body: %v", err)
//					return
//				}
//				data := strings.TrimSpace(string(bs))
//				if data != tt.want.body {
//					t.Fatalf("\nbody: %v\nwant: %v", data, tt.want.body)
//				}
//			}
//		})
//	}
//}
//
//func Test_RestApi_Update_CastMember(t *testing.T) {
//	teardownTestCase, _, err := test.SetupTestCase(t, testdata.FakeCastMembers)
//	if err != nil {
//		t.Fatalf("test: failed to setup test case: %v\n", err)
//		return
//	}
//	defer teardownTestCase(t)
//	fakeExistentId := testdata.FakeCastMembers[0].Id()
//	fakeNonExistentId := "fakeNonExistentId"
//	fakeUrl := func(name string) string {
//		return fmt.Sprintf("http://%s/%s/%s", cfg.AddressServer, "cast_members", name)
//	}
//	type request struct {
//		url         string
//		contentType string
//		body        io.Reader
//	}
//	type statusResponse struct {
//		status int
//		body   string
//	}
//	tests := []struct {
//		name    string
//		req     request
//		want    statusResponse
//		wantErr bool
//	}{
//		{
//			name: "When name doesn't exist",
//			req: request{
//				url:         fakeUrl(fakeNonExistentId),
//				contentType: "application/json; charset=UTF-8",
//				body: strings.NewReader(fmt.Sprintf(
//					`{"name": "%s"}`,
//					faker.Name(),
//				)),
//			},
//			want: statusResponse{
//				status: http.StatusNotFound,
//				body:   "Not Found",
//			},
//			wantErr: false,
//		},
//		{
//			name: "When name exists",
//			req: request{
//				url:         fakeUrl(fakeExistentId),
//				contentType: "application/json; charset=UTF-8",
//				body: strings.NewReader(fmt.Sprintf(
//					`{"name": "%s", "avatar": "%s", "whatsapp": "%s", "bio": "%s"}`,
//					faker.Name(),
//					faker.URL(),
//					faker.Phonenumber(),
//					faker.Sentence(),
//				)),
//			},
//			want: statusResponse{
//				status: http.StatusOK,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			client := &http.Client{}
//			req, err := http.NewRequest(http.MethodPut, tt.req.url, tt.req.body)
//			if (err != nil) != tt.wantErr {
//				t.Fatalf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			got, err := client.Do(req)
//			if (err != nil) != tt.wantErr {
//				t.Fatalf("error: %v, wantErr: %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if got.StatusCode != tt.want.status {
//					t.Fatalf("statusCode: %v, want: %v", got.StatusCode, tt.want.status)
//					return
//				}
//				bs, err := ioutil.ReadAll(got.Body)
//				if err != nil {
//					t.Fatalf("read body: %v", err)
//					return
//				}
//				data := strings.TrimSpace(string(bs))
//				if data != tt.want.body {
//					t.Fatalf("\nbody: %v\nwant: %v", data, tt.want.body)
//				}
//			}
//		})
//	}
//}

//func newCastMemberDTOToCastMemberDTO(id string, dto cast_member.NewCastMemberDTO) *cast_member.castMemberDTO {
//	return &cast_member.castMemberDTO{
//		Id:   id,
//		Name: dto.Name,
//		Type: dto.Type,
//	}
//}

func toJSON(i interface{}) []byte {
	s, _ := json.Marshal(i)
	return s
}
