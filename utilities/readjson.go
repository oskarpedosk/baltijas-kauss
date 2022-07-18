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

	return players
}
