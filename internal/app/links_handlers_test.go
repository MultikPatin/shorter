package app

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"main/internal/config"
	"main/internal/constants"
	"main/internal/mocks"
	"main/internal/models"
	"main/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExampleLinksHandlers_AddLink demonstrates how to use the AddLink method to add a single link.
func ExampleLinksHandlers_AddLink() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	// Define expectation: expect one call to Add with specific parameters.
	linkToCreate := models.OriginLink{
		URL: "https://example.com",
	}
	resultURL := "http://localhost/test-short-url"
	mockLinksService.EXPECT().Add(gomock.Any(), linkToCreate, "localhost").Return(resultURL, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	reqBody := `{"url": "https://example.com"}`
	request, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(reqBody))

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the AddLink method.
	handlers.AddLink(recorder, request)

	// Output:
	// Response Status: 201 Created
	// Content Type: application/json
	// {"result":"http://localhost/test-short-url"}
}

// ExampleLinksHandlers_AddLinks demonstrates how to use the AddLinks method to add multiple links at once.
func ExampleLinksHandlers_AddLinks() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	// Define expectation: expect one call to AddBatch with specific parameters.
	batchInput := []models.OriginLink{
		{URL: "https://example1.com"},
		{URL: "https://example2.com"},
	}
	batchResult := []string{"http://localhost/link1", "http://localhost/link2"}
	mockLinksService.EXPECT().AddBatch(gomock.Any(), batchInput, "localhost").Return(batchResult, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	reqBody := `[{"url": "https://example1.com"}, {"url": "https://example2.com"}]`
	request, _ := http.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(reqBody))

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the AddLinks method.
	handlers.AddLinks(recorder, request)

	// Output:
	// Response Status: 201 Created
	// Content Type: application/json
	// ["http://localhost/link1","http://localhost/link2"]
}

// ExampleGetLink demonstrates how to resolve a short link to its original URL.
func ExampleLinksHandlers_GetLink() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	// Define expectation: expect one call to Get with specific ID.
	id := "abc123"
	originalURL := "https://example.com"
	mockLinksService.EXPECT().Get(gomock.Any(), id).Return(originalURL, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	request, _ := http.NewRequest(http.MethodGet, "/"+id, nil)

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the GetLink method.
	handlers.GetLink(recorder, request)

	// Output:
	// Response Status: 307 Temporary Redirect
	// Location: https://example.com
}

// ExampleLinksHandlers_AddLinkInText demonstrates how to create a link from raw text input.
func ExampleLinksHandlers_AddLinkInText() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	// Define expectation: expect one call to Add with specific parameters.
	linkToCreate := models.OriginLink{
		URL: "https://example.com",
	}
	resultURL := "http://localhost/test-short-url"
	mockLinksService.EXPECT().Add(gomock.Any(), linkToCreate, "localhost").Return(resultURL, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	reqBody := "https://example.com"
	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the AddLinkInText method.
	handlers.AddLinkInText(recorder, request)

	// Output:
	// Response Status: 201 Created
	// Content Type: text/plain
	// http://localhost/test-short-url
}

var c = &config.Config{
	StorageFilePaths: "test.json",
}

func TestAddLinkInText(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	type req struct {
		method string
	}
	tests := []struct {
		name string
		want want
		req  req
	}{
		{
			name: "positive case",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusCreated,
			},
			req: req{
				method: http.MethodPost,
			},
		},
		{
			name: "wrong method",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusMethodNotAllowed,
			},
			req: req{
				method: http.MethodGet,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.req.method, "/", nil)
			w := httptest.NewRecorder()

			r, _ := NewRepository(c)
			l := services.NewLinksService(c, r.links)
			h := NewLinksHandlers(l)

			h.AddLinkInText(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			res.Body.Close()
		})
	}
}

func TestAddLink(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	type req struct {
		method string
	}
	tests := []struct {
		name string
		want want
		req  req
		body string
	}{
		{
			name: "positive case",
			want: want{
				contentType: constants.JSONContentType,
				statusCode:  http.StatusCreated,
			},
			req: req{
				method: http.MethodPost,
			},
			body: `{"url":"https://go.dev/blog/package-names"}`,
		},
		{
			name: "wrong method",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusMethodNotAllowed,
			},
			req: req{
				method: http.MethodGet,
			},
			body: ``,
		},
		{
			name: "wrong body",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusBadRequest,
			},
			req: req{
				method: http.MethodPost,
			},
			body: ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString(test.body)

			request := httptest.NewRequest(test.req.method, "/api/shorten", &buf)
			w := httptest.NewRecorder()

			r, _ := NewRepository(c)
			l := services.NewLinksService(c, r.links)
			h := NewLinksHandlers(l)

			h.AddLink(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			res.Body.Close()
		})
	}
}

func TestAddLinks(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	type req struct {
		method string
	}
	tests := []struct {
		name string
		want want
		req  req
		body string
	}{
		{
			name: "positive case",
			want: want{
				contentType: constants.JSONContentType,
				statusCode:  http.StatusCreated,
			},
			req: req{
				method: http.MethodPost,
			},
			body: `[
						{"correlation_id": "wf","original_url": "https://go.dev/blog/package-names/1"},
						{"correlation_id": "wf","original_url": "https://go.dev/blog/package-names"}
					]`,
		},
		{
			name: "wrong method",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusMethodNotAllowed,
			},
			req: req{
				method: http.MethodGet,
			},
			body: ``,
		},
		{
			name: "wrong body",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusBadRequest,
			},
			req: req{
				method: http.MethodPost,
			},
			body: ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString(test.body)

			request := httptest.NewRequest(test.req.method, "/api/shorten/batch", &buf)
			w := httptest.NewRecorder()

			r, _ := NewRepository(c)
			l := services.NewLinksService(c, r.links)
			h := NewLinksHandlers(l)

			h.AddLinks(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			res.Body.Close()
		})
	}
}

func TestGetLink(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	type req struct {
		method string
	}
	tests := []struct {
		name string
		want want
		req  req
	}{
		{
			name: "positive case",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusTemporaryRedirect,
			},
			req: req{
				method: http.MethodGet,
			},
		},
		{
			name: "wrong method",
			want: want{
				contentType: constants.TextContentType,
				statusCode:  http.StatusMethodNotAllowed,
			},
			req: req{
				method: http.MethodPost,
			},
		},
	}
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			u, err := uuid.NewRandom()
			if err != nil {
				t.Fatalf("Failed to generate UUID")
				return
			}

			r, _ := NewRepository(c)
			l := services.NewLinksService(c, r.links)
			h := NewLinksHandlers(l)

			addedLink := models.AddedLink{
				Short:  u.String(),
				Origin: "test.com",
			}

			id, err := r.links.Add(ctx, addedLink)
			if err != nil {
				t.Fatalf("Failed to add link")
				return
			}

			request := httptest.NewRequest(test.req.method, "/"+id+"/", nil)
			request.SetPathValue("id", id)

			w := httptest.NewRecorder()
			h.GetLink(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			err = res.Body.Close()
			if err != nil {
				t.Fatalf("Failed to close response body")
				return
			}
		})
	}
}
