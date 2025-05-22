package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/adapters"
	"main/internal/adapters/database/memory"
	"main/internal/adapters/database/psql"
	"main/internal/config"
	"main/internal/interfaces"
	"main/internal/middleware"
	"main/internal/services"
)

// App encapsulates the core application state and dependencies.
type App struct {
	Addr     string    // Binding address for the HTTP server.
	Router   *chi.Mux  // Main router for handling HTTP requests.
	Services *Services // Aggregation of application services.
}

// Close gracefully cleans up running services and dependencies.
func (a *App) Close() error {
	err := a.Services.Close()
	if err != nil {
		return err
	}
	return nil
}

// Handlers organizes HTTP handlers into a coherent structure.
type Handlers struct {
	links  interfaces.LinkHandlers   // Handler for link-related operations.
	health interfaces.HealthHandlers // Handler for health check endpoints.
	users  interfaces.UsersHandlers  // Handler for user-specific operations.
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

// NewApp constructs a fully-configured application instance.
func NewApp(c *config.Config) (*App, error) {
	s, err := NewServices(c)
	if err != nil {
		return nil, err
	}
	h := NewHandlers(s)
	r := NewRouters(h)

	app := &App{
		Addr:     c.Addr,
		Router:   r,
		Services: s,
	}
	return app, nil
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
func NewServices(c *config.Config) (*Services, error) {
	repository, err := NewRepository(c)
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
func NewRepository(c *config.Config) (*Repository, error) {
	var repository *Repository

	logger := adapters.GetLogger()

	if c.PostgresDNS == nil {
		db, err := memory.NewInMemoryDB(c.StorageFilePaths, logger)
		if err != nil {
			return nil, err
		}
		logger.Info("Create InMemoryDB Connection")
		repository = NewInMemoryRepository(db)
	} else {
		db, err := psql.NewPostgresDB(c.PostgresDNS, logger)
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
