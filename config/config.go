package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	config *configStruct

	Database        string
	Collection      string
	Path            string
	Servername      string
	Serverport      string
	Afkkickvaluemin int
	Afkkickvaluemax int
	Kickreason      string
)

type configStruct struct {
	Database        string `json:"databasepath"`
	Collection      string `json:"databasename"`
	Path            string `json:"statspath"`
	Servername      string `json:"servername"`
	Serverport      string `json:"serverport"`
	Afkkickvaluemin int    `json:"afkkickvaluemin"`
	Afkkickvaluemax int    `json:"afkkickvaluemax"`
	Kickreason      string `json:"kickreason"`
}

func ReadConfig() error {

	// read in config.json file and export out required variables.

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading config.json", err)
		os.Exit(2)
	}
	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error converting json to variables", err)
		os.Exit(2)
	}

	Database = config.Database
	Collection = config.Collection
	Path = config.Path
	Servername = config.Servername
	Serverport = config.Serverport
	Afkkickvaluemin = config.Afkkickvaluemin
	Afkkickvaluemax = config.Afkkickvaluemax
	Kickreason = config.Kickreason

	return nil
}
