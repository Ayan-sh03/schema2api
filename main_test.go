package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper to create a sample schema for testing.
func createSampleSchema() *Schema {
	return &Schema{
		Title: "User",
		Type:  "object",
		Properties: map[string]Property{
			"id":    {Type: "integer"},
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
		Required: []string{"id", "name", "email"},
	}
}

// Helper to perform a request and check the response.
func performRequest(t *testing.T, handler http.HandlerFunc, method, path string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func TestUploadHandler(t *testing.T) {
	// Reset schema before tests
	currentSchema = nil

	t.Run("Successful Upload", func(t *testing.T) {
		schema := createSampleSchema()
		schemaJSON, _ := json.Marshal(schema)
		rr := performRequest(t, uploadHandler, http.MethodPost, "/upload", schemaJSON)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"message":"Schema uploaded successfully","title":"User"}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
		if currentSchema == nil || currentSchema.Title != "User" {
			t.Errorf("currentSchema was not updated correctly")
		}
	})

	t.Run("Invalid Method", func(t *testing.T) {
		rr := performRequest(t, uploadHandler, http.MethodGet, "/upload", nil)
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		rr := performRequest(t, uploadHandler, http.MethodPost, "/upload", []byte("{invalid json"))
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestCatchAllHandler(t *testing.T) {
	// Reset schema before tests
	currentSchema = nil

	t.Run("No Schema Loaded", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodGet, "/users", nil)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		expected := "No schema uploaded. Please POST your JSON schema to /upload"
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	// Load schema for subsequent tests
	currentSchema = createSampleSchema()
	entityPlural := "users" // Based on schema title "User"

	t.Run("GET List", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodGet, "/"+entityPlural, nil)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		// Check if it's a JSON array
		if !strings.HasPrefix(rr.Body.String(), "[") {
			t.Errorf("handler returned non-array body for list: got %v", rr.Body.String())
		}
	})

	t.Run("GET Single", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodGet, "/"+entityPlural+"/123", nil)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		// Check if it's a JSON object and contains the ID
		if !strings.HasPrefix(rr.Body.String(), "{") || !strings.Contains(rr.Body.String(), `"id":123`) {
			t.Errorf("handler returned unexpected body for single item: got %v", rr.Body.String())
		}
	})

	t.Run("GET Invalid ID", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodGet, "/"+entityPlural+"/abc", nil)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("GET Non-existent Entity", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodGet, "/products", nil)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("POST", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodPost, "/"+entityPlural, []byte(`{"name":"test"}`)) // Body content doesn't matter for mock
		if status := rr.Code; status != http.StatusOK { // Should be 201 Created ideally, but OK for mock
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if !strings.HasPrefix(rr.Body.String(), "{") || !strings.Contains(rr.Body.String(), `"id":1`) {
			t.Errorf("handler returned unexpected body for POST: got %v", rr.Body.String())
		}
	})

	t.Run("PUT", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodPut, "/"+entityPlural+"/456", []byte(`{"name":"updated"}`)) // Body content doesn't matter
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if !strings.HasPrefix(rr.Body.String(), "{") || !strings.Contains(rr.Body.String(), `"id":456`) {
			t.Errorf("handler returned unexpected body for PUT: got %v", rr.Body.String())
		}
	})

	t.Run("PUT Invalid ID", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodPut, "/"+entityPlural+"/abc", nil)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodDelete, "/"+entityPlural+"/789", nil)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := `{"message":"Deleted successfully"}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body for DELETE: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("DELETE Invalid ID", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodDelete, "/"+entityPlural+"/abc", nil)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Unsupported Method", func(t *testing.T) {
		rr := performRequest(t, catchAllHandler, http.MethodPatch, "/"+entityPlural+"/1", nil)
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})
}