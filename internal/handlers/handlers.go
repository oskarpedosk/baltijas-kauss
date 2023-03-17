package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
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

var positions = []models.Positions{
	{Name: "PG", Number: 1}, {Name: "SG", Number: 2}, {Name: "SF", Number: 3}, {Name: "PF", Number: 4}, {Name: "C", Number: 5},
	{Name: "PG", Number: 6}, {Name: "SG", Number: 7}, {Name: "SF", Number: 8}, {Name: "PF", Number: 9}, {Name: "C", Number: 10},
	{Name: "Res", Number: 11}, {Name: "Res", Number: 12}}

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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUser(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "user_name", user.FirstName)
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	m.App.Session.Put(r.Context(), "access_level", accessLevel)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	seasons, err := m.DB.GetSeasons()
	if err != nil {
		helpers.ServerError(w, err)
	}

	results, err := m.DB.GetSeasonResults(0)
	if err != nil {
		helpers.ServerError(w, err)
	}
	
	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}
	var teamsWithoutFA = []models.Team{}
	for _, team := range teams {
		if team.TeamID != 1 {
			teamsWithoutFA = append(teamsWithoutFA, team)
		}
	}
	standings := CalculateStandings(teamsWithoutFA, results)


	data := make(map[string]interface{})
	data["teams"] = teamsWithoutFA
	data["standings"] = standings
	data["activeSeason"] = seasons[0].SeasonID

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
