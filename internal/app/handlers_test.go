package app

import (
	"github.com/google/uuid"
	"main/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostLink(t *testing.T) {
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
				contentType: contentType,
				statusCode:  http.StatusCreated,
			},
			req: req{
				method: http.MethodPost,
			},
		},
		{
			name: "wrong method",
			want: want{
				contentType: contentType,
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

			d := db.NewInMemoryDB()
			h := GetHeaders(d)

			h.postLink(w, request)

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
				contentType: contentType,
				statusCode:  http.StatusTemporaryRedirect,
			},
			req: req{
				method: http.MethodGet,
			},
		},
		{
			name: "wrong method",
			want: want{
				contentType: contentType,
				statusCode:  http.StatusMethodNotAllowed,
			},
			req: req{
				method: http.MethodPost,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			u, err := uuid.NewRandom()
			if err != nil {
				t.Fatalf("Failed to generate UUID")
				return
			}

			d := db.NewInMemoryDB()
			h := GetHeaders(d)

			id, err := d.AddLink(u.String(), urlPrefix+"test.com")
			if err != nil {
				t.Fatalf("Failed to add link")
				return
			}

			request := httptest.NewRequest(test.req.method, "/"+id+"/", nil)
			request.SetPathValue("id", id)

			w := httptest.NewRecorder()
			h.getLink(w, request)

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
