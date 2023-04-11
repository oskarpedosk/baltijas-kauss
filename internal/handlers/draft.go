package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

var (
	timeLimit      = 0
	pick           = 1
	rounds         = 2
	randomPlayer   = 5
	draft          = false
	pause          = false
	draftGenerated = false
	draftCompleted = false
	countdown      time.Duration
	draftOrder     = []models.Team{}
	draftPicks     = []models.DraftPick{}
)

type WebSocketConnection struct {
	*websocket.Conn
}

var draftChan = make(chan DraftPayload)
var messengerChan = make(chan MessengerPayload)
var clients = make(map[WebSocketConnection]string)
var clientsUserIDs = make(map[int]WebSocketConnection)

// Define the response sent back from websocket
type DraftJsonResponse struct {
	Action     string             `json:"action"`
	Message    string             `json:"message"`
	Countdown  int                `json:"countdown"`
	PlayerID   int                `json:"player_id"`
	Row        int                `json:"row"`
	Col        int                `json:"col"`
	NextRow    int                `json:"next_row"`
	NextCol    int                `json:"next_col"`
	Pick       int                `json:"pick"`
	TimeLimit  int                `json:"time_limit"`
	TeamName   string             `json:"team_name"`
	Teams      []models.Team      `json:"teams"`
	DraftPicks []models.DraftPick `json:"draft_picks"`
}

type DraftPayload struct {
	Action     string              `json:"action"`
	UserID     int                 `json:"user_id"`
	PlayerID   int                 `json:"player_id"`
	PlayerInfo []string            `json:"player_info"`
	TimeLimit  int                 `json:"time_limit"`
	Conn       WebSocketConnection `json:"-"`
}

type MessengerJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	ConnectedUsers []string `json:"connected_users"`
}

type MessengerPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// Upgrade connection to websocket
func (m *Repository) DraftEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var draftResponse DraftJsonResponse
	draftResponse.Message = `<em><small>Connected to server</small><em>`
	draftResponse.DraftPicks = draftPicks
	draftResponse.Teams = draftOrder

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(draftResponse)
	if err != nil {
		log.Println(err)
	}

	go ListenForDraftWs(&conn)
}

// Upgrade connection to websocket
func (m *Repository) MessengerEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var messengerResponse MessengerJsonResponse
	messengerResponse.Message = `<em><small>Connected to server</small><em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(messengerResponse)
	if err != nil {
		log.Println(err)
	}

	go ListenForMessengerWs(&conn)
}

func ListenForDraftWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var draftPayload DraftPayload

	for {
		err := conn.ReadJSON(&draftPayload)
		if err != nil {
			// do nothing
		} else {
			draftPayload.Conn = *conn
			draftChan <- draftPayload
		}
	}
}

func ListenForMessengerWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var messengerPayload MessengerPayload

	for {
		err := conn.ReadJSON(&messengerPayload)
		if err != nil {
			// do nothing
		} else {
			messengerPayload.Conn = *conn
			messengerChan <- messengerPayload
		}
	}
}

func BroadcastToAll(response interface{}) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func BroadcastToUser(userID int, response interface{}) {
	_, ok := clientsUserIDs[userID]
	if ok {
		err := clientsUserIDs[userID].WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = clientsUserIDs[userID].Close()
			delete(clientsUserIDs, userID)
		}
	}
}

func GetUserList() []string {
	var userList []string
	for _, user := range clients {
		if user != "" {
			userList = append(userList, user)
		}
	}
	sort.Strings(userList)
	return userList
}

func ListenToDraftWsChannel() {
	var response DraftJsonResponse

	for {
		e := <-draftChan

		switch e.Action {

		case "connected":
			clientsUserIDs[e.UserID] = e.Conn

		case "generate_draft":
			if !draft {
				resetDraft()
				generateDraft()
				draftGenerated = true
			}

		case "stop":
			draft = false

		case "pause":
			if draft {
				pause = !pause
			}

		case "draft_player":
			if draft {
				draftPlayer(e)
				pick++
				if pick > len(draftPicks) {
					draft = false
					draftCompleted = true
				}
			}

		case "start":
			if !draft && draftGenerated {
				startDraft(e)
				go draftCountdown()
			}

		case "reset_players":
			if !draft {
				response.Action = "reset_players"
				Repo.ResetPlayers()
				BroadcastToAll(response)
			}
		}
	}
}

func draftPlayer(e DraftPayload) {
	var response DraftJsonResponse
	response.Action = "draft_player"

	name := e.PlayerInfo[0]
	positions := e.PlayerInfo[1]

	countdown = time.Duration(timeLimit) * time.Second

	if pick < len(draftPicks) {
		BroadcastToUser(draftPicks[pick].TeamID, DraftJsonResponse{Action: "your_turn"})
		response.NextRow = draftPicks[pick].Row
		response.NextCol = draftPicks[pick].Col
		response.TeamName = draftPicks[pick].TeamName
	}

	
	response.Pick = pick + 1
	response.PlayerID = e.PlayerID
	response.Row = draftPicks[pick-1].Row
	response.Col = draftPicks[pick-1].Col
	response.Message = fmt.Sprintf("<span style=\"font-size: 14px\">%s</span><br>%s", name, positions)
	BroadcastToAll(response)

	Repo.DraftPlayer(draftPicks[pick-1].TeamID, e.PlayerID)
	draftPicks[pick-1].PlayerID = e.PlayerID
	draftPicks[pick-1].Name = name
	draftPicks[pick-1].Positions = positions
}

func generateDraft() {
	teams := Repo.GenerateDraftOrder()
	for row := 1; row <= rounds; row++ {
		if row%2 != 0 {
			col := 1
			for j := 0; j < len(teams); j++ {
				draftPick := models.DraftPick{
					Row:      row,
					Col:      col,
					Pick:     pick,
					TeamID:   teams[j].TeamID,
					TeamName: teams[j].Name,
				}
				draftPicks = append(draftPicks, draftPick)
				pick++
				col++
			}
		} else {
			col := len(teams)
			for j := len(teams) - 1; j >= 0; j-- {
				draftPick := models.DraftPick{
					Row:      row,
					Col:      col,
					Pick:     pick,
					TeamID:   teams[j].TeamID,
					TeamName: teams[j].Name,
				}
				draftPicks = append(draftPicks, draftPick)
				pick++
				col--
			}
		}
	}
	pick = 1
	var response DraftJsonResponse
	draftOrder = teams
	response.Teams = teams
	response.Action = "generate_draft"
	BroadcastToAll(response)
}

func draftCountdown() {
	var response DraftJsonResponse
	for draft {
		if !pause {
			response.Action = "countdown"
			response.Countdown = int(countdown / time.Second)
			response.TimeLimit = timeLimit
			BroadcastToAll(response)
			countdown -= time.Second
			time.Sleep(time.Second)
			if countdown < 0 {
				playerID, firstName, lastName, primary, secondary := Repo.SelectRandomPlayer()
				name := firstName + " " + lastName
				positions := primary
				if secondary != "" {
					positions += "/" + secondary
				}
				draftPlayer(DraftPayload{
					PlayerID:   playerID,
					PlayerInfo: []string{name, positions},
				})
				pick++
				if pick > len(draftPicks) {
					draft = false
					draftCompleted = true
				}
				countdown = time.Duration(timeLimit) * time.Second
			}
		}
	}
	if draftCompleted {
		response = DraftJsonResponse{}
		response.Action = "draft_complete"
		BroadcastToAll(response)
		draftID, err := Repo.DB.GetDraftID()
		if err != nil {
			log.Println(err)
			return
		}
		for _, pick := range draftPicks {
			err = Repo.DB.AddDraftPick(draftID+1, pick)
			if err != nil {
				log.Println(err)
				return
			}
		}
		draftGenerated = false
	}
}

func startDraft(e DraftPayload) {
	pick = 1
	draft = true
	pause = false
	timeLimit = 0
	countdown = 0

	var response DraftJsonResponse
	response.Action = "start"
	response.TeamName = draftPicks[pick-1].TeamName
	BroadcastToAll(response)
	BroadcastToUser(draftPicks[pick-1].TeamID, DraftJsonResponse{Action: "your_turn"})

	timeLimit = e.TimeLimit
	countdown = time.Duration(timeLimit) * time.Second
}

func resetDraft() {
	pick = 1
	timeLimit = 0
	countdown = 0
	pause = false
	draftGenerated = false
	draftCompleted = false
	draftOrder = []models.Team{}
	draftPicks = []models.DraftPick{}
}

func ListenToMessengerWsChannel() {
	var response MessengerJsonResponse

	for {
		e := <-messengerChan

		switch e.Action {
		case "username":
			// get a list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := GetUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			BroadcastToAll(response)

		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := GetUserList()
			response.ConnectedUsers = users
			BroadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			BroadcastToAll(response)
		}
	}
}

func (m *Repository) SelectRandomPlayer() (playerID int, firstName, lastName, primary, secondary string) {
	random := rand.Intn(randomPlayer)
	player, err := m.DB.SelectRandomPlayer(random)
	if err != nil {
		log.Println(err)
	}
	return player.PlayerID, player.FirstName, player.LastName, player.PrimaryPosition, player.SecondaryPosition
}

func (m *Repository) DraftPlayer(teamID, playerID int) {
	err := m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) ResetPlayers() {
	err := m.DB.ResetPlayers()
	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) GenerateDraftOrder() []models.Team {
	teams, err := m.DB.GetTeams()
	if err != nil {
		return nil
	}
	teamsWithoutFA := teams[1:]

	// Shuffle teams
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(teamsWithoutFA), func(i, j int) { teamsWithoutFA[i], teamsWithoutFA[j] = teamsWithoutFA[j], teamsWithoutFA[i] })
	return teamsWithoutFA
}

func (m *Repository) Draft(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	filter := models.Filter{
		TeamID:              0,
		HeightMin:           150,
		HeightMax:           250,
		WeightMin:           50,
		WeightMax:           150,
		OverallMin:          1,
		OverallMax:          99,
		ThreePointShotMin:   1,
		ThreePointShotMax:   99,
		DrivingDunkMin:      1,
		DrivingDunkMax:      99,
		AthleticismMin:      1,
		AthleticismMax:      99,
		PerimeterDefenseMin: 1,
		PerimeterDefenseMax: 99,
		InteriorDefenseMin:  1,
		InteriorDefenseMax:  99,
		ReboundingMin:       1,
		ReboundingMax:       99,
		Position1:           1,
		Position2:           1,
		Position3:           1,
		Position4:           1,
		Position5:           1,
		Limit:               1000,
		Offset:              0,
		Col1:                "overall",
		Col2:                "\"attributes/TotalAttributes\"",
		Order:               "desc",
	}

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teamsWithoutFA := teams[1:]

	data["players"] = players
	data["teams"] = teamsWithoutFA

	render.Template(w, r, "draft.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
