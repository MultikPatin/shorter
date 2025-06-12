package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

// cmdConfig holds configuration settings obtained from command-line flags.
type cmdConfig struct {
	Addr             string // Command-line argument for server address.
	StorageFilePaths string // Command-line option specifying file storage paths.
	ShortLinkPrefix  string // Base URL for short links passed via command-line.
	PostgresDSN      string // Postgres DSN given on the command line.
	HTTPSEnable      string // Indicates whether HTTPS is enabled for the server.
	ConfFile         string // Name of the configuration file.
}

// servHost encapsulates information about the network service's host and port.
type servHost struct {
	Host string // Hostname part of the net address.
	Port int    // Port number component of the net address.
}

// parseCmd processes command-line flags to extract configuration values.
func parseCmd() (*cmdConfig, error) {
	cfg := &cmdConfig{}

	hostPort := new(servHost)
	_ = flag.Value(hostPort)

	flag.StringVar(&cfg.PostgresDSN, "d", "", "Postgres DSN")
	flag.StringVar(&cfg.ShortLinkPrefix, "b", "", "Short link server")
	flag.StringVar(&cfg.StorageFilePaths, "f", "", "Path to storage file")
	flag.StringVar(&cfg.HTTPSEnable, "s", "0", "HTTPS is enabled")
	flag.StringVar(&cfg.HTTPSEnable, "c", "", "Name of the configuration file")
	flag.Var(hostPort, "a", "Network address host:port")
	flag.Parse()

	cfg.Addr = hostPort.String()
	return cfg, nil
}

// String returns the formatted representation of the ServHost instance.
func (a *servHost) String() string {
	a.normalize()
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Set parses a string input representing a host-port pair and updates the ServHost fields accordingly.
func (a *servHost) Set(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return errors.New("address must be in format host:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	a.Host = parts[0]
	a.Port = port
	a.normalize()
	return nil
}

// normalize ensures valid defaults for empty or zero-value fields.
func (a *servHost) normalize() {
	if a.Port == 0 {
		a.Port = 8080
	}
	if a.Host == "" {
		a.Host = "localhost"
	}
}
