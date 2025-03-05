package main

import (
	"main/internal/adapters"
	"main/internal/adapters/database/memory"
	"main/internal/app"
	"main/internal/services"
	"net/http"
)

func main() {
	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c, err := app.ParseConfig(logger)
	if err != nil {
		logger.Error(err)
	}

	InMemoryDB, err := memory.NewInMemoryDB(c.StorageFilePaths, logger)
	if err != nil {
		logger.Infow(
			"Create in memory DB",
			"error", err.Error(),
		)
	}
	defer InMemoryDB.Close()

	linksService := services.NewLinksService(InMemoryDB, c.ShortLinkPrefix)

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
