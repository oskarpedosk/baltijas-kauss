package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"signIn", "/", "GET", []postData{}, http.StatusOK},
	{"nbaHome", "/nba", "GET", []postData{}, http.StatusOK},
	{"nbaPlayers", "/nba/players", "GET", []postData{}, http.StatusOK},
	{"nbaTeams", "/nba/teams", "GET", []postData{}, http.StatusOK},
	{"nbaTeamsJSON", "/nba/teams-json", "GET", []postData{}, http.StatusOK},
	{"nbaResults", "/nba/results", "GET", []postData{}, http.StatusOK},
}


func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			response, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {

		}
	}
}
