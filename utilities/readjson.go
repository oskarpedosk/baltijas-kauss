package utilities

import (
	"fmt"
	"os"
)

func ReadJson(fileName string) {
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened ", fileName)
	// Defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

}
