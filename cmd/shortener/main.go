// The package provides the ability to create, store, and receive short links.
//
// environment variables:
//
//	FILE_STORAGE_PATH | File storage paths specified via an environment variable.
//	SERVER_ADDRESS    | Server address defined by an environment variable.
//	BASE_URL          | Short link base URL configured via an environment variable.
//	DATABASE_DSN      | PostgreSQL Data Source Name received from an environment variable.
//
// command-line arguments:
//
//	-a | Command-line argument for server address.
//	-f | Command-line option specifying file storage paths.
//	-b | Base URL for short links passed via command-line.
//	-d | Postgres DSN given on the command line.
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
		err := http.ListenAndServe(c.PProfAddr, nil)
		if err != nil {
			logger.Errorf("error starting PProf listener: %s", err)
		}
	}()

	if err := http.ListenAndServe(shorterApp.Addr, shorterApp.Router); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
