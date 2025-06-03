package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"main/internal/adapters/database/memory"
	"main/internal/adapters/database/psql"
	"main/internal/config"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/middleware"
	"main/internal/services"
	"net/http"
	"sync"
)

// Handlers organizes HTTP handlers into a coherent structure.
type Handlers struct {
	links  interfaces.LinkHandlers   // Handler for link-related operations.
	health interfaces.HealthHandlers // Handler for health check endpoints.
	users  interfaces.UsersHandlers  // Handler for user-specific operations.
}

// App encapsulates the core application state and dependencies.
type App struct {
	Router   *chi.Mux           // Main router for handling HTTP requests.
	Services *Services          // Business logic and service instances.
	log      *zap.SugaredLogger // Configuration settings.
	conf     *config.Config     // Logger for application-wide logging.
	cancel   context.CancelFunc // Function to cancel the application context.
	ctx      context.Context    // Application context for signal propagation.
	wg       sync.WaitGroup     // Wait group for tracking background tasks.
}

// NewApp constructs a fully-configured application instance.
func NewApp(c *config.Config, l *zap.SugaredLogger) (*App, error) {
	s, err := NewServices(c, l)
	if err != nil {
		return nil, err
	}
	h := NewHandlers(s)
	r := NewRouters(h)

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		log:      l,
		conf:     c,
		Router:   r,
		Services: s,
		ctx:      ctx,
		cancel:   cancel,
	}
	return app, nil
}

// StartServer boots the primary HTTP server and handles graceful shutdowns.
func (a *App) StartServer() error {
	a.wg.Add(1)
	go a.startPPROFServer()
	a.log.Infow("Starting server", "addr", a.conf.Addr)

	srv := &http.Server{
		Addr:    a.conf.Addr,
		Handler: a.Router,
	}

	errCh := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("ListenAndServe failed: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-a.ctx.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), constants.ServerShutdownTime)
		defer cancelShutdown()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			a.log.Fatalw(err.Error(), "event", "server shutdown")
		}
		return nil
	}
}

// startPPROFServer launches a secondary server dedicated to performance profiling tools.
func (a *App) startPPROFServer() {
	defer a.wg.Done()

	a.log.Infow("Starting PPROF server", "addr", a.conf.PProfAddr)

	srv := &http.Server{
		Addr:    a.conf.PProfAddr,
		Handler: nil,
	}

	errCh := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("ListenAndServe PPROF failed: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		log.Println("Error in PPROF server:", err)
	case <-a.ctx.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), constants.ServerShutdownTime)
		defer cancelShutdown()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			a.log.Fatalw(err.Error(), "event", "server PPROF shutdown")
		}
	}
}

// Close gracefully cleans up running services and dependencies.
func (a *App) Close() error {
	a.cancel()
	a.wg.Wait()
	err := a.Services.Close()
	if err != nil {
		return err
	}
	return nil
}

// Services orchestrates service-level behavior and lifecycle management.
type Services struct {
	links      interfaces.LinksService  // Service for link-related operations.
	health     interfaces.HealthService // Service for health-related operations.
	users      interfaces.UsersService  // Service for user-specific operations.
	Repository *Repository              // Encapsulation of repository access.
}

// Close shuts down the services and propagates cleanup.
func (s *Services) Close() error {
	err := s.Repository.Close()
	if err != nil {
		return err
	}
	return nil
}

// Repository abstracts the interaction with the underlying data store.
type Repository struct {
	links    interfaces.LinksRepository  // Repository for link operations.
	users    interfaces.UsersRepository  // Repository for user operations.
	health   interfaces.HealthRepository // Repository for health checks.
	Database interfaces.DB               // Low-level database connection.
}

// Close terminates the underlying database connection.
func (s *Repository) Close() error {
	err := s.Database.Close()
	if err != nil {
		return err
	}
	return nil
}

// NewHandlers builds a set of HTTP handlers from the provided services.
func NewHandlers(s *Services) *Handlers {
	return &Handlers{
		links:  NewLinksHandlers(s.links),
		health: NewHealthHandlers(s.health),
		users:  NewUsersHandlers(s.users),
	}
}

// NewServices configures the application's service layer based on the configuration.
func NewServices(c *config.Config, l *zap.SugaredLogger) (*Services, error) {
	repository, err := NewRepository(c, l)
	if err != nil {
		return nil, err
	}
	if repository.users != nil {
		middleware.UserService = services.NewUserService(repository.users)
	} else {
		middleware.UserService = nil
	}
	return &Services{
		links:      services.NewLinksService(c, repository.links),
		health:     services.NewHealthService(repository.health),
		users:      services.NewUserService(repository.users),
		Repository: repository,
	}, nil
}

// NewRepository selects and initializes the appropriate repository based on configuration.
func NewRepository(c *config.Config, logger *zap.SugaredLogger) (*Repository, error) {
	var repository *Repository

	if c.PostgresDSN == nil {
		db, err := memory.NewInMemoryDB(c.StorageFilePaths, logger)
		if err != nil {
			return nil, err
		}
		logger.Info("Create InMemoryDB Connection")
		repository = NewInMemoryRepository(db)
	} else {
		db, err := psql.NewPostgresDB(c.PostgresDSN, logger)
		if err != nil {
			return nil, err
		}
		logger.Info("Create PostgresDB Connection")
		repository = NewPostgresRepository(db)
	}
	return repository, nil
}

// NewInMemoryRepository constructs a repository using an in-memory database.
func NewInMemoryRepository(db *memory.InMemoryDB) *Repository {
	return &Repository{
		links:    memory.NewLinksRepository(db),
		users:    nil,
		health:   memory.NewHealthRepository(db),
		Database: db,
	}
}

// NewPostgresRepository constructs a repository using a PostgreSQL database.
func NewPostgresRepository(db *psql.PostgresDB) *Repository {
	return &Repository{
		links:    psql.NewLinksRepository(db),
		users:    psql.NewUsersRepository(db),
		health:   psql.NewHealthRepository(db),
		Database: db,
	}
}
