package main

import (
	"main/internal/adapters"
	"main/internal/app"
	"main/internal/config"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c := config.Parse(logger)

	shorterApp, err := app.NewApp(c)
	if err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
		return
	}
	defer shorterApp.Close()

	logger.Infow(
		"Starting server",
		"addr", c.Addr,
	)

	go func() {
		logger.Infof("PProf endpoints available at /debug/pprof/*")
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			logger.Errorf("error starting PProf listener:", err)
		}
	}()

	if err := http.ListenAndServe(shorterApp.Addr, shorterApp.Router); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
