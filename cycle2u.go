package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Cycle2U v2016")

	_loadConfig()

	_initSMT()
	_initDB()
	_initEcho()
	_initRoutes()

	// Start the web server
	if Config.Debug {
		log.Printf("... Starting Web Server on port %d", Config.WebPort)
	}
	e.Run(fmt.Sprintf(":%d", Config.WebPort))

}
