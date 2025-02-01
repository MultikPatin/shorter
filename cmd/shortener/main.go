package main

import (
	"fmt"
	"log"
	"main/internal/app"
	"net/http"
)

func main() {
	err := app.CmdConfig.Parse()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	fmt.Println(app.CmdConfig)

	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(app.CmdConfig.ServHost.String(), r))
}
