package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"nbaTeamInfo", "/nba/team-info-summary", "GET", []postData{}, http.StatusOK},
	{"nbaResults", "/nba/results", "GET", []postData{}, http.StatusOK},

	{"post-nbaTeams", "/nba/teams", "POST", []postData{
		{key: "TeamName", value: "Jannseni Krakenid"},
		{key: "Abbreviation", value: "JNSN"},
		{key: "TeamColor", value: "#FFFFFF"},
		{key: "DarkText", value: "true"},
	}, http.StatusOK},
	{"post-nbaResults", "/nba/results", "POST", []postData{
		{key: "home_team", value: "Jannseni Krakenid"},
		{key: "home_score", value: "21"},
		{key: "away_score", value: "13"},
		{key: "away_team", value: "Hiiu Kalur"},
	}, http.StatusOK},
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
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			response, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
