package app

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"main/internal/adapters"
	"main/internal/config"
	"main/internal/constants"
	"main/internal/models"
	"main/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c = &config.Config{
	StorageFilePaths: "test.json",
}
var logger = adapters.GetLogger()

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
