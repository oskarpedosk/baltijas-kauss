package utilities

import (
	"encoding/json"
	"io/ioutil"
	"os"
)


func writeToJson(fileName string, sliceOfPlayers interface{}) {
	// Create a file
	file,_ := os.Create(fileName)
	defer file.Close()
	
	// Write scraped data to .json
	data,_ := json.MarshalIndent(sliceOfPlayers, "", "	")	
	ioutil.WriteFile(fileName, data, 0644)
}

