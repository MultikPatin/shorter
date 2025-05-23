package app

import (
	"bytes"
	"context"
	"fmt"
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
	resultURL := "http://localhost:8080/test-short-url"
	mockLinksService.EXPECT().Add(gomock.Any(), linkToCreate, gomock.Any()).Return(resultURL, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	reqBody := `{"url": "https://example.com"}`
	request, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(reqBody))

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the AddLink method.
	handlers.AddLink(recorder, request)

	fmt.Printf("Response Status: %v\n", recorder.Code)
	fmt.Printf("Content Type: %v\n", recorder.Header().Get("Content-Type"))

	// Output:
	// Response Status: 201
	// Content Type: application/json

}

// ExampleLinksHandlers_AddLinks demonstrates how to use the AddLinks method to add multiple links at once.
func ExampleLinksHandlers_AddLinks() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	var originLinks []models.OriginLink
	for i := 1; i < 3; {
		url := fmt.Sprintf("https://example%d.com", i)
		originLink := models.OriginLink{
			URL: url,
		}
		originLinks = append(originLinks, originLink)
		i++
	}

	var results []models.Result
	for i := 1; i < 3; {
		url := fmt.Sprintf("http://localhost/link%d", i)
		result := models.Result{
			Result: url,
		}
		results = append(results, result)
		i++
	}

	mockLinksService.EXPECT().AddBatch(gomock.Any(), originLinks, gomock.Any()).Return(results, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	reqBody := `[{"original_url": "https://example1.com"},{"original_url": "https://example2.com"}]`
	request, _ := http.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(reqBody))

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the AddLinks method.
	handlers.AddLinks(recorder, request)

	fmt.Printf("Response Status: %v\n", recorder.Code)
	fmt.Printf("Content Type: %v\n", recorder.Header().Get("Content-Type"))

	// Output:
	// Response Status: 201
	// Content Type: application/json
}

// ExampleGetLink demonstrates how to resolve a short link to its original URL.
func ExampleLinksHandlers_GetLink() {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock LinksService.
	mockLinksService := mocks.NewMockLinksService(ctrl)

	// Define expectation: expect one call to Get with specific ID.
	id := "5a8afeee-412d-4fbe-a059-79e4dc737e1d"
	originalURL := "https://example.com"
	mockLinksService.EXPECT().Get(gomock.Any(), id).Return(originalURL, nil)

	// Create a new LinksHandlers instance using the mock service.
	handlers := NewLinksHandlers(mockLinksService)

	// Create a simple HTTP request.
	request, _ := http.NewRequest(http.MethodGet, "/"+id+"/", nil)
	request.SetPathValue("id", id)

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the GetLink method.
	handlers.GetLink(recorder, request)

	fmt.Printf("Response Status: %v\n", recorder.Code)
	fmt.Printf("Content Type: %v\n", recorder.Header().Get("Content-Type"))

	// Output:
	// Response Status: 307
	// Content Type: text/plain; charset=utf-8

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
