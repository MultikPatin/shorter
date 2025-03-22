package main

import (
	"main/internal/adapters"
	"main/internal/app"
	"main/internal/config"
	"net/http"
)

func main() {
	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c := config.Parse(logger)

	shorterApp := app.NewApp(c)
	defer shorterApp.Close()

	logger.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	if err := http.ListenAndServe(shorterApp.Addr, shorterApp.Router); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
