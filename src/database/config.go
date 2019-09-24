package database

// Config is a structure for storing data to connect to Postgres
type Config struct {
	Dialect      string `json:"dialect"`
	User         string `json:"user"`
	DataBaseName string `json:"db_name"`
	SSLMode      string `json:"ssl_mode"`
	Password     string `json:"password"`
}
