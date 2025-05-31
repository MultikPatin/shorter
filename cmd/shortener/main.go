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
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger := adapters.GetLogger()
	defer adapters.SyncLogger()

	c := config.Parse(logger)

	a, err := app.NewApp(c, logger)
	if err != nil {
		logger.Fatalw(err.Error(), "event", "initialize application")
		return
	}
	defer a.Close()

	var wg sync.WaitGroup // Используем WaitGroup для ожидания завершения процессов

	wg.Add(1)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopChan
		a.Close()
		wg.Done() // Уведомляем о завершении процесса закрытия приложения
	}()

	if err := a.StartServer(); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}

	wg.Wait() // Ждем пока завершится закрытие сервера
}
