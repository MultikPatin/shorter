package main

import (
	"main/internal/app"
	"net/http"
)

func main() {
	router := app.GetRouter()

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
