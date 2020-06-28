package test_data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Locations struct {
	Locations	[]Location `json:"locations"`
}
type Location struct {
	PlaceName string  `json:"placeName"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

var LocationsData Locations

func LoadCSV(){
	path := "test_data/zips.json"
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &LocationsData)

	if err != nil{
		fmt.Println(err)
	}
	//
	//for i := 0; i < len(LocationsData.Locations); i++ {
	//	fmt.Println(LocationsData.Locations[i].Longitude)
	//}
}

