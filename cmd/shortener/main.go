package main

import (
	"log"
	"main/internal/app"
	"net/http"
)

func main() {
	err := app.CmdConfig.Parse()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%v\n", app.CmdConfig)
	err = app.EnvConfig.Parse()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%v\n", app.CmdConfig)
	addr := app.EnvConfig.ServHost
	if addr == "" {
		addr = app.CmdConfig.ServHost.String()
	}

	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(addr, r))
}
