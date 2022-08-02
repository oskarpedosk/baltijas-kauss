package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadNBAPlayerData(fileName string) []NBAPlayerData {
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened ", fileName)
	}

	// Defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initialize our players array
	var players []NBAPlayerData

	// Unmarshal our byteArray which contains our
	// jsonFile's content into 'players' which we defined above
	json.Unmarshal(byteValue, &players)

	return players
}