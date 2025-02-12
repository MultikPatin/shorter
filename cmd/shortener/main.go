package main

import (
	"log"
	"main/internal/app"
	"main/internal/db"
	"net/http"
)

func main() {
	c, err := app.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.ShortPre = c.ShortLinkPrefix

	d := db.NewInMemoryDB()
	h := app.GetHeaders(d)
	r := app.GetRouter(h)

	log.Fatal(http.ListenAndServe(c.Addr, r))
}
