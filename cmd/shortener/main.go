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
//
// Compile the program into a binary named 'shortenerapp', embedding version, build timestamp, and Git commit hash,
// then immediately execute the compiled binary.
//
//	go build -ldflags="\
//	-X 'main.buildVersion=v1.0' \
//	-X 'main.buildDate=`date '+%Y-%m-%dT%H:%M:%SZ'`' \
//	-X 'main.buildCommit=`git rev-parse HEAD`'" \
//	-o shortenerapp && ./shortenerapp
package main

import (
	"fmt"
	"main/internal/adapters"
	"main/internal/app"
	"main/internal/config"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c := config.Parse(logger)

	a, err := app.NewApp(c, logger)
	if err != nil {
		logger.Fatalw(err.Error(), "event", "initialize application")
		return
	}
	defer a.Close()

	doneCh := make(chan struct{})

	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-stopChan:
			logger.Info("Received shutdown signal.")
			a.Close()
		case <-doneCh:
			logger.Info("Application closed normally.")
		}
		close(doneCh)
	}()

	if err := a.StartServer(); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}

	<-doneCh
}
