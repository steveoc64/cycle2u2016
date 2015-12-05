package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// Runtime variables, held in external file config.json
type ConfigType struct {
	Debug          bool
	DataSourceName string
	WebPort        int
}

var Config ConfigType

// Load the config.json file, and override with runtime flags
func _loadConfig() {
	cf, err := os.Open("config.json")
	if err != nil {
		log.Println("Could not open config.json :", err.Error())
	}

	data := json.NewDecoder(cf)
	if err = data.Decode(&Config); err != nil {
		log.Fatalln("Failed to load config.json :", err.Error())
	}

	flag.BoolVar(&Config.Debug, "debug", Config.Debug, "Enable Debugging")
	flag.StringVar(&Config.DataSourceName, "", Config.DataSourceName, "DataSourceName for SQLServer")
	flag.IntVar(&Config.WebPort, "webport", Config.WebPort, "Port Number for Web Server")
	flag.Parse()

	log.Printf("Starting\n\tDebug: \t\t%t\n\tSQLServer: \t%s\n\tWeb Port: \t%d\n",
		Config.Debug,
		Config.DataSourceName,
		Config.WebPort)
}
