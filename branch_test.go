package GitLabApiExperiment

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetBranches(t *testing.T) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// close the server at the end of the function
	defer server.Close()

	// Add an url and its handler to the test HTTP server.
	mux.HandleFunc(
		"/api/v3/projects/test_project/repository/branches",
		func(w http.ResponseWriter, r *http.Request) {

			method := r.Method
			if method != "GET" {
				t.Errorf("Request method: %s, want %s", method, "GET")
			}

			fmt.Fprint(w, `
				[
				  {
				    "name": "master"
				  }
				]`)
		})

	expectedBranches := []Branch{{Name: "master"}}

	returnedBranches, err := GetBranches("gitlab_token", server.URL, "test_project")
	if err != nil {
		t.Errorf("getBranches returned error: %v", err)
	}

	if reflect.DeepEqual(expectedBranches, returnedBranches) == false {
		t.Errorf("getBranches returned %+v, expected %+v", returnedBranches, expectedBranches)
	}
}
