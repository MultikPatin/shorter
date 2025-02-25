package main

import (
	"go.uber.org/zap"
	"main/internal/app"
	"main/internal/database"
	"net/http"
)

var sugar zap.SugaredLogger

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	c, err := app.ParseConfig()
	if err != nil {
		sugar.Error(err)
	}

	producerFS, err := database.NewFileStorage(c.StorageFilePaths, true)
	if err != nil {
		sugar.Infow(
			"File storage",
			"error", err.Error(),
		)
	} else {
		sugar.Infow(
			"File producerFS created",
			"path", c.StorageFilePaths,
		)
	}
	defer producerFS.Close()

	consumerFS, err := database.NewFileStorage(c.StorageFilePaths, false)
	if err != nil {
		sugar.Infow(
			"File storage",
			"error", err.Error(),
		)
	} else {
		sugar.Infow(
			"File consumerFS created",
			"path", c.StorageFilePaths,
		)
	}
	defer consumerFS.Close()

	d := database.NewInMemoryDB(producerFS, consumerFS)

	err = d.LoadFromFile()
	if err != nil {
		sugar.Infow(
			"Load in memory DB",
			"error", err.Error(),
		)
	}

	app.ShortPre = c.ShortLinkPrefix
	h := app.GetHandlers(d)
	r := app.GetRouters(h)

	sugar.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	if err := http.ListenAndServe(c.Addr, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
