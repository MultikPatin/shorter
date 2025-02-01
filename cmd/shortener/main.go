package main

import (
	"log"
	"main/internal/app"
	"net/http"
)

func main() {
	r := app.GetRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
