package main

import (
	"main/internal/adapters"
	"main/internal/app"
	"main/internal/database"
	"net/http"
)

func main() {
	logger := adapters.GetLogger()

	c, err := app.ParseConfig(logger)
	if err != nil {
		logger.Error(err)
	}

	InMemoryDB, err := database.NewInMemoryDB(c.StorageFilePaths, logger)
	if err != nil {
		logger.Infow(
			"Create in memory DB",
			"error", err.Error(),
		)
	}
	defer InMemoryDB.Close()

	app.ShortPre = c.ShortLinkPrefix
	h := app.GetHandlers(InMemoryDB)
	r := app.GetRouters(h)

	logger.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	if err := http.ListenAndServe(c.Addr, r); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
