package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

)

func WriteToJson(fileName string, sliceOfPlayers interface{}) {
	// Set file path and type
	fileName = "../../static/jsondata/" + fileName + ".json"

	// Create a file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	// Write scraped data to .json
	data, err := json.MarshalIndent(sliceOfPlayers, "", "	")
	if err != nil {
		fmt.Println(err.Error())
	}
	ioutil.WriteFile(fileName, data, 0644)
}
