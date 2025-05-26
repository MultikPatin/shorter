package app

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"main/internal/mocks"
	"main/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
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
	// Content Type: text/plain; charset=utf-8

}

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
	fmt.Printf("Response body: %v\n", recorder.Body.String())

	// Output:
	// Response Status: 201
	// Content Type: application/json
	// Response body: {"result":"http://localhost:8080/test-short-url"}

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
	fmt.Printf("Response body: %v\n", recorder.Body.String())

	// Output:
	// Response Status: 201
	// Content Type: application/json
	// Response body: [{"short_url":"http://localhost/link1"},{"short_url":"http://localhost/link2"}]
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
