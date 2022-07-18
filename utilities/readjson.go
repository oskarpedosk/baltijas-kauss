package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJson(fileName string) []PlayerData {
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened ", fileName)
	// Defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initialize our players array
	var players []PlayerData

	// Unmarshal our byteArray which contains our
	// jsonFile's content into 'players' which we defined above
	json.Unmarshal(byteValue, &players)

	// Iterate through every player within our players array
	for i := 0; i < len(players); i++ {
		fmt.Println("Rank: ", players[i].Rank)
		fmt.Println(players[i].Name, players[i].Team)
		fmt.Println("Positions: ", players[i].Positions)
		fmt.Println("Height: ", players[i].Height)
		fmt.Println("Overall: ", players[i].OverallRating)
		fmt.Println("3pt: ", players[i].ThreePointRating)
		fmt.Println("Dunk: ", players[i].DunkRating)
		fmt.Println("--------------")
	}
	return players
}
