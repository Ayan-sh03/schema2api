package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Schema defines the JSON schema structure.
type Schema struct {
	Title      string              `json:"title"`
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

// Property defines each property's type.
type Property struct {
	Type string `json:"type"`
}

// currentSchema holds the uploaded JSON schema.
var currentSchema *Schema

// dummyData generates a dummy data object based on the schema.
func dummyData() map[string]interface{} {
	data := make(map[string]interface{})
	if currentSchema == nil {
		return data
	}
	for key, prop := range currentSchema.Properties {
		switch prop.Type {
		case "string":
			data[key] = "example"
		case "integer":
			data[key] = 1
		case "number":
			data[key] = 0.0
		case "boolean":
			data[key] = false
		default:
			data[key] = nil
		}
	}
	return data
}

// uploadHandler handles uploading and parsing JSON schema.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var schema Schema
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		http.Error(w, "Invalid JSON schema: "+err.Error(), http.StatusBadRequest)
		return
	}
	currentSchema = &schema
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Schema uploaded successfully",
		"title":   schema.Title,
	}
	json.NewEncoder(w).Encode(response)
}

// catchAllHandler handles all other routes.
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure a schema is loaded.
	if currentSchema == nil {
		http.Error(w, "No schema uploaded. Please POST your JSON schema to /upload", http.StatusBadRequest)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")
	entity := strings.ToLower(currentSchema.Title) + "s" // simple pluralization
	var responseObj interface{}

	switch r.Method {
	case http.MethodGet:
		if len(segments) == 1 && segments[0] == entity {
			// Return a list of dummy objects
			var list []map[string]interface{}
			for i := 1; i <= 3; i++ {
				obj := dummyData()
				obj["id"] = i
				list = append(list, obj)
			}
			responseObj = list
		} else if len(segments) == 2 && segments[0] == entity {
			         // Return single dummy object reflecting the requested ID
			         requestedID := segments[1]
			         obj := dummyData()

			         // Check schema for expected ID type (simple check for "id" property)
			         idProp, hasIntegerId := currentSchema.Properties["id"]
			         isIntegerExpected := hasIntegerId && idProp.Type == "integer"

			         if isIntegerExpected {
			             // Expecting an integer ID
			             id, err := strconv.Atoi(requestedID)
			             if err != nil {
			                 http.Error(w, "Invalid ID format: expected integer", http.StatusBadRequest)
			                 return
			             }
			             obj["id"] = id
			         } else {
			             // Expecting a string ID (or no specific "id" field)
			             // Find the first string property to use as key, or default to "id"
			             stringKey := "id" // Default key
			             foundKey := false
			             for key, prop := range currentSchema.Properties {
			                  // Use explicit "id" if string, or first string property otherwise
			                 if key == "id" && prop.Type == "string" {
			                      stringKey = key
			                      foundKey = true
			                      break
			                 }
			                 if prop.Type == "string" && !foundKey {
			                     stringKey = key
			                     // Don't break, prefer "id" if found later
			                 }
			             }
			              obj[stringKey] = requestedID
			         }
			         responseObj = obj
		} else {
			http.NotFound(w, r)
			return
		}
	case http.MethodPost:
		// Simulate creation and echo back dummy object
		obj := dummyData()
		obj["id"] = 1 // simulate new id
		responseObj = obj
	case http.MethodPut:
		      // Simulate update and return updated dummy object reflecting the ID
		      if len(segments) == 2 && segments[0] == entity {
		          requestedID := segments[1]
		          obj := dummyData()

		           // Check schema for expected ID type
		          idProp, hasIntegerId := currentSchema.Properties["id"]
		          isIntegerExpected := hasIntegerId && idProp.Type == "integer"

		          if isIntegerExpected {
		               // Expecting an integer ID
		              id, err := strconv.Atoi(requestedID)
		              if err != nil {
		                  http.Error(w, "Invalid ID format: expected integer", http.StatusBadRequest)
		                  return
		              }
		              obj["id"] = id
		          } else {
		              // Expecting a string ID
		               stringKey := "id"
		               foundKey := false
		               for key, prop := range currentSchema.Properties {
		                   if key == "id" && prop.Type == "string" {
		                       stringKey = key
		                       foundKey = true
		                       break
		                   }
		                   if prop.Type == "string" && !foundKey {
		                       stringKey = key
		                   }
		               }
		               obj[stringKey] = requestedID
		          }
		          responseObj = obj
		      } else {
		          http.NotFound(w, r)
			return
		}
	case http.MethodDelete:
		// Simulate deletion by returning a success message.
		if len(segments) == 2 && segments[0] == entity {
			// Validate ID format based on schema expectation
			requestedID := segments[1]
			idProp, hasIntegerId := currentSchema.Properties["id"]
			isIntegerExpected := hasIntegerId && idProp.Type == "integer"

			if isIntegerExpected {
			     // Expecting an integer ID
			    _, err := strconv.Atoi(requestedID)
			    if err != nil {
			        http.Error(w, "Invalid ID format: expected integer", http.StatusBadRequest)
			        return
			    }
			}
			// If not expecting integer, any string is considered valid for DELETE

			// Validation passed
			// In a real scenario, might check against schema type here

			responseObj = map[string]string{"message": "Deleted successfully"}
		} else {
			http.NotFound(w, r)
			return
		}
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseObj); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func main() {
	// Endpoint to upload JSON schema.
	http.HandleFunc("/upload", uploadHandler)
	// Catch-all route handler.
	http.HandleFunc("/", catchAllHandler)

	fmt.Println("Server started on port :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
