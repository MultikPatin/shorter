package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"main/internal/adapters/database/memory"
	"main/internal/adapters/database/psql"
	"main/internal/cert"
	"main/internal/config"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/middleware"
	pb "main/internal/proto"
	"main/internal/servers"
	"main/internal/services"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// TLSFiles contain paths to the TLS certificate and private key files
type TLSFiles struct {
	CertFile string
	KeyFile  string
}

// getTLSFiles returns the paths to the TLS certificate and private key files.
// If either file is missing, it generates them using cert.GenerateTLSFiles.
func getTLSFiles(exeDir string) (*TLSFiles, error) {
	certFile := filepath.Join(exeDir, constants.CertFile)
	keyFile := filepath.Join(exeDir, constants.KeyFile)

	if _, err := os.Stat(certFile); err != nil {
		if os.IsNotExist(err) {
			err := cert.GenerateTLSFiles(certFile, keyFile)
			if err != nil {
				return nil, err
			}
			return &TLSFiles{
				CertFile: certFile,
				KeyFile:  keyFile,
			}, nil
		}
		return nil, err
	}
	if _, err := os.Stat(keyFile); err != nil {
		if os.IsNotExist(err) {
			err := cert.GenerateTLSFiles(certFile, keyFile)
			if err != nil {
				return nil, err
			}
			return &TLSFiles{
				CertFile: certFile,
				KeyFile:  keyFile,
			}, nil
		}
		return nil, err
	}

	return &TLSFiles{
		CertFile: certFile,
		KeyFile:  keyFile,
	}, nil
}

// App encapsulates the core application state and dependencies.
type App struct {
	Router   *chi.Mux           // Main router for handling HTTP requests.
	Services *Services          // Business logic and service instances.
	Servers  *Servers           // gRPC servers for handling gRPC requests.
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
	g := NewGRPCServers(c, s)

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		log:      l,
		conf:     c,
		Router:   r,
		Services: s,
		Servers:  g,
		ctx:      ctx,
		cancel:   cancel,
	}
	return app, nil
}

// StartServer boots the primary HTTP server and handles graceful shutdowns.
func (a *App) StartServer() error {
	a.wg.Add(1)

	go a.startPPROFServer()
	go a.startGRPCServer()

	a.log.Infow("Starting server", "addr", a.conf.Addr)
	a.log.Info("HTTPS status: ", a.conf.HTTPSEnable)

	srv := &http.Server{
		Addr:    a.conf.Addr,
		Handler: a.Router,
	}
	errCh := make(chan error)

	if a.conf.HTTPSEnable {
		tlsFiles, err := getTLSFiles(a.conf.ExecutableDir)
		if err != nil {
			return err
		}
		go func() {
			if err := srv.ListenAndServeTLS(tlsFiles.CertFile, tlsFiles.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("ListenAndServeTLS failed: %w", err)
			}
			close(errCh)
		}()
	} else {
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("ListenAndServe failed: %w", err)
			}
			close(errCh)
		}()
	}

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

// startGRPCServer launches a secondary server dedicated to performance profiling tools.
func (a *App) startGRPCServer() {
	defer a.wg.Done()

	a.log.Infow("Starting gRPC server", "addr", a.conf.GRPCAddr)

	listen, err := net.Listen("tcp", a.conf.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	pb.RegisterLinksServer(srv, &servers.LinksServer{})

	errCh := make(chan error)
	go func() {
		if err := srv.Serve(listen); err != nil {
			errCh <- fmt.Errorf("ListenAndServe gRPC failed: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		log.Println("Error in gRPC server:", err)
	case <-a.ctx.Done():
		_, cancelShutdown := context.WithTimeout(context.Background(), constants.ServerShutdownTime)
		defer cancelShutdown()
		srv.GracefulStop()
		if err := listen.Close(); err != nil {
			a.log.Fatalw(err.Error(), "event", "TCP listen shutdown")
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

// Handlers organizes HTTP handlers into a coherent structure.
type Handlers struct {
	links  interfaces.LinkHandlers   // Handler for link-related operations.
	health interfaces.HealthHandlers // Handler for health check endpoints.
	users  interfaces.UsersHandlers  // Handler for user-specific operations.
	stats  interfaces.StatsHandlers  // Handler for statistic data operations.
}

// NewHandlers builds a set of HTTP handlers from the provided services.
func NewHandlers(s *Services) *Handlers {
	return &Handlers{
		links:  NewLinksHandlers(s.links),
		health: NewHealthHandlers(s.health),
		users:  NewUsersHandlers(s.users),
		stats:  NewStatsHandlers(s.stats),
	}
}

// Servers organizes gRPC handlers into a coherent structure.
type Servers struct {
	links interfaces.LinksServer
}

// NewGRPCServers builds a set of gRPC handlers from the provided services.
func NewGRPCServers(c *config.Config, s *Services) *Servers {
	return &Servers{
		links: servers.NewLinksServer(s.links, c.Addr),
	}
}

// Services orchestrates service-level behavior and lifecycle management.
type Services struct {
	links      interfaces.LinksService  // Service for link-related operations.
	health     interfaces.HealthService // Service for health-related operations.
	users      interfaces.UsersService  // Service for user-specific operations.
	stats      interfaces.StatsService  // Service for statistics-related operations.
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
	stats    interfaces.StatsRepository  // Repository for statistics operations.
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
		stats:      services.NewStatsService(c, repository.stats),
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
		stats:    memory.NewStatsRepository(db),
		Database: db,
	}
}

// NewPostgresRepository constructs a repository using a PostgreSQL database.
func NewPostgresRepository(db *psql.PostgresDB) *Repository {
	return &Repository{
		links:    psql.NewLinksRepository(db),
		users:    psql.NewUsersRepository(db),
		health:   psql.NewHealthRepository(db),
		stats:    psql.NewStatsRepository(db),
		Database: db,
	}
}
