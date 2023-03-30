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

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

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

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
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

	data["players"] = players
	data["teams"] = teams

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