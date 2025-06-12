// Package config handles parsing environment variables and command-line arguments into a unified configuration structure.
//
// environment variables:
//
//	FILE_STORAGE_PATH | File storage paths specified via an environment variable.
//	SERVER_ADDRESS    | Server address defined by an environment variable.
//	BASE_URL          | Short link base URL configured via an environment variable.
//	DATABASE_DSN      | PostgreSQL Data Source Name received from an environment variable.
//	ENABLE_HTTPS      | Indicates whether HTTPS is enabled for the server.
//	CONFIG      	  | Name of the configuration file.
//
// command-line arguments:
//
//	-a | Command-line argument for server address.
//	-f | Command-line option specifying file storage paths.
//	-b | Base URL for short links passed via command-line.
//	-d | Postgres DSN given on the command line.
//	-s | Indicates whether HTTPS is enabled for the server ("true", "yes", "1" -> true, "false", "no", "0" -> false).
//	-c | Name of the configuration file.
//
// config file:
//
//	config.json | Configuration file in JSON format.
//
//	file_storage_path | File storage paths specified via an environment variable.
//	server_address    | Server address defined by an environment variable.
//	base_url          | Short link base URL configured via an environment variable.
//	database_dsn      | PostgreSQL Data Source Name received from an environment variable.
//	enable_https      | Indicates whether HTTPS is enabled for the server.
package config
