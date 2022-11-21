package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/driver"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

var positions = []models.NBAPosition{
	{
		Name:   "PG",
		Number: 1,
	},
	{
		Name:   "SG",
		Number: 2,
	},
	{
		Name:   "SF",
		Number: 3,
	},
	{
		Name:   "PF",
		Number: 4,
	},
	{
		Name:   "C",
		Number: 5,
	},
	{
		Name:   "PG",
		Number: 6,
	},
	{
		Name:   "SG",
		Number: 7,
	},
	{
		Name:   "SF",
		Number: 8,
	},
	{
		Name:   "PF",
		Number: 9,
	},
	{
		Name:   "C",
		Number: 10,
	},
	{
		Name:   "Res",
		Number: 11,
	},
	{
		Name:   "Res",
		Number: 12,
	},
}

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action         string           `json:"action"`
	Message        string           `json:"message"`
	Countdown      int              `json:"countdown"`
	PlayerID       int              `json:"player_id"`
	PlayerInfo     []string         `json:"player_info"`
	Color          string           `json:"color"`
	Row            int              `json:"row"`
	Column         int              `json:"column"`
	DraftSeconds   int              `json:"draft_seconds"`
	Teams          []models.NBATeam `json:"teams"`
	MessageType    string           `json:"message_type"`
	ConnectedUsers []string         `json:"connected_users"`
}

type WsPayload struct {
	Action       string              `json:"action"`
	Username     string              `json:"username"`
	Countdown    int                 `json:"countdown"`
	PlayerID     int                 `json:"player_id"`
	PlayerInfo   []string            `json:"player_info"`
	Color        string              `json:"color"`
	Row          int                 `json:"row"`
	Column       int                 `json:"column"`
	DraftSeconds int                 `json:"draft_seconds"`
	Teams        []models.NBATeam    `json:"nba_teams"`
	Message      string              `json:"message"`
	Conn         WebSocketConnection `json:"-"`
}

// WsEndPoint upgrades connection to websocket
func (m *Repository) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small><em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	reset := false
	quit := false
	draftOrder := []int{}
	color := "transparent"
	var rowCounter int
	var colCounter int
	var draftCountdown int

	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// get a list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)

		case "timer":
			response.Action = "timer"
			response.Countdown = e.Countdown
			broadcastToAll(response)

		case "generate_order":
			response.Action = "generate_order"
			teams := Repo.getDraftOrder()
			draftOrder = []int{}
			for _, team := range teams {
				draftOrder = append(draftOrder, int(team.TeamID.Int64))
			}
			response.Teams = teams
			broadcastToAll(response)

		case "stop_draft":
			fmt.Println("draft stopped")
			// response.Action = "stop_draft"
			quit = true
			broadcastToAll(response)

		case "start_draft":
			fmt.Println("draft started")
			rowCounter = 1
			colCounter = 1
			draftCountdown = e.Countdown
			go func() {
				timeLeft := draftCountdown
				for {
					switch {
					case reset:
						timeLeft = draftCountdown
						reset = false
						continue
					case quit:
						response.Action = "draft_ended"
						broadcastToAll(response)
						reset = false
						quit = false
						rowCounter = 1
						colCounter = 1
						fmt.Println("draft ended")
						return
					default:
						response.Action = "timer"
						response.Countdown = timeLeft
						response.DraftSeconds = draftCountdown
						broadcastToAll(response)
						time.Sleep(1000 * time.Millisecond)
						if timeLeft <= 0 {
							playerID, firstName, lastName, primary, secondary := Repo.getRandomPlayer()
							Repo.draftPlayer(draftOrder[colCounter-1], playerID)
							response.Action = "draft_player"
							reset = true
							response.Row = rowCounter
							response.Column = colCounter
							response.PlayerID = playerID
							positions := primary
							switch positions {
							case "PG":
								color = "#FDD8E6"
							case "SG":
								color = "#FDD8E6"
							case "SF":
								color = "#C1EBE7"
							case "PF":
								color = "#C4E7FD"
							case "C":
								color = "#C4E7FD"
							}
							if secondary != "" {
								if secondary == "SF" {
									color = "#C1EBE7"
								}
								positions += "/" + secondary
							}
							response.Color = color
							response.Message = fmt.Sprintf("%s<br><strong>%s</strong><br>%s", firstName, lastName, positions)
							broadcastToAll(response)
							if rowCounter%2 == 0 {
								colCounter -= 1
							} else {
								colCounter += 1
							}
							if rowCounter == 12 && colCounter == 0 {
								quit = true
							}
							if colCounter == 5 {
								colCounter = 4
								rowCounter += 1
							} else if colCounter == 0 {
								colCounter = 1
								rowCounter += 1
							}
							continue
						}
						timeLeft -= 1
					}
				}
			}()

		case "reset_players":
			response.Action = "reset_players"
			quit = true
			Repo.resetPlayers()
			broadcastToAll(response)

		case "draft_player":
			Repo.draftPlayer(draftOrder[colCounter-1], e.PlayerID)
			response.Action = "draft_player"
			reset = true
			response.Row = rowCounter
			response.Column = colCounter
			response.PlayerID = e.PlayerID
			firstName := e.PlayerInfo[0]
			lastName := e.PlayerInfo[1]
			positions := e.PlayerInfo[2]
			switch positions {
			case "PG":
				color = "#FDD8E6"
			case "SG":
				color = "#FDD8E6"
			case "SF":
				color = "#C1EBE7"
			case "PF":
				color = "#C4E7FD"
			case "C":
				color = "#C4E7FD"
			}
			if e.PlayerInfo[3] != "" {
				if e.PlayerInfo[3] == "SF" {
					color = "#C1EBE7"
				}
				positions += "/" + e.PlayerInfo[3]
			}
			response.Color = color
			response.Message = fmt.Sprintf("%s<br><strong>%s</strong><br>%s", firstName, lastName, positions)
			broadcastToAll(response)
			if rowCounter%2 == 0 {
				colCounter -= 1
			} else {
				colCounter += 1
			}
			if rowCounter == 12 && colCounter == 0 {
				quit = true
			}
			if colCounter == 5 {
				colCounter = 4
				rowCounter += 1
			} else if colCounter == 0 {
				colCounter = 1
				rowCounter += 1
			}
		}
	}
}

func (m *Repository) getRandomPlayer() (playerID int, firstName, lastName, primary, secondary string) {
	random := rand.Intn(5)
	player, err := m.DB.GetRandomNBAPlayer(random)
	if err != nil {
		fmt.Println(err)
	}
	return player.PlayerID, player.FirstName, player.LastName, player.PrimaryPosition, player.SecondaryPosition
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}

func (m *Repository) draftPlayer(teamID, playerID int) {
	err := m.DB.AddNBAPlayer(playerID, teamID)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Repository) resetPlayers() {
	err := m.DB.DropAllNBAPlayers()
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Repository) getDraftOrder() []models.NBATeam {
	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		return nil
	}
	// Shuffle teams
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(teams), func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })
	return teams
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signin.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles logging in the user
func (m *Repository) PostSignIn(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "signin.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, _, accessLevel, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "user_name", user.FirstName)
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", accessLevel)
	http.Redirect(w, r, "/nba", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) NBAHome(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "nba_home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAPlayers(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "update" {
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		nullInt := true

		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			log.Println(err)
			nullInt = false
			// helpers.ServerError(w, err)
		}

		player := models.NBAPlayer{
			PlayerID: playerID,
			TeamID:   sql.NullInt64{int64(teamID), nullInt},
			Assigned: 0,
		}

		err = m.DB.UpdateNBAPlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	players, err := m.DB.GetNBAPlayersWithBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	/* badges, err := m.DB.GetNBABadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	playersBadges, err := m.DB.GetNBAPlayersWithBadgesBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	} */
	data := make(map[string]interface{})
	data["nba_players"] = players
	data["nba_teams"] = teams
	// data["nba_badges"] = badges
	// data["nba_players_badges"] = playersBadges

	render.Template(w, r, "nba_players.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostNBAPlayers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(r.Form.Get("player_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	nullInt := true

	teamID, err := strconv.Atoi(r.Form.Get("team_id"))
	if err != nil {
		log.Println(err)
		nullInt = false
		// helpers.ServerError(w, err)
	}

	player := models.NBAPlayer{
		PlayerID: playerID,
		TeamID:   sql.NullInt64{int64(teamID), nullInt},
		Assigned: 0,
	}

	err = m.DB.UpdateNBAPlayer(player)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "nba_players", player)

	http.Redirect(w, r, "/nba/players", http.StatusSeeOther)
}

func (m *Repository) NBATeams(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "drop" {
		playerID, err := strconv.Atoi(r.FormValue("playerID"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		err = m.DB.DropNBAPlayer(playerID)
		if err != nil {
			helpers.ServerError(w, err)
		}
	} else if r.FormValue("action") == "add" {
		playerID, err := strconv.Atoi(r.FormValue("playerID"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		teamID, err := strconv.Atoi(r.FormValue("teamID"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		err = m.DB.AddNBAPlayer(playerID, teamID)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	if r.FormValue("player_id") != "" {
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		assigned, err := strconv.Atoi(r.FormValue("assigned"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		player := models.NBAPlayer{
			PlayerID: playerID,
			TeamID:   sql.NullInt64{int64(teamID), true},
			Assigned: assigned,
		}
		err = m.DB.AssignNBAPlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyTeamInfo models.NBATeamInfo
	data := make(map[string]interface{})

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	players, err := m.DB.GetNBAPlayersWithoutBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["teamInfo"] = emptyTeamInfo
	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "nba_teams.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNBATeams(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	text := r.Form.Get("dark_text")
	darkText := "false"
	if text == "true" {
		darkText = "true"
	}

	teamID, err := strconv.Atoi(r.Form.Get("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	teamInfo := models.NBATeamInfo{
		ID:           teamID,
		Name:         r.Form.Get("team_name"),
		Abbreviation: r.Form.Get("abbreviation"),
		Color1:       r.Form.Get("team_color1"),
		Color2:       r.Form.Get("team_color2"),
		DarkText:     darkText,
	}

	form := forms.New(r.PostForm)

	form.Required("team_name", "abbreviation")
	form.MaxLength("abbreviation", 4)

	if !form.Valid() {
		data := make(map[string]interface{})

		teams, err := m.DB.GetNBATeamInfo()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		players, err := m.DB.GetNBAPlayersWithoutBadges()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data["nba_players"] = players
		data["nba_teams"] = teams
		data["teamInfo"] = teamInfo
		data["positions"] = positions

		render.Template(w, r, "nba_teams.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateNBATeamInfo(teamInfo)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "team_info", teamInfo)

	http.Redirect(w, r, "/nba/teams", http.StatusSeeOther)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Handles request for availability and sends JSON response
func (m *Repository) NBATeamsAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

func (m *Repository) NBATeamInfoSummary(w http.ResponseWriter, r *http.Request) {
	teamInfo, ok := m.App.Session.Get(r.Context(), "team_info").(models.NBATeamInfo)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Cant get team info from session")
		http.Redirect(w, r, "/nba/teams", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "team_info")
	data := make(map[string]interface{})
	data["team_info"] = teamInfo

	render.Template(w, r, "nba_team_info.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NBAResults(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "add" {
		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		result := models.Result{
			HomeTeam:  homeTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			AwayTeam:  awayTeam,
		}
		err = m.DB.AddNBAResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	} else if r.FormValue("action") == "update" {
		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		timestampString := r.FormValue("timestamp")
		layout := "2006-01-02 15:04:05 -0700 MST"
		timestamp, err := time.Parse(layout, timestampString)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		result := models.Result{
			HomeTeam:  homeTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			AwayTeam:  awayTeam,
			Time:      timestamp,
		}
		err = m.DB.UpdateNBAResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}

	} else if r.FormValue("action") == "delete" {
		timestampString := r.FormValue("timestamp")
		layout := "2006-01-02 15:04:05 -0700 MST"
		timestamp, err := time.Parse(layout, timestampString)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		result := models.Result{
			Time: timestamp,
		}
		err = m.DB.DeleteNBAResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyStandings models.Result
	data := make(map[string]interface{})

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	standings, err := m.DB.GetNBAStandings()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	lastResults, err := m.DB.GetLastResults(10)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["result"] = emptyStandings
	data["teams"] = teams
	data["standings"] = standings
	data["last_results"] = lastResults

	render.Template(w, r, "nba_results.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNBAResults(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	homeTeam, err := strconv.Atoi(r.Form.Get("home_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	homeScore, err := strconv.Atoi(r.Form.Get("home_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayScore, err := strconv.Atoi(r.Form.Get("away_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayTeam, err := strconv.Atoi(r.Form.Get("away_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	result := models.Result{
		HomeTeam:  homeTeam,
		HomeScore: homeScore,
		AwayScore: awayScore,
		AwayTeam:  awayTeam,
	}

	form := forms.New(r.PostForm)

	form.Required("home_team", "home_score", "away_score", "away_team")
	form.IsDuplicate("home_team", "away_team", "Home and away have to be different")
	form.IsDuplicate("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["NBAresult"] = result

		render.Template(w, r, "nba_results.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.AddNBAResult(result)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "nba_result", result)

	http.Redirect(w, r, "/nba/results", http.StatusSeeOther)
}

func (m *Repository) NBADraft(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	players, err := m.DB.GetNBAPlayersWithoutBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "nba_draft.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNBATeams(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-nba-teams.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNBAPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := m.DB.GetNBAPlayersWithoutBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["nba_players"] = players

	render.Template(w, r, "admin-nba-players.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminNBAResults(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-nba-results.page.tmpl", &models.TemplateData{})
}

// Shows a single players stats
func (m *Repository) AdminShowNBAPlayer(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	player, err := m.DB.GetNBAPlayerByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["nba_player"] = player
	data["nba_teams"] = teams

	render.Template(w, r, "admin-nba-player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
