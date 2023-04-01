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

type WebSocketConnection struct {
	*websocket.Conn
}

var draftChan = make(chan DraftPayload)
var messengerChan = make(chan MessengerPayload)
var clients = make(map[WebSocketConnection]string)

// Define the response sent back from websocket
type DraftJsonResponse struct {
	Action         string        `json:"action"`
	Message        string        `json:"message"`
	Countdown      int           `json:"countdown"`
	PlayerID       int           `json:"player_id"`
	PlayerInfo     []string      `json:"player_info"`
	Row            int           `json:"row"`
	Col            int           `json:"col"`
	TimeLimit      int           `json:"time_limit"`
	Teams          []models.Team `json:"teams"`
	ConnectedUsers []string      `json:"connected_users"`
}

type DraftPayload struct {
	Action     string              `json:"action"`
	Username   string              `json:"username"`
	Countdown  int                 `json:"countdown"`
	PlayerID   int                 `json:"player_id"`
	PlayerInfo []string            `json:"player_info"`
	Row        int                 `json:"row"`
	Col        int                 `json:"col"`
	TimeLimit  int                 `json:"time_limit"`
	Teams      []models.Team       `json:"nba_teams"`
	Message    string              `json:"message"`
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

	log.Println("Client connected to draft endpoint")

	var draftResponse DraftJsonResponse
	draftResponse.Message = `<em><small>Connected to server</small><em>`

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

	log.Println("Client connected to messenger endpoint")

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

func GetUserList() []string {
	var userList []string
	for _, user := range clients {
		fmt.Println(user)
		if user != "" {
			userList = append(userList, user)
		}
	}
	sort.Strings(userList)
	return userList
}

func ListenToDraftWsChannel() {
	var response DraftJsonResponse

	// reset := false
	// quit := false
	draftOrder := []int{}
	draft := false
	// var rowCounter int
	// var colCounter int
	// var draftCountdown int

	for {
		e := <-draftChan

		switch e.Action {
		case "timer":
			response.Action = "timer"
			response.Countdown = e.Countdown
			BroadcastToAll(response)

		case "generate_order":
			response.Action = "generate_order"
			teams := Repo.GenerateDraftOrder()
			draftOrder = []int{}
			for _, team := range teams {
				draftOrder = append(draftOrder, int(team.TeamID))
			}
			response.Teams = teams
			BroadcastToAll(response)

		case "stop_draft":
			fmt.Println("draft stopped")
			// response.Action = "stop_draft"
			// quit = true
			BroadcastToAll(response)

		case "start":
			response.Action = "start"
			draft = true
			BroadcastToAll(response)
			fmt.Println("draft started")
			// rowCounter = 1
			// colCounter = 1
			// response.Action = ""
			// draftCountdown = e.Countdown
			// go func() {
			// 	timeLeft := draftCountdown
			// 	for {
			// 		switch {
			// 		case reset:
			// 			timeLeft = draftCountdown
			// 			reset = false
			// 			continue
			// 		case quit:
			// 			response.Action = "draft_ended"
			// 			BroadcastToAll(response)
			// 			reset = false
			// 			quit = false
			// 			rowCounter = 1
			// 			colCounter = 1
			// 			fmt.Println("draft ended")
			// 			return
			// 		default:
			// 			response.Action = "timer"
			// 			response.Countdown = timeLeft
			// 			response.TimeLimit = draftCountdown
			// 			BroadcastToAll(response)
			// 			time.Sleep(1000 * time.Millisecond)
			// 			if timeLeft <= 0 {
			// 				playerID, firstName, lastName, primary, secondary := Repo.GetRandomPlayer()
			// 				Repo.DraftPlayer(draftOrder[colCounter-1], playerID)
			// 				response.Action = "draft_player"
			// 				reset = true
			// 				response.Row = rowCounter
			// 				response.Col = colCounter
			// 				response.PlayerID = playerID
			// 				positions := primary
			// 				if secondary != "" {
			// 					positions += "/" + secondary
			// 				}
			// 				response.Message = fmt.Sprintf("%s <span class=\"fw-semibold\">%s</span><br>%s", firstName, lastName, positions)
			// 				BroadcastToAll(response)
			// 				if rowCounter%2 == 0 {
			// 					colCounter -= 1
			// 				} else {
			// 					colCounter += 1
			// 				}
			// 				if rowCounter == 12 && colCounter == 0 {
			// 					quit = true
			// 				}
			// 				if colCounter == 5 {
			// 					colCounter = 4
			// 					rowCounter += 1
			// 				} else if colCounter == 0 {
			// 					colCounter = 1
			// 					rowCounter += 1
			// 				}
			// 				continue
			// 			}
			// 			timeLeft -= 1
			// 		}
			// 	}
			// }()

		case "reset_players":
			response.Action = "reset_players"
			// quit = true
			Repo.ResetPlayers()
			BroadcastToAll(response)

		case "draft_player":
			// Repo.DraftPlayer(draftOrder[colCounter-1], e.PlayerID)
			response.Action = "draft_player"
			// reset = true
			// response.Row = rowCounter
			// response.Col = colCounter
			response.PlayerID = e.PlayerID
			firstName := e.PlayerInfo[0]
			lastName := e.PlayerInfo[1]
			positions := e.PlayerInfo[2]
			response.Message = fmt.Sprintf("%s <span class=\"fw-semibold\">%s</span><br>%s", firstName, lastName, positions)
			BroadcastToAll(response)
			// if rowCounter%2 == 0 {
			// 	colCounter -= 1
			// } else {
			// 	colCounter += 1
			// }
			// if rowCounter == 12 && colCounter == 0 {
			// 	quit = true
			// }
			// if colCounter == 5 {
			// 	colCounter = 4
			// 	rowCounter += 1
			// } else if colCounter == 0 {
			// 	colCounter = 1
			// 	rowCounter += 1
			// }
		}
	}
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

func (m *Repository) GetRandomPlayer() (playerID int, firstName, lastName, primary, secondary string) {
	random := rand.Intn(5)
	player, err := m.DB.GetRandomPlayer(random)
	if err != nil {
		fmt.Println(err)
	}
	return player.PlayerID, player.FirstName, player.LastName, player.PrimaryPosition, player.SecondaryPosition
}

func (m *Repository) DraftPlayer(teamID, playerID int) {
	err := m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Repository) ResetPlayers() {
	err := m.DB.ResetPlayers()
	if err != nil {
		fmt.Println(err)
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

// duration := 5 * time.Second
// 	draft := true
// 	counter := 3

// 	for draft {
// 		if duration < 0 {
// 			fmt.Println("Time's up!")
// 			duration = 3 * time.Second
// 			counter--
// 			if counter == 0 {
// 				draft = false
// 				continue
// 			}
// 		}

// 		fmt.Printf("Time remaining: %d\n", duration/time.Second)
// 		duration -= time.Second
// 		time.Sleep(time.Second)
// 	}
// 	fmt.Println("Time's up!")
