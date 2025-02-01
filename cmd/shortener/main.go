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

	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(app.CmdConfig.ServHost.String(), r))
}
