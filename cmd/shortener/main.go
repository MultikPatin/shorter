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
	err = app.EnvConfig.Parse()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	addr := app.EnvConfig.ServHost
	app.ShortPre = app.EnvConfig.ShorLink
	if addr == "" {
		addr = app.CmdConfig.ServHost.String()
		app.ShortPre = app.CmdConfig.ShorLink.Addr
	}

	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(addr, r))
}
