package main

import (
	"log"
	"main/internal/app"
	"net/http"
)

func main() {
	c, err := app.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.ShortPre = c.ShortLinkPrefix

	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(c.Addr, r))
}
