package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
	Action         string        `json:"action"`
	Message        string        `json:"message"`
	Countdown      int           `json:"countdown"`
	PlayerID       int           `json:"player_id"`
	PlayerInfo     []string      `json:"player_info"`
	Color          string        `json:"color"`
	Row            int           `json:"row"`
	Column         int           `json:"column"`
	DraftSeconds   int           `json:"draft_seconds"`
	Teams          []models.Team `json:"teams"`
	MessageType    string        `json:"message_type"`
	ConnectedUsers []string      `json:"connected_users"`
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
	Teams        []models.Team       `json:"nba_teams"`
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
				draftOrder = append(draftOrder, int(team.TeamID))
			}
			response.Teams = teams
			broadcastToAll(response)

		case "stop_draft":
			fmt.Println("draft stopped")
			// response.Action = "stop_draft"
			quit = true
			broadcastToAll(response)

		case "start_draft":
			response.Action = "draft_started"
			broadcastToAll(response)
			fmt.Println("draft started")
			rowCounter = 1
			colCounter = 1
			response.Action = ""
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
							playerID, firstName, lastName, primary, secondary := Repo.GetRandomPlayer()
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
							response.Message = fmt.Sprintf("%s <span class=\"fw-semibold\">%s</span><br>%s", firstName, lastName, positions)
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
			response.Message = fmt.Sprintf("%s <span class=\"fw-semibold\">%s</span><br>%s", firstName, lastName, positions)
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

func (m *Repository) GetRandomPlayer() (playerID int, firstName, lastName, primary, secondary string) {
	random := rand.Intn(5)
	player, err := m.DB.GetRandomPlayer(random)
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
	err := m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Repository) resetPlayers() {
	err := m.DB.ResetPlayers()
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Repository) getDraftOrder() []models.Team {
	teams, err := m.DB.GetTeams()
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

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles logging in the user
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
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

	user, err := m.DB.GetUser(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "user_name", user.FirstName)
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", accessLevel)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
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

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// Player is the single player handler
func (m *Repository) Player(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	player, err := m.DB.GetPlayer(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var playersTeam models.Team
	for _, team := range teams {
		if team.TeamID == player.TeamID {
			playersTeam = team
			break
		}
	}

	data := make(map[string]interface{})
	data["player"] = player
	data["team"] = playersTeam
	data["teams"] = teams[1:]
	data["FA"] = teams[0]

	render.Template(w, r, "player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostPlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(r.FormValue("player_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	teamID, err := strconv.Atoi(r.FormValue("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

func (m *Repository) Players(w http.ResponseWriter, r *http.Request) {
	page := 1
	perPage := 20
	// Get page number
	re := regexp.MustCompile(`\/page=(\d+)`)
	match := re.FindStringSubmatch(r.RequestURI)
	if len(match) > 1 {
		page, _ = strconv.Atoi(match[1])
	}

	pagination, err := m.DB.GetPaginationData(page, perPage, "players", "/players")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	players, err := m.DB.GetPlayers(perPage, pagination.Offset)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var playersWithTeamInfo []models.PlayerWithTeamInfo

	for _, player := range players {
		for _, team := range teams {
			if player.TeamID == team.TeamID {
				playerWithTeamInfo := models.PlayerWithTeamInfo{
					Player: player,
					Team:   team,
				}
				playersWithTeamInfo = append(playersWithTeamInfo, playerWithTeamInfo)
				break
			}
		}
	}

	ranking := []int{}
	for i := 1; i <= len(players); i++ {
		ranking = append(ranking, i+pagination.Offset)
	}

	data := make(map[string]interface{})
	data["players"] = playersWithTeamInfo
	data["teams"] = teams[1:]
	data["FA"] = teams[0]
	data["ranking"] = ranking
	data["pagination"] = pagination

	render.Template(w, r, "players.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostPlayers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(r.FormValue("player_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	teamID, err := strconv.Atoi(r.FormValue("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	player := models.Player{
		PlayerID:         playerID,
		TeamID:           teamID,
		AssignedPosition: 0,
	}

	err = m.DB.SwitchTeam(player)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "nba_players", player)

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func (m *Repository) PostUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	playerID := r.FormValue("player_id")
	ratingsURL := r.FormValue("ratings_url")

	go func(playerID, ratingsURL string) {
		filePath := "./static/js/script/updateplayer.js"
		cmd := exec.Command("node", filePath, playerID, ratingsURL)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		var player models.Player
		json.Unmarshal(output, &player)

		err = m.DB.UpdatePlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}(playerID, ratingsURL)

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func (m *Repository) Team(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	fmt.Println(teamID)
	if r.FormValue("action") == "drop" {
		playerID, err := strconv.Atoi(r.FormValue("playerID"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		err = m.DB.DropPlayer(playerID)
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

		err = m.DB.AddPlayer(playerID, teamID)
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
		assigned, err := strconv.Atoi(r.FormValue("AssignedPosition"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		player := models.Player{
			PlayerID:         playerID,
			TeamID:           teamID,
			AssignedPosition: assigned,
		}
		err = m.DB.AssignPosition(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyTeamInfo models.Team
	data := make(map[string]interface{})

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	players, err := m.DB.GetPlayers(150, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["teamInfo"] = emptyTeamInfo
	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "team.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostTeam(w http.ResponseWriter, r *http.Request) {
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

	teamInfo := models.Team{
		TeamID:       teamID,
		Name:         r.Form.Get("team_name"),
		Abbreviation: r.Form.Get("abbreviation"),
		Color1:       r.Form.Get("team_color1"),
		Color2:       r.Form.Get("team_color2"),
		DarkText:     darkText,
	}

	form := forms.New(r.PostForm)

	form.Required("team_name", "abbreviation")
	form.AlphaNumeric("team_name", "abbreviation")
	form.MaxLength("team_name", 20)
	form.IsUpper("abbreviation")
	form.MaxLength("abbreviation", 4)

	if !form.Valid() {
		data := make(map[string]interface{})

		teams, err := m.DB.GetTeams()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		players, err := m.DB.GetPlayers(150, 0)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data["nba_players"] = players
		data["nba_teams"] = teams
		data["teamInfo"] = teamInfo
		data["positions"] = positions

		errMsg := form.Errors.Get("team_name")
		if errMsg == "" {
			errMsg = form.Errors.Get("abbreviation")
		}
		m.App.Session.Put(r.Context(), "error", errMsg)
		render.Template(w, r, "team.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateTeam(teamInfo)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "flash", "Team updated successfully!")
	http.Redirect(w, r, "/team/"+string(teamID), http.StatusSeeOther)
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
			HomeTeamID: homeTeam,
			HomeScore:  homeScore,
			AwayScore:  awayScore,
			AwayTeamID: awayTeam,
		}
		err = m.DB.AddResult(result)
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

		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		result := models.Result{
			HomeTeamID: homeTeam,
			HomeScore:  homeScore,
			AwayScore:  awayScore,
			AwayTeamID: awayTeam,
		}
		err = m.DB.UpdateResult(result)
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
			CreatedAt: timestamp,
		}
		err = m.DB.DeleteResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyStandings models.Result
	data := make(map[string]interface{})

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	standings, err := m.DB.GetStandings()
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

	render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
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
		HomeTeamID: homeTeam,
		HomeScore:  homeScore,
		AwayScore:  awayScore,
		AwayTeamID: awayTeam,
	}

	form := forms.New(r.PostForm)

	form.Required("home_team", "home_score", "away_score", "away_team")
	form.IsDuplicate("home_team", "away_team", "Home and away have to be different")
	form.IsDuplicate("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["NBAresult"] = result

		render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.AddResult(result)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "result", result)

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func (m *Repository) NBADraft(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	players, err := m.DB.GetPlayers(200, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "draft.page.tmpl", &models.TemplateData{
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
	players, err := m.DB.GetPlayers(120, 0)
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

	player, err := m.DB.GetPlayer(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
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
