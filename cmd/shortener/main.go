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

	producerFS, err := database.NewFileStorage(c.StorageFilePaths, true)
	if err != nil {
		logger.Infow(
			"File storage",
			"error", err.Error(),
		)
	} else {
		logger.Infow(
			"File producerFS created",
			"path", c.StorageFilePaths,
		)
	}
	defer producerFS.Close()

	consumerFS, err := database.NewFileStorage(c.StorageFilePaths, false)
	if err != nil {
		logger.Infow(
			"File storage",
			"error", err.Error(),
		)
	} else {
		logger.Infow(
			"File consumerFS created",
			"path", c.StorageFilePaths,
		)
	}
	defer consumerFS.Close()

	d := database.NewInMemoryDB(producerFS, consumerFS)

	err = d.LoadFromFile()
	if err != nil {
		logger.Infow(
			"Load in memory DB",
			"error", err.Error(),
		)
	}

	app.ShortPre = c.ShortLinkPrefix
	h := app.GetHandlers(d)
	r := app.GetRouters(h)

	logger.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	if err := http.ListenAndServe(c.Addr, r); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
