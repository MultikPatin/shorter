package main

import (
	"main/internal/adapters"
	"main/internal/app"
	"main/internal/config"
	"main/internal/services"
	"net/http"
)

func main() {
	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c := config.Parse(logger)

	database, err := adapters.GetDatabase(c, logger)
	if err != nil {
		panic(err)
	}

	linksService := services.NewLinksService(c, database)
	defer linksService.Close()

	h := app.GetHandlers(linksService)
	r := app.GetRouters(h)

	logger.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	if err := http.ListenAndServe(c.Addr, r); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
