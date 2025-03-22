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

type App struct {
	Addr     string
	Router   *chi.Mux
	Services *Services
}

func (a *App) Close() error {
	err := a.Services.Close()
	if err != nil {
		return err
	}
	return nil
}

type Handlers struct {
	links  interfaces.LinkHandlers
	health interfaces.HealthHandlers
}

type Services struct {
	links      interfaces.LinksService
	users      interfaces.UsersService
	health     interfaces.HealthService
	Repository *Repository
}

func (s *Services) Close() error {
	err := s.Repository.Close()
	if err != nil {
		return err
	}
	return nil
}

type Repository struct {
	links    interfaces.LinksRepository
	users    interfaces.UsersRepository
	health   interfaces.HealthRepository
	Database interfaces.DB
}

func (s *Repository) Close() error {
	err := s.Database.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewApp(c *config.Config) *App {
	s := NewServices(c)
	h := NewHandlers(s)
	r := NewRouters(h)

	return &App{
		Addr:     c.Addr,
		Router:   r,
		Services: s,
	}
}

func NewHandlers(s *Services) *Handlers {
	middleware.UserService = s.users
	return &Handlers{
		links:  NewLinksHandlers(s.links),
		health: NewHealthHandlers(s.health),
	}
}

func NewServices(c *config.Config) *Services {
	repository, err := NewRepository(c)
	if err != nil {
		panic(err)
	}
	return &Services{
		links:      services.NewLinksService(c, repository.links),
		users:      services.NewUserService(repository.users),
		health:     services.NewHealthService(repository.health),
		Repository: repository,
	}
}

func NewRepository(c *config.Config) (*Repository, error) {
	var repository *Repository

	logger := adapters.GetLogger()

	if c.PostgresDNS == nil {
		db, err := memory.NewInMemoryDB(c.StorageFilePaths, logger)
		if err != nil {
			return nil, err
		}
		repository = NewInMemoryRepository(db)
	} else {
		db, err := psql.NewPostgresDB(c.PostgresDNS, logger)
		if err != nil {
			return nil, err
		}
		repository = NewPostgresRepository(db)
	}
	return repository, nil
}

func NewInMemoryRepository(db *memory.InMemoryDB) *Repository {
	return &Repository{
		links: memory.NewLinksRepository(db),
		//users:    memory.NewUserRepository(db),
		health:   memory.NewHealthRepository(db),
		Database: db,
	}
}

func NewPostgresRepository(db *psql.PostgresDB) *Repository {
	return &Repository{
		links:    psql.NewLinksRepository(db),
		users:    psql.NewUsersRepository(db),
		health:   psql.NewHealthRepository(db),
		Database: db,
	}
}
