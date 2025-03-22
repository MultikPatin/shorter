package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/adapters"
	"main/internal/adapters/database/psql"
	"main/internal/config"
	"main/internal/interfaces"
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
	links interfaces.LinkHandlers
	users interfaces.UsersHandlers
}

type Services struct {
	links interfaces.LinksService
	users interfaces.UsersService
}

func (s *Services) Close() error {
	var err error

	err = s.links.Close()
	if err != nil {
		return err
	}
	err = s.users.Close()
	if err != nil {
		return err
	}
	return nil
}

type Repository struct {
	links interfaces.LinksRepository
	users interfaces.UsersRepository
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
	return &Handlers{
		links: NewLinksHandlers(s.links),
		users: NewUsersHandlers(s.users),
	}
}

func NewServices(c *config.Config) *Services {
	repository, err := NewRepository(c)
	if err != nil {
		panic(err)
	}
	return &Services{
		links: services.NewLinksService(c, repository.links),
		users: services.NewUserService(c, repository.users),
	}
}

func NewRepository(c *config.Config) (*Repository, error) {
	logger := adapters.GetLogger()

	conn, err := psql.NewPostgresConnection(c.PostgresDNS, logger)

	links, err := adapters.NewLinksRepository(c, logger)
	if err != nil {
		return nil, err
	}
	users, err := adapters.NewUserRepository(c, logger)
	if err != nil {
		return nil, err
	}

	return &Repository{
		links: links,
		users: users,
	}, nil
}
