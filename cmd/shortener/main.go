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
		sugar.Fatal(err)
	}
	app.ShortPre = c.ShortLinkPrefix

	d := database.NewInMemoryDB()
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
