package app

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"main/internal/constants"
	"main/internal/mocks"
	"main/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExampleHealthHandlers_Ping shows how to use the Ping method from HealthHandlers.
func ExampleHealthHandlers_Ping() {
	// Setup controller and mock objects.
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	// Create a mock HealthService.
	mockHealthService := mocks.NewMockHealthService(ctrl)

	// Define expectation: expect one call to Ping returning no error.
	mockHealthService.EXPECT().Ping().Return(nil)

	// Create a new HealthHandlers instance using the mock service.
	handlers := NewHealthHandlers(mockHealthService)

	// Create a simple HTTP request.
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	// Record the response using httptest.
	recorder := httptest.NewRecorder()

	// Execute the Ping method.
	handlers.Ping(recorder, req)

	// Access the response details properly via Header().
	fmt.Printf("Response Status: %v\n", recorder.Code)
	fmt.Printf("Content Type: %v\n", recorder.Header().Get("Content-Type"))

	// Output:
	// Response Status: 200
	// Content Type: text/plain
}

func TestPing(t *testing.T) {
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
				statusCode:  http.StatusOK,
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.req.method, "/ping", nil)
			w := httptest.NewRecorder()

			r, _ := NewRepository(c)
			l := services.NewHealthService(r.health)
			h := NewHealthHandlers(l)

			h.Ping(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			res.Body.Close()
		})
	}
}
